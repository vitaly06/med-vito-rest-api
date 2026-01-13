import { ApiProperty } from '@nestjs/swagger';

import { Transform } from 'class-transformer';
import { IsNumber, IsOptional } from 'class-validator';

export class StartChatDto {
  @ApiProperty({
    description: 'ID товара',
    example: 1,
  })
  @Transform(({ value }) => parseInt(value))
  @IsNumber({}, { message: 'ID товара должно быть числом' })
  productId: number;
}

export class GetMessagesDto {
  @ApiProperty({
    description: 'Номер страницы',
    example: 1,
    required: false,
  })
  @IsOptional()
  @Transform(({ value }) => parseInt(value))
  @IsNumber({}, { message: 'Номер страницы должен быть числом' })
  page?: number = 1;

  @ApiProperty({
    description: 'Количество сообщений на странице',
    example: 50,
    required: false,
  })
  @IsOptional()
  @Transform(({ value }) => parseInt(value))
  @IsNumber({}, { message: 'Лимит должен быть числом' })
  limit?: number = 50;
}
