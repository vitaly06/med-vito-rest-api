# API модерации — документация для фронта

**Base URL:** `https://your-api.com`  
**Авторизация:** cookie `session_id` (тот же механизм, что и у других admin-эндпоинтов). Требуется роль `admin`.

---

## 1. Список товаров на модерации

```
GET /admin/moderation/products
```

### Query-параметры

| Параметр | Тип | Значения | По умолчанию | Описание |
|---|---|---|---|---|
| `filter` | string | `ALL` / `DENIED` / `MANUAL` / `APPROVED_AI` | `ALL` | `DENIED` — только отклонённые ИИ; `MANUAL` — ручная проверка (в т.ч. уже проверенные ИИ); `APPROVED_AI` — только одобренные ИИ и отображаемые в ленте (`isHide=false`); `ALL` — отклонённые + ручная проверка + одобренные ИИ |
| `page` | number | `1, 2, …` | `1` | Пагинация (20 штук на странице) |

### Примеры запросов

```http
GET /admin/moderation/products
GET /admin/moderation/products?filter=DENIED
GET /admin/moderation/products?filter=MANUAL&page=2
GET /admin/moderation/products?filter=APPROVED_AI
```

### Ответ `200`

```json
{
  "items": [
    {
      "id": 1234567,
      "name": "Тонометр Omron",
      "price": 100,
      "images": ["https://s3.../photo.jpg"],
      "moderateState": "DENIDED",
      "moderationRejectionReason": "В описании обнаружен номер телефона",
      "createdAt": "2026-03-23T10:00:00.000Z",
      "updatedAt": "2026-03-23T10:00:15.000Z",
      "category": { "id": 1, "name": "Диагностика" },
      "subCategory": { "id": 3, "name": "Тонометры" },
      "user": {
        "id": 42,
        "fullName": "Иван Иванов",
        "email": "ivan@example.com",
        "phoneNumber": "+79001234567"
      }
    }
  ],
  "total": 45,
  "page": 1,
  "pages": 3
}
```

### Значения `moderateState`

| Значение | Смысл |
|---|---|
| `"DENIDED"` | Отклонено ИИ автоматически — поле `moderationRejectionReason` содержит причину |
| `"MODERATE"` | Новый товар, ещё ожидает первичной проверки ИИ |
| `"AI_REVIEWED"` | Товар уже проверен ИИ и отправлен на ручную проверку; повторно ИИ не обрабатывается |
| `"APPROVED"` + `moderationRejectionReason = "Одобрено ИИ автоматически"` | Товар одобрен ИИ автоматически и отображается в этом же API списке |

---

## 2. Детали одного товара

```
GET /admin/moderation/products/:id
```

### Пример запроса

```http
GET /admin/moderation/products/1234567
```

### Ответ `200`

```json
{
  "id": 1234567,
  "name": "Тонометр Omron",
  "price": 100,
  "description": "Продам тонометр, звоните ...",
  "images": ["https://s3.../photo.jpg"],
  "videoUrl": null,
  "moderateState": "DENIDED",
  "moderationRejectionReason": "В описании обнаружен номер телефона",
  "createdAt": "2026-03-23T10:00:00.000Z",
  "updatedAt": "2026-03-23T10:00:15.000Z",
  "category": { "id": 1, "name": "Диагностика" },
  "subCategory": { "id": 3, "name": "Тонометры" },
  "type": { "id": 7, "name": "Автоматические" },
  "user": {
    "id": 42,
    "fullName": "Иван Иванов",
    "email": "ivan@example.com",
    "phoneNumber": "+79001234567",
    "profileType": "INDIVIDUAL"
  },
  "fieldValues": [
    { "value": "Японский", "field": { "id": 5, "name": "Производитель" } }
  ]
}
```

### Ответ `404` — товар не найден

```json
{ "statusCode": 404, "message": "No Product found" }
```

---

## Пример интеграции (fetch)

```ts
// Все товары на модерации (1-я страница)
const res = await fetch('/admin/moderation/products?filter=ALL&page=1', {
  credentials: 'include', // передаёт cookie session_id
});
const data = await res.json();
// data.items — массив товаров
// data.total — всего записей
// data.page  — текущая страница
// data.pages — всего страниц

// Только отклонённые
const denied = await fetch('/admin/moderation/products?filter=DENIED', {
  credentials: 'include',
});

// Только на ручной проверке
const manual = await fetch('/admin/moderation/products?filter=MANUAL&page=1', {
  credentials: 'include',
});

// Детали конкретного товара
const product = await fetch('/admin/moderation/products/1234567', {
  credentials: 'include',
});
```
