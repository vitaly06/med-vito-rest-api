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

import { WsSessionAuthGuard } from '../auth/guards/ws-session-auth.guard';
import { SendSupportMessageDto } from './dto';
import { SupportService } from './support.service';

interface AuthenticatedSocket extends Socket {
  userId?: number;
  user?: any;
}

@WebSocketGateway({
  cors: {
    origin: true,
    credentials: true,
  },
  namespace: '/support',
})
export class SupportGateway
  implements OnGatewayInit, OnGatewayConnection, OnGatewayDisconnect
{
  @WebSocketServer()
  server: Server;

  constructor(private readonly supportService: SupportService) {}

  afterInit(server: Server) {
    console.log('Support WebSocket Server Initialized');
  }

  async handleConnection(client: AuthenticatedSocket) {
    console.log('Support WebSocket client connected:', client.id);

    // Получаем session_id из cookies
    const cookies = client.handshake.headers.cookie;
    let sessionId: string | null = null;

    if (cookies) {
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
      // Проверяем сессию через SupportService
      const sessionData = await this.supportService.getSessionData(sessionId);

      if (!sessionData) {
        console.log('Invalid session, disconnecting');
        client.disconnect();
        return;
      }

      // Получаем пользователя
      const user = await this.supportService.getUserById(sessionData.userId);

      if (!user) {
        console.log('User not found, disconnecting');
        client.disconnect();
        return;
      }

      client.userId = user.id;
      client.user = user;

      console.log(
        `Support: User ${client.userId} connected with valid session`,
      );
    } catch (error) {
      console.log('Support authentication failed:', error.message);
      client.disconnect();
    }
  }

  handleDisconnect(client: AuthenticatedSocket) {
    console.log('Support WebSocket client disconnected:', client.id);
  }

  // Присоединение к комнате тикета
  @SubscribeMessage('joinTicket')
  async handleJoinTicket(
    @ConnectedSocket() client: AuthenticatedSocket,
    @MessageBody() data: { ticketId: number },
  ) {
    try {
      const userRole = client.user?.role?.name;
      const isModerator = userRole === 'moderator' || userRole === 'admin';

      // Проверяем доступ к тикету
      const ticket = await this.supportService.getTicket(
        data.ticketId,
        client.userId!,
        isModerator,
      );

      if (ticket) {
        const roomName = `support_ticket_${data.ticketId}`;
        await client.join(roomName);

        client.emit('joinedTicket', {
          ticketId: data.ticketId,
          message: 'Присоединились к тикету',
        });

        console.log(
          `Client ${client.id} joined support ticket room: ${roomName}`,
        );
      }
    } catch (error) {
      client.emit('error', {
        message: error.message || 'Ошибка при присоединении к тикету',
      });
    }
  }

  // Покидание комнаты тикета
  @SubscribeMessage('leaveTicket')
  async handleLeaveTicket(
    @ConnectedSocket() client: AuthenticatedSocket,
    @MessageBody() data: { ticketId: number },
  ) {
    const roomName = `support_ticket_${data.ticketId}`;
    await client.leave(roomName);

    client.emit('leftTicket', {
      ticketId: data.ticketId,
      message: 'Покинули тикет',
    });

    console.log(`Client ${client.id} left support ticket room: ${roomName}`);
  }

  // Отправка сообщения в тикет
  @UseGuards(WsSessionAuthGuard)
  @SubscribeMessage('sendSupportMessage')
  async handleSendMessage(
    @ConnectedSocket() client: AuthenticatedSocket,
    @MessageBody() data: { ticketId: number; message: SendSupportMessageDto },
  ) {
    try {
      const userRole = client.user?.role?.name;
      const isModerator = userRole === 'moderator' || userRole === 'admin';

      // Отправляем сообщение через сервис
      const message = await this.supportService.sendMessage(
        data.ticketId,
        client.userId!,
        data.message,
        isModerator,
      );

      const roomName = `support_ticket_${data.ticketId}`;

      // Отправляем сообщение всем участникам тикета
      this.server.to(roomName).emit('newSupportMessage', {
        ticketId: data.ticketId,
        message,
      });

      console.log(`New support message sent to room: ${roomName}`);

      // Если это модератор отвечает, уведомляем пользователя
      if (isModerator) {
        // Можно добавить дополнительные уведомления
        this.server.to(roomName).emit('moderatorResponse', {
          ticketId: data.ticketId,
          message: 'Модератор ответил на ваш запрос',
        });
      }
    } catch (error) {
      client.emit('error', {
        message: error.message || 'Ошибка при отправке сообщения',
      });
    }
  }

  // Уведомление о том, что модератор печатает
  @SubscribeMessage('supportTyping')
  async handleTyping(
    @ConnectedSocket() client: AuthenticatedSocket,
    @MessageBody() data: { ticketId: number; isTyping: boolean },
  ) {
    const roomName = `support_ticket_${data.ticketId}`;

    // Отправляем уведомление другим участникам тикета (кроме отправителя)
    client.to(roomName).emit('supportUserTyping', {
      ticketId: data.ticketId,
      userId: client.userId,
      userName: client.user?.fullName,
      isTyping: data.isTyping,
    });
  }

  // Уведомление об обновлении статуса тикета
  @SubscribeMessage('updateTicketStatus')
  async handleUpdateTicketStatus(
    @ConnectedSocket() client: AuthenticatedSocket,
    @MessageBody() data: { ticketId: number; status: string },
  ) {
    try {
      const userRole = client.user?.role?.name;
      if (userRole !== 'moderator' && userRole !== 'admin') {
        throw new Error('Нет прав для изменения статуса тикета');
      }

      const roomName = `support_ticket_${data.ticketId}`;

      // Уведомляем всех участников об изменении статуса
      this.server.to(roomName).emit('ticketStatusUpdated', {
        ticketId: data.ticketId,
        status: data.status,
        updatedBy: client.user?.fullName,
      });

      console.log(`Ticket ${data.ticketId} status updated to: ${data.status}`);
    } catch (error) {
      client.emit('error', {
        message: error.message || 'Ошибка при обновлении статуса тикета',
      });
    }
  }

  // Метод для отправки уведомлений о новых тикетах модераторам
  public notifyModeratorsNewTicket(ticket: any) {
    this.server.emit('newSupportTicket', {
      ticket,
      message: 'Создан новый тикет поддержки',
    });
  }

  // Метод для уведомления о назначении тикета модератору
  public notifyTicketAssigned(ticketId: number, moderatorId: number) {
    const roomName = `support_ticket_${ticketId}`;

    this.server.to(roomName).emit('ticketAssigned', {
      ticketId,
      moderatorId,
      message: 'Тикет назначен модератору',
    });
  }
}
