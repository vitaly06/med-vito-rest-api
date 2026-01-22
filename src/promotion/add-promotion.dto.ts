import { ApiProperty } from '@nestjs/swagger';
import { Transform } from 'class-transformer';
import { IsInt, IsNumber, IsPositive } from 'class-validator';

export class AddPromotionDto {
  @ApiProperty({
    description: 'ID товара для продвижения',
    example: 123,
    type: Number,
  })
  @IsInt({ message: 'Id товара должен быть целым числом' })
  @IsNumber({}, { message: 'Id товара должен быть числом' })
  @IsPositive({ message: 'Id товара должен быть положительным числом' })
  @Transform(({ value }) => parseInt(value))
  productId: number;

  @ApiProperty({
    description: 'ID типа продвижения (1 - стандарт, 2 - люкс)',
    example: 1,
    type: Number,
  })
  @IsInt({ message: 'Id продвижения должен быть целым числом' })
  @IsNumber({}, { message: 'Id продвижения должен быть числом' })
  @IsPositive({ message: 'Id продвижения должен быть положительным числом' })
  @Transform(({ value }) => parseInt(value))
  promotionId: number;

  @ApiProperty({
    description: 'Количество дней продвижения',
    example: 7,
    type: Number,
  })
  @IsInt({ message: 'Кол-во дней должно быть целым числом' })
  @IsNumber({}, { message: 'Кол-во дней должно быть числом' })
  @IsPositive({ message: 'Кол-во дней должно быть положительным числом' })
  @Transform(({ value }) => parseInt(value))
  days: number;
}
