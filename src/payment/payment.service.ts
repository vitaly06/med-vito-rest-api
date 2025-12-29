import { Injectable, BadRequestException } from '@nestjs/common';
import { ConfigService } from '@nestjs/config';
import { PrismaService } from 'src/prisma/prisma.service';
import * as crypto from 'crypto';

@Injectable()
export class PaymentService {
  private readonly terminalKey: string;
  private readonly secretKey: string;
  private readonly apiUrl = 'https://securepay.tinkoff.ru/v2';

  constructor(
    private readonly configService: ConfigService,
    private readonly prisma: PrismaService,
  ) {
    this.terminalKey = this.configService.get<string>(
      'TINKOFF_TERMINAL_KEY',
      'not found',
    );
    this.secretKey = this.configService.get<string>(
      'TINKOFF_SECRET_KEY',
      'not found',
    );
  }

  private generateToken(params: Record<string, any>): string {
    const values = { ...params, Password: this.secretKey };
    const sortedKeys = Object.keys(values).sort();
    const concatenated = sortedKeys.map((key) => values[key]).join('');
    console.log('Token generation:', { sortedKeys, concatenated });
    return crypto.createHash('sha256').update(concatenated).digest('hex');
  }

  async createPayment(userId: number, amount: number, description: string) {
    if (amount < 1) {
      throw new BadRequestException('Минимальная сумма пополнения 1 рубль');
    }

    const orderId = `${userId}-${Date.now()}`;
    const amountInKopecks = Math.round(amount * 100);

    // email пользователя для чека
    const user = await this.prisma.user.findUnique({
      where: { id: userId },
      select: { email: true },
    });

    const params = {
      TerminalKey: this.terminalKey,
      Amount: amountInKopecks,
      OrderId: orderId,
      Description: description,
    };

    const token = this.generateToken(params);

    const requestBody = { ...params, Token: token };
    console.log('Request to Tinkoff:', JSON.stringify(requestBody, null, 2));

    try {
      const response = await fetch(`${this.apiUrl}/Init`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(requestBody),
      });

      const data = await response.json();
      console.log('Response from Tinkoff:', JSON.stringify(data, null, 2));

      if (!data.Success) {
        throw new BadRequestException(
          `Ошибка создания платежа: ${data.Message || data.Details}`,
        );
      }

      // Сохраняем платеж в БД
      await this.prisma.payment.create({
        data: {
          orderId,
          paymentId: data.PaymentId.toString(),
          userId,
          amount,
          status: 'PENDING',
          paymentUrl: data.PaymentURL,
        },
      });

      return {
        paymentId: data.PaymentId.toString(),
        paymentUrl: data.PaymentURL,
        orderId,
        amount,
      };
    } catch (error) {
      throw new BadRequestException(
        `Ошибка при создании платежа: ${error.message}`,
      );
    }
  }

  async handleNotification(notification: any) {
    const { Token, ...params } = notification;
    const expectedToken = this.generateToken(params);

    // Проверка подписи
    if (Token !== expectedToken) {
      throw new BadRequestException('Неверная подпись уведомления');
    }

    const paymentId = notification.PaymentId.toString();

    const payment = await this.prisma.payment.findUnique({
      where: { paymentId },
      include: { user: true },
    });

    if (!payment) {
      throw new BadRequestException('Платеж не найден');
    }

    // Обновляем статус платежа
    if (notification.Status === 'CONFIRMED') {
      await this.prisma.$transaction([
        this.prisma.payment.update({
          where: { paymentId },
          data: { status: 'COMPLETED' },
        }),
        this.prisma.user.update({
          where: { id: payment.userId },
          data: {
            balance: { increment: payment.amount },
          },
        }),
        this.prisma.log.create({
          data: {
            userId: payment.userId,
            action: `Пополнение баланса: id: ${payment.userId}; email: ${payment.user.email};\nсумма: ${payment.amount}; баланс: ${payment.user.balance}; бонусный баланс: ${payment.user.bonusBalance}`,
          },
        }),
      ]);

      return { success: true, message: 'Баланс успешно пополнен' };
    } else if (
      notification.Status === 'REJECTED' ||
      notification.Status === 'CANCELED'
    ) {
      await this.prisma.payment.update({
        where: { paymentId },
        data: { status: 'FAILED' },
      });
    }

    return { success: true };
  }

  async getUserPayments(userId: number) {
    return this.prisma.payment.findMany({
      where: { userId },
      orderBy: { createdAt: 'desc' },
      take: 50,
    });
  }

  async checkPaymentStatus(paymentId: string) {
    const params = {
      TerminalKey: this.terminalKey,
      PaymentId: paymentId,
    };

    const token = this.generateToken(params);

    try {
      const response = await fetch(`${this.apiUrl}/GetState`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ ...params, Token: token }),
      });

      const data = await response.json();

      if (!data.Success) {
        throw new BadRequestException('Ошибка проверки статуса платежа');
      }

      return {
        status: data.Status,
        amount: data.Amount / 100,
        orderId: data.OrderId,
      };
    } catch (error) {
      throw new BadRequestException(
        `Ошибка при проверке статуса: ${error.message}`,
      );
    }
  }
}
