import {
  IsString,
  IsNotEmpty,
  IsNumber,
  IsOptional,
  IsEnum,
  Min,
} from 'class-validator';
import { Transform } from 'class-transformer';
import { ApiProperty } from '@nestjs/swagger';
import { ProductState } from '../enum/product-state.enum';

export class createProductDto {
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
    description: 'Состояние товара',
    enum: ProductState,
    example: ProductState.NEW,
  })
  @IsEnum(ProductState, { message: 'Состояние должно быть NEW или USED' })
  state: ProductState;

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
    description: 'ID категории товара',
    example: 1,
  })
  @Transform(({ value }) => parseInt(value))
  @IsNumber({}, { message: 'Id категории должно быть числом' })
  categoryId: number;

  @ApiProperty({
    description: 'ID подкатегории товара',
    example: 1,
  })
  @Transform(({ value }) => parseInt(value))
  @IsNumber({}, { message: 'Id подкатегории должно быть числом' })
  subcategoryId: number;

  @ApiProperty({
    description: 'ID типа подкатегории (например: Измерительные приборы)',
    example: 1,
    required: false,
  })
  @IsOptional()
  @Transform(({ value }) => (value ? parseInt(value) : undefined))
  @IsNumber({}, { message: 'Id типа должно быть числом' })
  typeId?: number;

  @ApiProperty({
    description: 'Дополнительные поля товара (fieldId: value)',
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
