import {
  IsString,
  IsNumber,
  IsOptional,
  IsEnum,
  Min,
  IsUrl,
} from 'class-validator';
import { Transform } from 'class-transformer';
import { ApiProperty } from '@nestjs/swagger';
import { ProductState } from '../enum/product-state.enum';

export class UpdateProductDto {
  @ApiProperty({
    description: 'Название продукта',
    example: 'iPhone 15 Pro',
    required: false,
  })
  @IsOptional()
  @IsString({ message: 'Название должно быть строкой' })
  name?: string;

  @ApiProperty({
    description: 'Цена в рублях',
    example: 120000,
    minimum: 1,
    required: false,
  })
  @IsOptional()
  @Transform(({ value }) => parseInt(value))
  @IsNumber({}, { message: 'Цена должна быть числом' })
  @Min(1, { message: 'Цена должна быть больше 0' })
  price?: number;

  @ApiProperty({
    description: 'Состояние товара',
    enum: ProductState,
    example: ProductState.NEW,
    required: false,
  })
  @IsOptional()
  @IsEnum(ProductState, { message: 'Состояние должно быть NEW или USED' })
  state?: ProductState;

  @ApiProperty({
    description: 'Описание товара',
    example: 'Новый iPhone 15 Pro в отличном состоянии. Полная комплектация.',
    required: false,
  })
  @IsOptional()
  @IsString({ message: 'Описание должно быть строкой' })
  description?: string;

  @ApiProperty({
    description: 'Адрес продавца',
    example: 'г. Москва, ул. Тверская, д. 1',
    required: false,
  })
  @IsOptional()
  @IsString({ message: 'Адрес должен быть строкой' })
  address?: string;

  @ApiProperty({
    description: 'Ссылка на видео товара (YouTube, VK и т.д.)',
    example: 'https://www.youtube.com/watch?v=8u6xaEYGLa0',
    required: false,
  })
  @IsOptional()
  @IsUrl({}, { message: 'Неверный формат url' })
  videoUrl?: string;

  @ApiProperty({
    description:
      'Дополнительные поля товара для добавления или обновления (fieldId: value)',
    example: { '1': 'Тонометр', '2': 'Omron' },
    required: false,
  })
  @IsOptional()
  @Transform(({ value }) => {
    if (typeof value === 'string') {
      try {
        return JSON.parse(value);
      } catch {
        return value;
      }
    }
    return value;
  })
  fieldValues?: Record<string, string>;
}
