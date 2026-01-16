import { ApiProperty } from '@nestjs/swagger';

import { IsEnum, IsOptional, IsString } from 'class-validator';

import { BannerPlace } from '../entities/banner-place.enum';

export class UpdateBannerDto {
  @ApiProperty({
    description: 'Изображение баннера',
    type: 'string',
    format: 'binary',
    required: false,
  })
  image?: any;

  @ApiProperty({
    description: 'Текст баннера',
    required: false,
  })
  @IsOptional()
  @IsString({ message: 'Текст баннера должен быть строкой' })
  name?: string;

  @ApiProperty({
    description: 'Место размещения баннера на сайте',
    enum: BannerPlace,
    example: BannerPlace.PRODUCT_FEED,
    required: false,
  })
  @IsOptional()
  @IsEnum(BannerPlace, { message: 'Некорректное значение места размещения' })
  place?: BannerPlace;

  @ApiProperty({
    description: 'Url, куда переносит баннер',
  })
  @IsOptional()
  navigateToUrl?: string;
}
