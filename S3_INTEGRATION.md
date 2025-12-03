# Интеграция S3 хранилища Beget

## 📋 Содержание

- [Конфигурация](#конфигурация)
- [Архитектура](#архитектура)
- [API эндпоинты](#api-эндпоинты)
- [Примеры использования](#примеры-использования)
- [Миграция с локального хранилища](#миграция-с-локального-хранилища)

---

## Конфигурация

### Переменные окружения (.env)

```env
# S3 Beget Storage
S3_ENDPOINT=https://s3.ru1.storage.beget.cloud
S3_BUCKET_NAME=c15b4d655f70-medvito-data
S3_ACCESS_KEY=I6I3KOJ2YO3TN08TDJAI
S3_SECRET_KEY=5up6F9kLNHRGmPIczdqAVZgBNgKhFpAGJ1JnCJUY
S3_REGION=ru1
```

### Установка зависимостей

```bash
yarn add @aws-sdk/client-s3 uuid
yarn add -D @types/uuid
```

---

## Архитектура

### S3Service

Основной сервис для работы с S3 хранилищем.

**Методы:**

- `uploadFile(file, folder)` - загрузка одного файла
- `uploadFiles(files, folder)` - загрузка нескольких файлов
- `deleteFile(fileUrl)` - удаление файла
- `deleteFiles(fileUrls)` - удаление нескольких файлов
- `getFile(key)` - получение файла (для приватных файлов)

### Структура папок в бакете

```
c15b4d655f70-medvito-data/
├── products/          # Изображения товаров
│   ├── uuid-1.jpg
│   ├── uuid-2.png
│   └── ...
├── users/             # Аватары пользователей
│   ├── uuid-1.jpg
│   └── ...
└── uploads/           # Другие файлы
    └── ...
```

---

## API эндпоинты

### 1. Загрузка одного файла

```http
POST /s3/upload
Content-Type: multipart/form-data

Body:
- file: [binary]
- folder: "products" (optional)
```

**Ответ:**

```json
{
  "message": "Файл успешно загружен",
  "url": "https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/uuid-123.jpg"
}
```

### 2. Загрузка нескольких файлов

```http
POST /s3/upload-multiple
Content-Type: multipart/form-data

Body:
- files: [binary array, до 10 файлов]
- folder: "products" (optional)
```

**Ответ:**

```json
{
  "message": "Успешно загружено 3 файлов",
  "urls": [
    "https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/uuid-1.jpg",
    "https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/uuid-2.jpg",
    "https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/uuid-3.jpg"
  ]
}
```

### 3. Удаление файла

```http
DELETE /s3/delete
Content-Type: application/json

Body:
{
  "fileUrl": "https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/uuid-123.jpg"
}
```

**Ответ:**

```json
{
  "message": "Файл успешно удалён"
}
```

### 4. Удаление нескольких файлов

```http
DELETE /s3/delete-multiple
Content-Type: application/json

Body:
{
  "fileUrls": [
    "https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/uuid-1.jpg",
    "https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/uuid-2.jpg"
  ]
}
```

**Ответ:**

```json
{
  "message": "Успешно удалено 2 файлов"
}
```

---

## Примеры использования

### Пример 1: Создание товара с изображениями

**Frontend (JavaScript):**

```javascript
const formData = new FormData();
formData.append('name', 'iPhone 15 Pro');
formData.append('price', '120000');
formData.append('state', 'NEW');
formData.append('description', 'Новый iPhone');
formData.append('address', 'Москва');
formData.append('categoryId', '1');
formData.append('subcategoryId', '1');

// Добавляем изображения
const images = document.getElementById('images').files;
for (let i = 0; i < images.length; i++) {
  formData.append('images', images[i]);
}

const response = await fetch('http://localhost:3002/product/create', {
  method: 'POST',
  body: formData,
  credentials: 'include', // Для отправки cookies (session_id)
});

const result = await response.json();
console.log('Product created:', result);
```

**Ответ:**

```json
{
  "message": "Продукт успешно создан",
  "product": {
    "id": 1,
    "name": "iPhone 15 Pro",
    "price": 120000,
    "images": [
      "https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/uuid-1.jpg",
      "https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/uuid-2.jpg"
    ],
    ...
  }
}
```

### Пример 2: Обновление товара (добавление изображений)

**Frontend:**

```javascript
const formData = new FormData();
formData.append('description', 'Обновлённое описание');

// Добавляем новые изображения
const newImages = document.getElementById('new-images').files;
for (let i = 0; i < newImages.length; i++) {
  formData.append('images', newImages[i]);
}

const response = await fetch('http://localhost:3002/product/123', {
  method: 'PATCH',
  body: formData,
  credentials: 'include',
});

const result = await response.json();
console.log('Product updated:', result);
```

### Пример 3: Удаление товара (автоматически удаляются изображения из S3)

**Frontend:**

```javascript
const response = await fetch('http://localhost:3002/product/123', {
  method: 'DELETE',
  credentials: 'include',
});

const result = await response.json();
console.log('Product deleted:', result);
// Все изображения товара автоматически удалены из S3
```

### Пример 4: Использование в другом сервисе

**user.service.ts:**

```typescript
import { S3Service } from 'src/s3/s3.service';

@Injectable()
export class UserService {
  constructor(
    private readonly prisma: PrismaService,
    private readonly s3Service: S3Service,
  ) {}

  async updateAvatar(userId: number, file: Express.Multer.File) {
    // Получаем пользователя
    const user = await this.prisma.user.findUnique({
      where: { id: userId },
    });

    // Если есть старый аватар - удаляем его из S3
    if (user.avatar) {
      await this.s3Service.deleteFile(user.avatar);
    }

    // Загружаем новый аватар
    const avatarUrl = await this.s3Service.uploadFile(file, 'users');

    // Обновляем в БД
    await this.prisma.user.update({
      where: { id: userId },
      data: { avatar: avatarUrl },
    });

    return { message: 'Аватар обновлён', avatarUrl };
  }
}
```

**user.module.ts:**

```typescript
import { S3Module } from 'src/s3/s3.module';

@Module({
  imports: [
    PrismaModule,
    S3Module, // Импортируем S3Module
  ],
  controllers: [UserController],
  providers: [UserService],
  exports: [UserService],
})
export class UserModule {}
```

---

## Миграция с локального хранилища

### Изменения в ProductService

#### До (локальное хранилище):

```typescript
async createProduct(dto, fileNames: string[], userId: number) {
  const imagePaths = fileNames.map(
    (fileName) => `/uploads/product/${fileName}`,
  );

  const product = await this.prisma.product.create({
    data: {
      ...dto,
      images: imagePaths,
    },
  });

  return {
    product: {
      ...product,
      images: imagePaths.map((path) => `${this.baseUrl}${path}`),
    },
  };
}
```

#### После (S3 хранилище):

```typescript
async createProduct(dto, files: Express.Multer.File[], userId: number) {
  // Загружаем изображения в S3
  const imageUrls = files && files.length > 0
    ? await this.s3Service.uploadFiles(files, 'products')
    : [];

  const product = await this.prisma.product.create({
    data: {
      ...dto,
      images: imageUrls, // URL уже полные из S3
    },
  });

  return {
    product: {
      ...product,
      images: product.images, // URL уже полные
    },
  };
}
```

### Изменения в ProductModule

#### До:

```typescript
MulterModule.register({
  storage: diskStorage({
    destination: './uploads/product',
    filename: (req, file, callback) => {
      const uniqueSuffix = Date.now() + '-' + Math.round(Math.random() * 1e9);
      const ext = path.extname(file.originalname);
      callback(null, file.fieldname + '-' + uniqueSuffix + ext);
    },
  }),
});
```

#### После:

```typescript
MulterModule.register({
  storage: 'memory', // Используем память вместо диска для S3
});
```

### Изменения в ProductController

#### До:

```typescript
async createProduct(@Body() dto, @UploadedFiles() images) {
  const fileNames = images.map((file) => file.filename);
  return this.productService.createProduct(dto, fileNames, userId);
}
```

#### После:

```typescript
async createProduct(@Body() dto, @UploadedFiles() images) {
  return this.productService.createProduct(dto, images, userId);
}
```

---

## Преимущества S3

✅ **Масштабируемость** - не ограничены дисковым пространством сервера  
✅ **Надёжность** - резервное копирование и отказоустойчивость  
✅ **CDN** - быстрая доставка контента пользователям  
✅ **Безопасность** - контроль доступа и шифрование  
✅ **Экономия** - оплата только за использованное пространство

---

## Тестирование

### Тест загрузки файла:

```bash
curl -X POST http://localhost:3002/s3/upload \
  -F "file=@/path/to/image.jpg" \
  -F "folder=products"
```

### Тест удаления файла:

```bash
curl -X DELETE http://localhost:3002/s3/delete \
  -H "Content-Type: application/json" \
  -d '{"fileUrl": "https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/uuid-123.jpg"}'
```

---

## Важные моменты

1. **ACL: public-read** - все загружаемые файлы доступны публично
2. **UUID имена** - файлы получают уникальные имена, предотвращая конфликты
3. **Автоматическое удаление** - при удалении товара файлы удаляются из S3
4. **Ошибки удаления** - не блокируют основной функционал (логируются в консоль)
5. **Memory storage** - файлы загружаются в память, затем отправляются в S3

---

## Дальнейшее развитие

- [ ] Добавить генерацию thumbnails (миниатюр изображений)
- [ ] Реализовать подписанные URL для приватных файлов
- [ ] Добавить валидацию типов и размеров файлов
- [ ] Реализовать сжатие изображений перед загрузкой
- [ ] Добавить прогресс-бар загрузки на frontend
