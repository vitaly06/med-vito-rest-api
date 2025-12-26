import { ApiProperty } from '@nestjs/swagger';

export class PaymentNotificationDto {
  @ApiProperty({
    description: 'ID терминала',
    example: '1766153689307DEMO',
  })
  TerminalKey: string;

  @ApiProperty({
    description: 'ID заказа в системе',
    example: '123-1234567890',
  })
  OrderId: string;

  @ApiProperty({
    description: 'Успешность операции',
    example: true,
  })
  Success: boolean;

  @ApiProperty({
    description: 'Статус платежа',
    example: 'CONFIRMED',
    enum: [
      'NEW',
      'FORM_SHOWED',
      'AUTHORIZING',
      'AUTHORIZED',
      'CONFIRMING',
      'CONFIRMED',
      'REVERSING',
      'REVERSED',
      'REFUNDING',
      'REFUNDED',
      'REJECTED',
      'CANCELED',
    ],
  })
  Status: string;

  @ApiProperty({
    description: 'ID платежа в системе Т-Банк',
    example: '2673412345',
  })
  PaymentId: string;

  @ApiProperty({
    description: 'Код ошибки (если есть)',
    example: '0',
    required: false,
  })
  ErrorCode?: string;

  @ApiProperty({
    description: 'Сумма платежа в копейках',
    example: 100000,
  })
  Amount: number;

  @ApiProperty({
    description: 'Токен для проверки подлинности уведомления',
    example: 'abc123...',
  })
  Token: string;

  @ApiProperty({
    description: 'Номер карты (маска)',
    example: '430000******0777',
    required: false,
  })
  Pan?: string;
}
