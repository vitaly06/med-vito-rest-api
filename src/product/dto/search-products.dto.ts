import { ApiPropertyOptional } from '@nestjs/swagger';
import { IsOptional, IsString, IsNumber, IsIn } from 'class-validator';
import { Transform, Type } from 'class-transformer';

export class SearchProductsDto {
  @ApiPropertyOptional({
    description: 'Поисковый запрос по названию или описанию товара',
    example: 'iPhone 15 Pro',
  })
  @IsOptional()
  @IsString()
  search?: string;

  @ApiPropertyOptional({
    description: 'Slug категории для фильтрации',
    example: 'lichnye-veschi',
  })
  @IsOptional()
  @IsString()
  categorySlug?: string;

  @ApiPropertyOptional({
    description: 'Slug подкатегории для фильтрации',
    example: 'odezhda',
  })
  @IsOptional()
  @IsString()
  subCategorySlug?: string;

  @ApiPropertyOptional({
    description: 'Slug типа подкатегории для фильтрации',
    example: 'muzhskaya',
  })
  @IsOptional()
  @IsString()
  typeSlug?: string;

  @ApiPropertyOptional({
    description: 'Минимальная цена',
    example: 1000,
  })
  @IsOptional()
  @Type(() => Number)
  @IsNumber()
  minPrice?: number;

  @ApiPropertyOptional({
    description: 'Максимальная цена',
    example: 100000,
  })
  @IsOptional()
  @Type(() => Number)
  @IsNumber()
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
    description: 'Тип продавца',
    enum: ['INDIVIDUAL', 'OOO', 'IP'],
    example: 'INDIVIDUAL',
  })
  @IsOptional()
  @IsIn(['INDIVIDUAL', 'OOO', 'IP'])
  profileType?: 'INDIVIDUAL' | 'OOO' | 'IP';

  @ApiPropertyOptional({
    description:
      'Фильтр по характеристикам товара в формате JSON. Например: {"1":"XXL","2":"Чёрный"} где ключ - ID поля, значение - искомое значение',
    example: '{"1":"XXL","2":"Чёрный"}',
  })
  @IsOptional()
  @Transform(({ value }) => {
    if (typeof value === 'string') {
      try {
        return JSON.parse(value);
      } catch {
        return undefined;
      }
    }
    return value;
  })
  fieldValues?: Record<string, string>;

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
  @Type(() => Number)
  @IsNumber()
  page?: number;

  @ApiPropertyOptional({
    description: 'Количество товаров на странице',
    example: 20,
  })
  @IsOptional()
  @Type(() => Number)
  @IsNumber()
  limit?: number;
}
