# План рефакторинга системы категорий и фильтрации

> **Цель:** Улучшить фильтрацию товаров для фронтенда без поломки админки и текущего функционала

---

## 🎯 Проблемы текущей системы

### Архитектура

```
Category (Электроника)
  └── SubCategory (Телефоны)
      └── SubcategoryType (Смартфоны)
          └── TypeField (Память, Цвет, Процессор)
              └── ProductFieldValue (128GB, Черный, Snapdragon)
```

**Проблемы:**

- ❌ Фиксированная 3-уровневая вложенность (нельзя сделать 2 или 4 уровня)
- ❌ Опечатка в названии модели `SubcategotyType` → `SubcategoryType`
- ❌ Нет индексов для быстрого поиска
- ❌ N+1 проблемы при загрузке фильтров
- ❌ Нет кеширования структуры категорий
- ❌ Сложно добавить новые уровни вложенности

---

## 🚀 План действий (3 этапа)

### ✅ Этап 1: Быстрые улучшения (БЕЗ ЛОМКИ ФУНКЦИОНАЛА)

**Что можно сделать прямо сейчас:**

#### 1.1 Исправить схему базы данных

```prisma
// Исправить опечатку
model SubcategoryType {  // было: SubcategotyType
  id             Int       @id @default(autoincrement())
  name           String
  subcategoryId  Int

  subcategory    SubCategory     @relation(fields: [subcategoryId], references: [id])
  typeFields     TypeField[]
  products       Product[]

  // ✅ Добавить индексы для производительности
  @@index([subcategoryId])
  @@index([name])
  @@map("subcategory_types")
}

model Product {
  // ✅ Добавить индексы для фильтрации
  @@index([categoryId, subcategoryId])
  @@index([subcategoryTypeId])
  @@index([price])
  @@index([createdAt])
  @@index([moderateState, createdAt])
}

model ProductFieldValue {
  // ✅ Составной индекс для быстрого поиска
  @@index([fieldId, value])
  @@index([productId])
}
```

**Команда:**

```bash
# Применить миграцию
npm run prisma:migrate:dev -- --name fix_category_typo_and_add_indexes
```

---

#### 1.2 Создать эндпоинт для получения дерева категорий с кешированием

**Новый файл:** `src/category/dto/category-tree.dto.ts`

```typescript
export class CategoryTreeDto {
  id: number;
  name: string;
  image: string | null;
  subcategories: SubcategoryTreeDto[];
}

export class SubcategoryTreeDto {
  id: number;
  name: string;
  types: SubcategoryTypeTreeDto[];
}

export class SubcategoryTypeTreeDto {
  id: number;
  name: string;
  fieldsCount: number;
}
```

**Новый эндпоинт:** `src/category/category.controller.ts`

```typescript
@Get('tree')
@UseInterceptors(CacheInterceptor)
@CacheTTL(3600000) // 1 час
async getCategoryTree() {
  return this.categoryService.getCategoryTree();
}
```

**Реализация в сервисе:**

```typescript
async getCategoryTree(): Promise<CategoryTreeDto[]> {
  const categories = await this.prisma.category.findMany({
    include: {
      subcategories: {
        include: {
          types: {
            include: {
              _count: {
                select: { typeFields: true }
              }
            }
          }
        }
      }
    },
    orderBy: { name: 'asc' }
  });

  return categories.map(cat => ({
    id: cat.id,
    name: cat.name,
    image: cat.image,
    subcategories: cat.subcategories.map(sub => ({
      id: sub.id,
      name: sub.name,
      types: sub.types.map(type => ({
        id: type.id,
        name: type.name,
        fieldsCount: type._count.typeFields
      }))
    }))
  }));
}
```

---

#### 1.3 Создать эндпоинт для получения доступных фильтров

**Новый файл:** `src/product/dto/filter-options.dto.ts`

