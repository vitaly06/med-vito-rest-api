import { ApiProperty } from '@nestjs/swagger';
import { IsString, IsNotEmpty } from 'class-validator';

export class CheckStatusDto {
  @ApiProperty({
    description: 'ID платежа в системе Т-Банк',
    example: '2673412345',
  })
  @IsString()
  @IsNotEmpty()
  paymentId: string;
}
