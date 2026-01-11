import {
  Body,
  Controller,
  Get,
  Param,
  ParseIntPipe,
  Post,
  Query,
  Req,
  UseGuards,
} from '@nestjs/common';
import {
  ApiBearerAuth,
  ApiOperation,
  ApiParam,
  ApiQuery,
  ApiResponse,
  ApiTags,
} from '@nestjs/swagger';

import { Request } from 'express';
import { SessionAuthGuard } from 'src/auth/guards/session-auth.guard';

import { ChatService } from './chat.service';
import { GetMessagesDto, StartChatDto } from './dto/index';

@ApiTags('Chat')
@Controller('chat')
@UseGuards(SessionAuthGuard)
@ApiBearerAuth()
export class ChatController {
  constructor(private readonly chatService: ChatService) {}

  @Post('start')
  @ApiOperation({
    summary: 'Начать чат по товару',
    description:
      'Создает новый чат между покупателем и продавцом товара или возвращает существующий',
  })
  @ApiResponse({
    status: 201,
    description: 'Чат успешно создан или найден',
  })
  @ApiResponse({
    status: 400,
    description: 'Нельзя писать самому себе или товар не найден',
  })
  async startChat(
    @Body() dto: StartChatDto,
    @Req() req: Request & { user: any },
  ) {
    return await this.chatService.startChat(dto.productId, req.user.id);
  }

  @Get()
  @ApiOperation({
    summary: 'Получить список чатов пользователя',
    description:
      'Возвращает все чаты пользователя с информацией о собеседниках и последних сообщениях',
  })
  @ApiResponse({
    status: 200,
    description: 'Список чатов успешно получен',
  })
  async getUserChats(@Req() req: Request & { user: any }) {
    return await this.chatService.getUserChats(req.user.id);
  }

  @Get(':id')
  @ApiOperation({
    summary: 'Получить информацию о чате',
    description: 'Возвращает подробную информацию о чате, товаре и собеседнике',
  })
  @ApiParam({ name: 'id', description: 'ID чата', type: Number })
  @ApiResponse({
    status: 200,
    description: 'Информация о чате успешно получена',
  })
  @ApiResponse({
    status: 403,
    description: 'Нет доступа к чату',
  })
  @ApiResponse({
    status: 404,
    description: 'Чат не найден',
  })
  async getChatInfo(
    @Param('id', ParseIntPipe) chatId: number,
    @Req() req: Request & { user: any },
  ) {
    return await this.chatService.getChatInfo(chatId, req.user.id);
  }

  @Get(':id/messages')
  @ApiOperation({
    summary: 'Получить сообщения чата',
    description: 'Возвращает сообщения чата с пагинацией',
  })
  @ApiParam({ name: 'id', description: 'ID чата', type: Number })
  @ApiQuery({ name: 'page', description: 'Номер страницы', required: false })
  @ApiQuery({
    name: 'limit',
    description: 'Количество сообщений на странице',
    required: false,
  })
  @ApiResponse({
    status: 200,
    description: 'Сообщения успешно получены',
  })
  async getChatMessages(
    @Param('id', ParseIntPipe) chatId: number,
    @Query() query: GetMessagesDto,
    @Req() req: Request & { user: any },
  ) {
    return await this.chatService.getChatMessages(
      chatId,
      req.user.id,
      query.page,
      query.limit,
    );
  }
}