```typescript
export class FilterOptionsDto {
  categories: CategoryFilterDto[];
  priceRange: { min: number; max: number };
  fields: FieldFilterDto[];
}

export class CategoryFilterDto {
  id: number;
  name: string;
  productsCount: number;
}

export class FieldFilterDto {
  id: number;
  name: string;
  type: string;
  values: FieldValueDto[];
}

export class FieldValueDto {
  value: string;
  count: number;
}
```

**Новый эндпоинт:** `src/product/product.controller.ts`

```typescript
@Get('filters')
@UseInterceptors(CacheInterceptor)
@CacheTTL(600000) // 10 минут
async getFilterOptions(
  @Query('categoryId') categoryId?: string,
  @Query('subcategoryId') subcategoryId?: string,
  @Query('typeId') typeId?: string,
) {
  return this.productService.getFilterOptions({
    categoryId: categoryId ? +categoryId : undefined,
    subcategoryId: subcategoryId ? +subcategoryId : undefined,
    typeId: typeId ? +typeId : undefined,
  });
}
```

**Реализация:**

```typescript
async getFilterOptions(filters: {
  categoryId?: number;
  subcategoryId?: number;
  typeId?: number;
}): Promise<FilterOptionsDto> {
  const where = {
    moderateState: 'approved' as const,
    ...(filters.categoryId && { categoryId: filters.categoryId }),
    ...(filters.subcategoryId && { subcategoryId: filters.subcategoryId }),
    ...(filters.typeId && { subcategoryTypeId: filters.typeId }),
  };

  // Параллельные запросы для оптимизации
  const [priceStats, fieldValues] = await Promise.all([
    this.prisma.product.aggregate({
      where,
      _min: { price: true },
      _max: { price: true },
    }),

    this.prisma.productFieldValue.groupBy({
      by: ['fieldId', 'value'],
      where: {
        product: where
      },
      _count: true,
    })
  ]);

  // Получить информацию о полях
  const fieldIds = [...new Set(fieldValues.map(fv => fv.fieldId))];
  const fields = await this.prisma.typeField.findMany({
    where: { id: { in: fieldIds } },
    select: { id: true, name: true, type: true }
  });

  const fieldsMap = new Map(fields.map(f => [f.id, f]));
  const groupedValues = new Map<number, FieldValueDto[]>();

  fieldValues.forEach(fv => {
    if (!groupedValues.has(fv.fieldId)) {
      groupedValues.set(fv.fieldId, []);
    }
    groupedValues.get(fv.fieldId)!.push({
      value: fv.value,
      count: fv._count
    });
  });

  return {
    priceRange: {
      min: priceStats._min.price || 0,
      max: priceStats._max.price || 0,
    },
    fields: Array.from(fieldsMap.entries()).map(([id, field]) => ({
      id,
      name: field.name,
      type: field.type,
      values: groupedValues.get(id) || []
    }))
  };
}
```

---

#### 1.4 Улучшить эндпоинт поиска с фильтрацией

**Обновить:** `src/product/dto/product-query.dto.ts`

```typescript
export class ProductQueryDto {
  @IsOptional()
  @Type(() => Number)
  categoryId?: number;

  @IsOptional()
  @Type(() => Number)
  subcategoryId?: number;

  @IsOptional()
  @Type(() => Number)
  typeId?: number;

  @IsOptional()
  @Type(() => Number)
  minPrice?: number;

  @IsOptional()
  @Type(() => Number)
  maxPrice?: number;

  @IsOptional()
  search?: string;

  // Динамические фильтры по полям: ?fields[memory]=128GB&fields[color]=Черный
  @IsOptional()
  fields?: Record<string, string>;

  @IsOptional()
  @Type(() => Number)
  @Min(1)
  page?: number = 1;

  @IsOptional()
  @Type(() => Number)
  @Min(1)
  @Max(100)
  limit?: number = 20;

  @IsOptional()
  @IsIn(['price', 'createdAt', 'popularity'])
  sortBy?: string = 'createdAt';

  @IsOptional()
  @IsIn(['asc', 'desc'])
  sortOrder?: 'asc' | 'desc' = 'desc';
}
```

