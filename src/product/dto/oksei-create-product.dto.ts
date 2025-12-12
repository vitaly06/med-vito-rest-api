import {
  IsString,
  IsNotEmpty,
  IsNumber,
  IsOptional,
  IsEnum,
  Min,
  IsUrl,
} from 'class-validator';
import { Transform } from 'class-transformer';
import { ApiProperty } from '@nestjs/swagger';
import { ProductState } from '../enum/product-state.enum';

export class createOkseiProductDto {
  @ApiProperty({
    description: 'Название продукта',
    example: 'iPhone 15 Pro',
  })
  @IsString({ message: 'Название должно быть строкой' })
  @IsNotEmpty({ message: 'Название обязательно для заполнения' })
  name: string;

  @ApiProperty({
    description: 'Цена в рублях',
    example: 120000,
    minimum: 1,
  })
  @Transform(({ value }) => parseInt(value))
  @IsNumber({}, { message: 'Цена должна быть числом' })
  @Min(1, { message: 'Цена должна быть больше 0' })
  price: number;

  @ApiProperty({
    description: 'Описание товара',
    example: 'Новый iPhone 15 Pro в отличном состоянии. Полная комплектация.',
    required: false,
  })
  @IsOptional()
  @IsString({ message: 'Описание должно быть строкой' })
  description?: string;
}
