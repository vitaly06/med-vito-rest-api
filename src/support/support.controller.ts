import {
  Body,
  Controller,
  ForbiddenException,
  Get,
  Param,
  ParseIntPipe,
  Post,
  Put,
  Query,
  Request,
  UseGuards,
} from '@nestjs/common';

import { SessionAuthGuard } from '../auth/guards/session-auth.guard';
import {
  CreateTicketDto,
  GetTicketsQueryDto,
  SendSupportMessageDto,
  UpdateTicketDto,
} from './dto';
import { SupportService } from './support.service';

@Controller('support')
@UseGuards(SessionAuthGuard)
export class SupportController {
  constructor(private readonly supportService: SupportService) {}

  // Создание нового тикета
  @Post('tickets')
  async createTicket(@Request() req, @Body() createTicketDto: CreateTicketDto) {
    return this.supportService.createTicket(req.user.id, createTicketDto);
  }

  // Получение тикетов пользователя
  @Get('tickets/my')
  async getMyTickets(@Request() req, @Query() query: GetTicketsQueryDto) {
    return this.supportService.getUserTickets(req.user.id, query);
  }

  // Получение всех тикетов (только для модераторов и админов)
  @Get('tickets/all')
  async getAllTickets(@Request() req, @Query() query: GetTicketsQueryDto) {
    const userRole = req.user.role?.name;
    if (userRole !== 'moderator' && userRole !== 'admin') {
      throw new ForbiddenException(
        'Доступ запрещен. Требуется роль модератора или администратора',
      );
    }
    return this.supportService.getAllTickets(query);
  }

  // Получение конкретного тикета с сообщениями
  @Get('tickets/:id')
  async getTicket(@Request() req, @Param('id', ParseIntPipe) ticketId: number) {
    const userRole = req.user.role?.name;
    const isModerator = userRole === 'moderator' || userRole === 'admin';

    return this.supportService.getTicket(ticketId, req.user.id, isModerator);
  }

  // Отправка сообщения в тикет
  @Post('tickets/:id/messages')
  async sendMessage(
    @Request() req,
    @Param('id', ParseIntPipe) ticketId: number,
    @Body() sendMessageDto: SendSupportMessageDto,
  ) {
    const userRole = req.user.role?.name;
    const isModerator = userRole === 'moderator' || userRole === 'admin';

    return this.supportService.sendMessage(
      ticketId,
      req.user.id,
      sendMessageDto,
      isModerator,
    );
  }

  // Обновление тикета (только для модераторов и админов)
  @Put('tickets/:id')
  async updateTicket(
    @Request() req,
    @Param('id', ParseIntPipe) ticketId: number,
    @Body() updateTicketDto: UpdateTicketDto,
  ) {
    const userRole = req.user.role?.name;
    if (userRole !== 'moderator' && userRole !== 'admin') {
      throw new ForbiddenException(
        'Доступ запрещен. Требуется роль модератора или администратора',
      );
    }

    return this.supportService.updateTicket(
      ticketId,
      req.user.id,
      updateTicketDto,
    );
  }

  // Назначение тикета модератору
  @Put('tickets/:id/assign')
  async assignTicket(
    @Request() req,
    @Param('id', ParseIntPipe) ticketId: number,
  ) {
    const userRole = req.user.role?.name;
    if (userRole !== 'moderator' && userRole !== 'admin') {
      throw new ForbiddenException(
        'Доступ запрещен. Требуется роль модератора или администратора',
      );
    }

    return this.supportService.assignTicket(ticketId, req.user.id);
  }

  // Получение статистики тикетов (только для админов)
  @Get('stats')
  async getStats(@Request() req) {
    const userRole = req.user.role?.name;
    if (userRole !== 'admin') {
      throw new ForbiddenException(
        'Доступ запрещен. Требуется роль администратора',
      );
    }

    return this.supportService.getTicketStats();
  }
}