**Обновить эндпоинт:**

```typescript
@Get('search')
async searchProducts(@Query() query: ProductQueryDto) {
  return this.productService.searchProducts(query);
}
```

**Реализация с фильтрацией:**

```typescript
async searchProducts(query: ProductQueryDto) {
  const where: Prisma.ProductWhereInput = {
    moderateState: 'approved',
    ...(query.categoryId && { categoryId: query.categoryId }),
    ...(query.subcategoryId && { subcategoryId: query.subcategoryId }),
    ...(query.typeId && { subcategoryTypeId: query.typeId }),
    ...(query.minPrice && { price: { gte: query.minPrice } }),
    ...(query.maxPrice && { price: { lte: query.maxPrice } }),
    ...(query.search && {
      OR: [
        { name: { contains: query.search, mode: 'insensitive' } },
        { description: { contains: query.search, mode: 'insensitive' } }
      ]
    })
  };

  // Фильтрация по динамическим полям
  if (query.fields && Object.keys(query.fields).length > 0) {
    where.fieldValues = {
      some: {
        OR: Object.entries(query.fields).map(([fieldName, value]) => ({
          field: { name: fieldName },
          value: value
        }))
      }
    };
  }

  const skip = ((query.page || 1) - 1) * (query.limit || 20);
  const take = query.limit || 20;

  const [products, total] = await Promise.all([
    this.prisma.product.findMany({
      where,
      skip,
      take,
      orderBy: { [query.sortBy || 'createdAt']: query.sortOrder || 'desc' },
      include: {
        category: { select: { id: true, name: true } },
        subcategory: { select: { id: true, name: true } },
        subcategoryType: { select: { id: true, name: true } },
        user: { select: { id: true, name: true } },
        fieldValues: {
          include: {
            field: { select: { name: true, type: true } }
          }
        }
      }
    }),
    this.prisma.product.count({ where })
  ]);

  return {
    products,
    pagination: {
      total,
      page: query.page || 1,
      limit: query.limit || 20,
      pages: Math.ceil(total / (query.limit || 20))
    }
  };
}
```

---

### 📋 Что получит фронтенд после Этапа 1:

#### API для работы с категориями:

```
GET /category/tree
→ Полное дерево категорий с кешированием

GET /product/filters?categoryId=1&subcategoryId=2
→ Доступные фильтры (диапазон цен, динамические поля)

GET /product/search?categoryId=1&minPrice=1000&maxPrice=5000&fields[memory]=128GB&page=1&limit=20
→ Поиск с фильтрацией и пагинацией
```

#### Пример запроса с фронтенда:

```javascript
// 1. Получить дерево категорий (один раз при загрузке)
const categoryTree = await fetch('/category/tree');

// 2. Получить доступные фильтры для выбранной категории
const filters = await fetch('/product/filters?categoryId=1&subcategoryId=2');

// 3. Поиск с фильтрами
const results = await fetch(
  '/product/search?' +
    new URLSearchParams({
      categoryId: '1',
      subcategoryId: '2',
      minPrice: '1000',
      maxPrice: '5000',
      'fields[memory]': '128GB',
      'fields[color]': 'Черный',
      page: '1',
      limit: '20',
      sortBy: 'price',
      sortOrder: 'asc',
    }),
);
```

---

## 🔄 Этап 2: Гибкая структура категорий (В БУДУЩЕМ)

### Проблема:

Текущая структура не позволяет создать категории с разной глубиной:

- Одежда → Мужская → Футболки (3 уровня)
- Услуги → Ремонт (2 уровня)
- Недвижимость → Коммерческая → Офисы → Класса A (4 уровня)

### Решение: Adjacency List Pattern

**Новая модель (параллельно со старой):**

