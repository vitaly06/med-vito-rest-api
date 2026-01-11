import { ApiProperty } from '@nestjs/swagger';

import { Type } from 'class-transformer';
import {
  IsInt,
  IsNotEmpty,
  IsNumber,
  IsPositive,
  IsString,
} from 'class-validator';

export class CreateTypeFieldDto {
  @ApiProperty({
    example: 'Размер',
  })
  @IsString({ message: 'Название характеристики должно быть строкой' })
  @IsNotEmpty({
    message: 'Название характеристики обязательно для заполнения',
  })
  name: string;
  @ApiProperty({
    example: 1,
  })
  @IsNotEmpty({ message: 'Id типа подкатегории обязателен для заполнения' })
  @IsNumber({}, { message: 'Id типа подкатегории должен быть числом' })
  @IsPositive({
    message: 'Id типа подкатегории должен быть положительным числом',
  })
  @IsInt({ message: 'Id типа подкатегории должен быть целым числом' })
  @Type(() => Number)
  typeId: number;
}
