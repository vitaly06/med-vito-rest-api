import { ApiPropertyOptional } from '@nestjs/swagger';

import { Transform } from 'class-transformer';
import { IsEnum, IsNumberString, IsOptional, IsString } from 'class-validator';

export enum PeriodEnum {
  DAY = 'day',
  WEEK = 'week',
  MONTH = 'month',
  QUARTER = 'quarter',
  HALF_YEAR = 'half-year',
  YEAR = 'year',
}

export class StatisticsQueryDto {
  @ApiPropertyOptional({
    description: 'Период для анализа статистики',
    enum: PeriodEnum,
    example: 'month',
  })
  @IsOptional()
  @IsEnum(PeriodEnum)
  period?: PeriodEnum;

  @ApiPropertyOptional({
    description: 'ID категории товаров для фильтрации',
    example: 1,
  })
  @IsOptional()
  @IsNumberString()
  @Transform(({ value }) => (value ? parseInt(value, 10) : undefined))
  categoryId?: number;

  @ApiPropertyOptional({
    description: 'Фильтр по региону',
    example: 'Москва',
  })
  @IsOptional()
  @IsString()
  region?: string;

  @ApiPropertyOptional({
    description: 'ID конкретного товара для анализа',
    example: '1',
  })
  @IsOptional()
  @IsNumberString()
  @Transform(({ value }) => (value ? parseInt(value, 10) : undefined))
  productId?: number;
}
