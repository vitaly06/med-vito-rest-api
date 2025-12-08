import { Transform } from 'class-transformer';
import { IsInt, IsNumber, IsPositive } from 'class-validator';

export class AddPromotionDto {
  @IsInt({ message: 'Id товара должен быть целым числом' })
  @IsNumber({}, { message: 'Id товара должен быть числом' })
  @IsPositive({ message: 'Id товара должен быть положительным числом' })
  @Transform(({ value }) => parseInt(value))
  productId: number;
  @IsInt({ message: 'Id продвижения должен быть целым числом' })
  @IsNumber({}, { message: 'Id продвижения должен быть числом' })
  @IsPositive({ message: 'Id продвижения должен быть положительным числом' })
  @Transform(({ value }) => parseInt(value))
  promotionId: number;
  @IsInt({ message: 'Кол-во дней должно быть целым числом' })
  @IsNumber({}, { message: 'Кол-во дней должно быть числом' })
  @IsPositive({ message: 'Кол-во дней должно быть положительным числом' })
  @Transform(({ value }) => parseInt(value))
  days: number;
}