```prisma
model CategoryNode {
  id          Int            @id @default(autoincrement())
  name        String
  slug        String         @unique
  parentId    Int?
  level       Int            @default(0)
  order       Int            @default(0)
  isActive    Boolean        @default(true)

  // Рекурсивная связь
  parent      CategoryNode?  @relation("CategoryHierarchy", fields: [parentId], references: [id])
  children    CategoryNode[] @relation("CategoryHierarchy")

  // Связь с полями
  fields      CategoryField[]

  @@index([parentId])
  @@index([slug])
  @@index([level])
  @@map("category_nodes")
}

model CategoryField {
  id           Int          @id @default(autoincrement())
  categoryId   Int
  name         String
  type         String       // text, number, select, multiselect, range
  options      Json?        // для select/multiselect
  isRequired   Boolean      @default(false)
  order        Int          @default(0)

  category     CategoryNode @relation(fields: [categoryId], references: [id])

  @@index([categoryId])
  @@map("category_fields")
}
```

**Преимущества:**

- ✅ Любая глубина вложенности
- ✅ Легко добавлять/удалять уровни
- ✅ Можно получить путь: "Электроника / Телефоны / Смартфоны / Samsung"
- ✅ Можно получить все дочерние категории рекурсивно

**Синхронизация со старой структурой:**

```typescript
// Middleware для синхронизации при создании категории
async syncToNewStructure(oldCategory: Category) {
  await this.prisma.categoryNode.create({
    data: {
      name: oldCategory.name,
      slug: slugify(oldCategory.name),
      level: 0,
      // связь со старой моделью для обратной совместимости
    }
  });
}
```

---

## 🎨 Этап 3: Миграция (ПОСЛЕ ТЕСТИРОВАНИЯ)

После того как фронтенд и админка протестированы на новой структуре:

1. Миграция всех данных из старой структуры в новую
2. Обновление всех связей в Product
3. Удаление старых моделей (Category, SubCategory, SubcategoryType)
4. Удаление кода синхронизации

---

## 🚦 Текущий статус

### Что можно сделать БЕЗ РИСКА:

- ✅ Исправить опечатку в `SubcategotyType` → `SubcategoryType`
- ✅ Добавить индексы для производительности
- ✅ Создать эндпоинт `/category/tree`
- ✅ Создать эндпоинт `/product/filters`
- ✅ Улучшить эндпоинт `/product/search`
- ✅ Добавить кеширование

### Что требует планирования:

- ⚠️ Переход на древовидную структуру категорий
- ⚠️ Миграция данных
- ⚠️ Обновление админки

---

## 💡 Рекомендация

**Начни с Этапа 1** — это даст фронтенду все необходимое для фильтрации прямо сейчас:

- Быстрый поиск с фильтрами
- Кеширование
- Пагинация
- Производительность

**Этап 2 и 3** можно отложить до того момента, когда реально понадобится более гибкая структура категорий.

---

## 📝 Следующие шаги

1. **Исправить схему и создать миграцию:**

   ```bash
   npm run prisma:migrate:dev -- --name fix_category_typo_and_add_indexes
   ```

2. **Создать DTO для фильтрации:**
   - `category-tree.dto.ts`
   - `filter-options.dto.ts`
   - `product-query.dto.ts`

3. **Добавить эндпоинты:**
   - `GET /category/tree`
   - `GET /product/filters`
   - Улучшить `GET /product/search`

4. **Добавить кеширование:**
   - Настроить Redis для category tree
   - Настроить TTL для фильтров

5. **Протестировать на фронтенде**

---

## ❓ Вопросы для принятия решения

1. **Нужна ли гибкая древовидная структура прямо сейчас?**
   - Если нет → начинай с Этапа 1
   - Если да → обсуди план миграции с командой

2. **Какие категории планируются в будущем?**
   - Если всегда 3 уровня → текущая структура ОК
   - Если разная глубина → нужна древовидная структура

3. **Сколько времени на тестирование?**
   - Этап 1: 1-2 дня разработки + 1 день тестирования
   - Этап 2: 3-5 дней разработки + 2-3 дня тестирования
   - Этап 3: 2-3 дня миграции + неделя мониторинга
