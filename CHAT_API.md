# Chat API Documentation

## Обзор

Система чата включает REST API для получения данных и WebSocket для реального времени.

## REST API Endpoints

### 1. Создать чат

```
POST /chat/start
Authorization: Bearer {token}
Body: { "productId": 1 }
```

### 2. Получить список чатов пользователя

```
GET /chat
Authorization: Bearer {token}
```

### 3. Получить информацию о чате

```
GET /chat/:id
Authorization: Bearer {token}
```

### 4. Получить сообщения чата (с пагинацией)

```
GET /chat/:id/messages?page=1&limit=50
Authorization: Bearer {token}
```

## WebSocket Events

### Подключение

```javascript
const socket = io('ws://localhost:3000/chat', {
  auth: {
    token: 'your-jwt-token',
  },
});
```

### События отправки (клиент → сервер)

#### Присоединиться к чату

```javascript
socket.emit('joinChat', { chatId: 1 });
```

#### Покинуть чат

```javascript
socket.emit('leaveChat', { chatId: 1 });
```

#### Отправить сообщение

```javascript
socket.emit('sendMessage', {
  chatId: 1,
  content: 'Привет!',
});
```

#### Индикатор печатания

```javascript
socket.emit('typing', {
  chatId: 1,
  isTyping: true,
});
```

#### Отметить как прочитанное

```javascript
socket.emit('markAsRead', { chatId: 1 });
```

### События получения (сервер → клиент)

#### Новое сообщение в чате

```javascript
socket.on('newMessage', (data) => {
  console.log('Новое сообщение:', data);
  // data: { id, content, senderId, sender, createdAt, chatId }
});
```

#### Уведомление о новом сообщении

```javascript
socket.on('newChatMessage', (data) => {
  console.log('Новое сообщение в чате:', data);
  // data: { chatId, message, product }
});
```

#### Пользователь печатает

```javascript
socket.on('userTyping', (data) => {
  console.log('Пользователь печатает:', data);
  // data: { chatId, userId, isTyping }
});
```

#### Сообщения прочитаны

```javascript
socket.on('messagesRead', (data) => {
  console.log('Сообщения прочитаны:', data);
  // data: { chatId, readBy }
});
```

## Рекомендации по использованию

1. **Получение данных** - используйте REST API
   - Список чатов
   - История сообщений
   - Информация о чате

2. **Реальное время** - используйте WebSocket
   - Отправка сообщений
   - Получение новых сообщений
   - Индикаторы печатания
   - Статусы прочтения

3. **Безопасность**
   - Все endpoints защищены JWT аутентификацией
   - WebSocket требует токен при подключении
   - Доступ к чатам проверяется на уровне сервиса

4. **Производительность**
   - Сообщения загружаются с пагинацией (по 50 по умолчанию)
   - WebSocket события отправляются только участникам чата
   - Автоматическое подключение к личным комнатам пользователей
