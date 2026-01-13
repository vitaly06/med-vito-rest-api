import { ApiProperty } from '@nestjs/swagger';

import { IsNumber, IsOptional, IsString, Min } from 'class-validator';

export class CreatePaymentDto {
  @ApiProperty({
    description: 'Сумма пополнения в рублях',
    example: 1000,
    minimum: 1,
  })
  @IsNumber()
  @Min(1, { message: 'Минимальная сумма пополнения 1 рубль' })
  amount: number;

  @ApiProperty({
    description: 'Описание платежа',
    example: 'Пополнение баланса',
    required: false,
  })
  @IsString()
  @IsOptional()
  description?: string;
}
