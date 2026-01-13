import { ApiProperty } from '@nestjs/swagger';

import { IsEnum, IsNotEmpty } from 'class-validator';

import { BannerPlace } from '../entities/banner-place.enum';

export class CreateBannerDto {
  @ApiProperty({
    description: 'Изображение баннера',
    type: 'string',
    format: 'binary',
  })
  @IsNotEmpty({ message: 'Изображение баннера обязательно' })
  image: any;

  @ApiProperty({
    description: 'Место размещения баннера на сайте',
    enum: BannerPlace,
    example: BannerPlace.PRODUCT_FEED,
  })
  @IsEnum(BannerPlace, { message: 'Некорректное значение места размещения' })
  @IsNotEmpty({ message: 'Место размещения баннера обязательно' })
  place: BannerPlace;

  @ApiProperty({
    description: 'Url, куда переносит баннер',
  })
  @IsNotEmpty({ message: 'Navigate url обязателен для заполнения' })
  navigateToUrl: string;
}
