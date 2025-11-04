import { ApiPropertyOptional } from '@nestjs/swagger';
import { IsOptional, IsString, IsNumberString, IsIn } from 'class-validator';
import { Transform } from 'class-transformer';

export class SearchProductsDto {
  @ApiPropertyOptional({
    description: 'Поисковый запрос по названию, описанию или бренду товара',
    example: 'iPhone 15 Pro',
  })
  @IsOptional()
  @IsString()
  search?: string;

  @ApiPropertyOptional({
    description: 'ID категории для фильтрации',
    example: 1,
  })
  @IsOptional()
  @IsNumberString()
  @Transform(({ value }) => (value ? parseInt(value, 10) : undefined))
  categoryId?: number;

  @ApiPropertyOptional({
    description: 'ID подкатегории для фильтрации',
    example: 1,
  })
  @IsOptional()
  @IsNumberString()
  @Transform(({ value }) => (value ? parseInt(value, 10) : undefined))
  subCategoryId?: number;

  @ApiPropertyOptional({
    description: 'Минимальная цена',
    example: 1000,
  })
  @IsOptional()
  @IsNumberString()
  @Transform(({ value }) => (value ? parseInt(value, 10) : undefined))
  minPrice?: number;

  @ApiPropertyOptional({
    description: 'Максимальная цена',
    example: 100000,
  })
  @IsOptional()
  @IsNumberString()
  @Transform(({ value }) => (value ? parseInt(value, 10) : undefined))
  maxPrice?: number;

  @ApiPropertyOptional({
    description: 'Состояние товара',
    enum: ['NEW', 'USED'],
    example: 'NEW',
  })
  @IsOptional()
  @IsIn(['NEW', 'USED'])
  state?: 'NEW' | 'USED';

  @ApiPropertyOptional({
    description: 'Регион/город',
    example: 'Москва',
  })
  @IsOptional()
  @IsString()
  region?: string;

  @ApiPropertyOptional({
    description: 'Бренд товара',
    example: 'Apple',
  })
  @IsOptional()
  @IsString()
  brand?: string;

  @ApiPropertyOptional({
    description: 'Сортировка результатов',
    enum: ['price_asc', 'price_desc', 'date_desc', 'date_asc', 'relevance'],
    example: 'price_asc',
  })
  @IsOptional()
  @IsIn(['price_asc', 'price_desc', 'date_desc', 'date_asc', 'relevance'])
  sortBy?: 'price_asc' | 'price_desc' | 'date_desc' | 'date_asc' | 'relevance';

  @ApiPropertyOptional({
    description: 'Номер страницы',
    example: 1,
  })
  @IsOptional()
  @IsNumberString()
  @Transform(({ value }) => (value ? parseInt(value, 10) : 1))
  page?: number;

  @ApiPropertyOptional({
    description: 'Количество товаров на странице',
    example: 20,
  })
  @IsOptional()
  @IsNumberString()
  @Transform(({ value }) => (value ? parseInt(value, 10) : 20))
  limit?: number;
}
