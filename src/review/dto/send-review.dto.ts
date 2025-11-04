import { ApiProperty } from '@nestjs/swagger';
import {
  IsIn,
  IsInt,
  IsNumber,
  IsOptional,
  IsPositive,
  IsString,
} from 'class-validator';

const allowedRatings = [1, 1.5, 2, 2.5, 3, 3.5, 4, 4.5, 5];

export class SendReviewDto {
  @ApiProperty({
    example: 'Коляска очень вместительная! 5 звёзд!',
  })
  @IsString({ message: 'Текст отзыва должен быть строкой' })
  @IsOptional()
  text?: string;
  @ApiProperty({
    example: 5,
  })
  @IsNumber(
    {},
    {
      message: 'Рейтинг должен быть числом',
    },
  )
  @IsIn(allowedRatings, {
    message:
      'Рейтинг должен быть одним из следующих чисел: 1, 1.5, 2, 2.5, 3, 3.5, 4, 4.5, 5',
  })
  rating: number;
  @ApiProperty({
    example: 1,
  })
  @IsNumber({}, { message: 'Id оцениваемого продавца должен быть числом' })
  @IsInt({ message: 'Id оцениваемого продавца должен быть целым числом' })
  @IsPositive({
    message: 'Id оцениваемого продавца должен быть положительным числом',
  })
  reviewedUserId: number;
}
