import { UseGuards } from '@nestjs/common';
import {
  ConnectedSocket,
  MessageBody,
  OnGatewayConnection,
  OnGatewayDisconnect,
  OnGatewayInit,
  SubscribeMessage,
  WebSocketGateway,
  WebSocketServer,
} from '@nestjs/websockets';

import { Server, Socket } from 'socket.io';

import { WsSessionAuthGuard } from '@/auth/guards/ws-session-auth.guard';

import { ChatService } from './chat.service';

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

    // Получаем session_id из cookies
    const cookies = client.handshake.headers.cookie;
    let sessionId: string | null = null;

    if (cookies) {
      console.log('Found cookies');
      const sessionIdMatch = cookies.match(/session_id=([^;]+)/);
      if (sessionIdMatch) {
        sessionId = sessionIdMatch[1];
        console.log('Found session_id in cookies');
      }
    }

    if (!sessionId) {
      console.log('No session_id provided, disconnecting');
      client.disconnect();
      return;
    }

    try {
      // Проверяем сессию через AuthService
      const sessionData = await this.chatService.getSessionData(sessionId);

      if (!sessionData) {
        console.log('Invalid session, disconnecting');
        client.disconnect();
        return;
      }

      // Получаем пользователя
      const user = await this.chatService.getUserById(sessionData.userId);

      if (!user) {
        console.log('User not found, disconnecting');
        client.disconnect();
        return;
      }

      client.userId = user.id;
      client.user = user;

      // Подключаем пользователя к его личной комнате
      client.join(`user_${client.userId}`);

      console.log(`User ${client.userId} connected to chat with valid session`);
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
  @UseGuards(WsSessionAuthGuard)
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
  @UseGuards(WsSessionAuthGuard)
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
  @UseGuards(WsSessionAuthGuard)
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
