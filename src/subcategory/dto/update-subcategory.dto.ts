import { ApiProperty } from '@nestjs/swagger';
import { Type } from 'class-transformer';
import { IsInt, IsNumber, IsPositive, IsString } from 'class-validator';

export class UpdateSubcategoryDto {
  @ApiProperty({
    example: 'Резина и диски',
  })
  @IsString({ message: 'Название подкатегории должно быть строкой' })
  name?: string;

  @ApiProperty({
    example: 1,
  })
  @IsNumber({}, { message: 'Id категории должен быть числом' })
  @IsPositive({ message: 'Id категории должен быть положительным числом' })
  @IsInt({ message: 'Id категории должен быть целым числом' })
  @Type(() => Number)
  categoryId?: number;
}
