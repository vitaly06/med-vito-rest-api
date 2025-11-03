import {
  WebSocketGateway,
  SubscribeMessage,
  MessageBody,
  ConnectedSocket,
  OnGatewayInit,
  OnGatewayConnection,
  OnGatewayDisconnect,
  WebSocketServer,
} from '@nestjs/websockets';
import { UseGuards } from '@nestjs/common';
import { Server, Socket } from 'socket.io';
import { ChatService } from './chat.service';
import { WsJwtAuthGuard } from 'src/auth/guards/ws-jwt-auth.guard';

interface AuthenticatedSocket extends Socket {
  userId?: number;
  user?: any;
}

@WebSocketGateway({
  cors: {
    origin: true,
    credentials: true,
  },
  namespace: '/chat',
})
export class ChatGateway
  implements OnGatewayInit, OnGatewayConnection, OnGatewayDisconnect
{
  @WebSocketServer()
  server: Server;

  constructor(private readonly chatService: ChatService) {}

  afterInit(server: Server) {
    console.log('WebSocket Server Initialized');
  }

  async handleConnection(client: AuthenticatedSocket) {
    console.log('WebSocket client connected:', client.id);
    console.log('Headers:', client.handshake.headers);
    console.log('Auth:', client.handshake.auth);
    console.log('Query:', client.handshake.query);

    // Получаем токен из куков (приоритет) или из других источников
    let token: string | null = null;

    // 1. Пытаемся получить из куков (как в основной стратегии)
    const cookies = client.handshake.headers.cookie;
    if (cookies) {
      console.log('Found cookies:', cookies);
      const accessTokenMatch = cookies.match(/access_token=([^;]+)/);
      if (accessTokenMatch) {
        token = accessTokenMatch[1];
        console.log('Found access_token in cookies');
      }
    }

    // 2. Fallback к другим методам, если куки недоступны
    if (!token) {
      token =
        (client.handshake.auth?.token as string) ||
        client.handshake.headers?.authorization?.replace('Bearer ', '') ||
        (client.handshake.query?.token as string);

      if (token) {
        console.log('Found token in auth/headers/query');
      }
    }

    if (!token) {
      console.log('No token provided. For development, you can:');
      console.log(
        '1. Open http://localhost:8080/test-websocket.html instead of file://',
      );
      console.log('2. Or provide token in query: ?token=YOUR_JWT_TOKEN');
      console.log('Origin:', client.handshake.headers.origin);

      // Для разработки временно разрешим подключение без токена, если origin = null
      const origin = client.handshake.headers.origin;
      if (origin === null || origin === 'null') {
        console.log(
          'Development mode: allowing connection without token for null origin',
        );
        client.userId = 1; // Заглушка для разработки
        client.user = { id: 1 };
        client.join(`user_${client.userId}`);
        console.log(
          `User ${client.userId} connected to chat (development mode)`,
        );
        return;
      }

      client.disconnect();
      return;
    }

    try {
      // Проверяем токен
      const jwt = require('jsonwebtoken');
      const payload = jwt.verify(
        token,
        process.env.JWT_SECRET || 'fallback-secret',
      );

      client.userId = payload.sub;
      client.user = { id: payload.sub };

      // Подключаем пользователя к его личной комнате
      client.join(`user_${client.userId}`);

      console.log(`User ${client.userId} connected to chat with valid token`);
    } catch (error) {
      console.log('Authentication failed:', error.message);
      client.disconnect();
    }
  }

  handleDisconnect(client: AuthenticatedSocket) {
    console.log(`User ${client.userId} disconnected from chat`);
  }

  // Подключиться к комнате чата
  @SubscribeMessage('joinChat')
  @UseGuards(WsJwtAuthGuard)
  async handleJoinChat(
    @MessageBody() data: { chatId: number },
    @ConnectedSocket() client: AuthenticatedSocket,
  ) {
    try {
      console.log(
        `Attempting to join chat ${data.chatId} for user ${client.userId}`,
      );

      // Проверяем доступ к чату
      await this.chatService.getChatInfo(data.chatId, client.userId!);

      // Подключаем к комнате чата
      client.join(`chat_${data.chatId}`);

      console.log(`User ${client.userId} joined chat ${data.chatId}`);

      return { success: true, message: 'Подключен к чату' };
    } catch (error) {
      console.log('joinChat error:', error.message);
      return { success: false, message: error.message };
    }
  }

  // Отправить сообщение через WebSocket
  @SubscribeMessage('sendMessage')
  @UseGuards(WsJwtAuthGuard)
  async handleMessage(
    @MessageBody() data: { chatId: number; content: string },
    @ConnectedSocket() client: AuthenticatedSocket,
  ) {
    try {
      // Отправляем сообщение через сервис
      const message = await this.chatService.sendMessage(
        data.chatId,
        client.userId!,
        data.content,
      );

      // Отправляем сообщение всем участникам чата
      this.server.to(`chat_${data.chatId}`).emit('newMessage', {
        id: message.id,
        content: message.content,
        senderId: message.senderId,
        sender: message.sender,
        createdAt: message.createdAt,
        timeString: message.timeString,
        chatId: data.chatId,
      });

      // Получаем информацию о чате для уведомления
      const chatInfo = await this.chatService.getChatInfo(
        data.chatId,
        client.userId!,
      );

      // Уведомляем собеседника о новом сообщении
      const recipientId = chatInfo.companion.id;
      this.server.to(`user_${recipientId}`).emit('newChatMessage', {
        chatId: data.chatId,
        message: message,
        product: chatInfo.product,
        companion: chatInfo.companion,
      });

      return { success: true, message: 'Сообщение отправлено' };
    } catch (error) {
      return { success: false, message: error.message };
    }
  }

  // Отметить сообщения как прочитанные
  @SubscribeMessage('markAsRead')
  @UseGuards(WsJwtAuthGuard)
  async handleMarkAsRead(
    @MessageBody() data: { chatId: number },
    @ConnectedSocket() client: AuthenticatedSocket,
  ) {
    try {
      await this.chatService.markMessagesAsRead(data.chatId, client.userId!);

      // Уведомляем о прочтении
      this.server.to(`chat_${data.chatId}`).emit('messagesRead', {
        chatId: data.chatId,
        readBy: client.userId,
      });

      return { success: true };
    } catch (error) {
      return { success: false, message: error.message };
    }
  }
}
