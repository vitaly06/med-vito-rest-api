import { ApiProperty } from '@nestjs/swagger';
import { Type } from 'class-transformer';
import {
  IsInt,
  IsNotEmpty,
  IsNumber,
  IsPositive,
  IsString,
} from 'class-validator';

export class CreateSubcategoryDto {
  @ApiProperty({
    example: 'Резина и диски',
  })
  @IsString({ message: 'Название подкатегории должно быть строкой' })
  @IsNotEmpty({ message: 'Название подкатегории обязательно для заполнения' })
  name: string;

  @ApiProperty({
    example: 1,
  })
  @IsNotEmpty({ message: 'Id категории обязателен для заполнения' })
  @IsNumber({}, { message: 'Id категории должен быть числом' })
  @IsPositive({ message: 'Id категории должен быть положительным числом' })
  @IsInt({ message: 'Id категории должен быть целым числом' })
  @Type(() => Number)
  categoryId: number;
}
