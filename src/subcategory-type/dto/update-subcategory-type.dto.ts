import { ApiProperty } from '@nestjs/swagger';

import { Type } from 'class-transformer';
import { IsInt, IsNumber, IsPositive, IsString } from 'class-validator';

export class UpdateSubcategoryTypeDto {
  @ApiProperty({
    example: 'Мужская',
  })
  @IsString({ message: 'Название типа подкатегории должно быть строкой' })
  name?: string;
  @ApiProperty({
    example: 1,
  })
  @IsNumber({}, { message: 'Id подкатегории должен быть числом' })
  @IsPositive({ message: 'Id подкатегории должен быть положительным числом' })
  @IsInt({ message: 'Id подкатегории должен быть целым числом' })
  @Type(() => Number)
  subcategoryId?: number;
}
