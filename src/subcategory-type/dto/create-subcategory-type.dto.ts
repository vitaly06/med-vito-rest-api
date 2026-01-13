import { ApiProperty } from '@nestjs/swagger';

import { Type } from 'class-transformer';
import {
  IsInt,
  IsNotEmpty,
  IsNumber,
  IsPositive,
  IsString,
} from 'class-validator';

export class CreateSubcategoryTypeDto {
  @ApiProperty({
    example: 'Мужская',
  })
  @IsString({ message: 'Название типа подкатегории должно быть строкой' })
  @IsNotEmpty({
    message: 'Название типа подкатегории обязательно для заполнения',
  })
  name: string;
  @ApiProperty({
    example: 1,
  })
  @IsNotEmpty({ message: 'Id подкатегории обязателен для заполнения' })
  @IsNumber({}, { message: 'Id подкатегории должен быть числом' })
  @IsPositive({ message: 'Id подкатегории должен быть положительным числом' })
  @IsInt({ message: 'Id подкатегории должен быть целым числом' })
  @Type(() => Number)
  subcategoryId: number;
}
