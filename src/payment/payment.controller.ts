import {
  Body,
  Controller,
  Get,
  HttpCode,
  HttpStatus,
  Post,
  Req,
  UseGuards,
} from '@nestjs/common';
import {
  ApiBearerAuth,
  ApiOperation,
  ApiResponse,
  ApiTags,
} from '@nestjs/swagger';

import { SessionAuthGuard } from 'src/auth/guards/session-auth.guard';

import { CheckStatusDto } from './dto/check-status.dto';
import { CreatePaymentDto } from './dto/create-payment.dto';
import { PaymentNotificationDto } from './dto/payment-notification.dto';
import { PaymentService } from './payment.service';

@ApiTags('Платежи')
@Controller('payment')
export class PaymentController {
  constructor(private readonly paymentService: PaymentService) {}

  @Post('create')
  @UseGuards(SessionAuthGuard)
  @ApiBearerAuth()
  @ApiOperation({ summary: 'Создание платежа для пополнения баланса' })
  @ApiResponse({
    status: 201,
    description: 'Платеж успешно создан',
    schema: {
      example: {
        paymentId: '2673412345',
        paymentUrl: 'https://securepay.tinkoff.ru/...',
        orderId: '123-1234567890',
        amount: 1000,
      },
    },
  })
  @ApiResponse({
    status: 400,
    description: 'Ошибка валидации или создания платежа',
  })
  @ApiResponse({ status: 401, description: 'Не авторизован' })
  async createPayment(@Req() req, @Body() dto: CreatePaymentDto) {
    const userId = req.user.id;
    const description = dto.description || 'Пополнение баланса';

    return this.paymentService.createPayment(userId, dto.amount, description);
  }

  @Post('notification')
  @HttpCode(HttpStatus.OK)
  @ApiOperation({
    summary: 'Webhook для уведомлений от Т-Банка о статусе платежа',
  })
  @ApiResponse({
    status: 200,
    description: 'Уведомление успешно обработано',
    schema: {
      example: {
        success: true,
        message: 'Баланс успешно пополнен',
      },
    },
  })
  @ApiResponse({
    status: 400,
    description: 'Неверная подпись уведомления или платеж не найден',
  })
  async handleNotification(@Body() notification: PaymentNotificationDto) {
    return this.paymentService.handleNotification(notification);
  }

  @Get('history')
  @UseGuards(SessionAuthGuard)
  @ApiBearerAuth()
  @ApiOperation({ summary: 'Получение истории платежей пользователя' })
  @ApiResponse({
    status: 200,
    description: 'История платежей',
    schema: {
      example: [
        {
          id: 1,
          orderId: '123-1234567890',
          paymentId: '2673412345',
          userId: 123,
          amount: 1000,
          status: 'COMPLETED',
          paymentUrl: 'https://securepay.tinkoff.ru/...',
          createdAt: '2025-12-26T10:00:00.000Z',
          updatedAt: '2025-12-26T10:05:00.000Z',
        },
      ],
    },
  })
  @ApiResponse({ status: 401, description: 'Не авторизован' })
  async getPaymentHistory(@Req() req) {
    return this.paymentService.getUserPayments(req.user.id);
  }

  @Post('check-status')
  @UseGuards(SessionAuthGuard)
  @ApiBearerAuth()
  @ApiOperation({ summary: 'Проверка статуса платежа в Т-Банке' })
  @ApiResponse({
    status: 200,
    description: 'Статус платежа получен',
    schema: {
      example: {
        status: 'CONFIRMED',
        amount: 1000,
        orderId: '123-1234567890',
      },
    },
  })
  @ApiResponse({ status: 400, description: 'Ошибка проверки статуса' })
  @ApiResponse({ status: 401, description: 'Не авторизован' })
  async checkStatus(@Body() dto: CheckStatusDto) {
    return this.paymentService.checkPaymentStatus(dto.paymentId);
  }
}
