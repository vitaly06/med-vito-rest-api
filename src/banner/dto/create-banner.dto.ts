import { ApiProperty } from '@nestjs/swagger';
import { IsNotEmpty } from 'class-validator';

export class CreateBannerDto {
  @ApiProperty({
    description: 'Изображение баннера',
    type: 'string',
    format: 'binary',
  })
  @IsNotEmpty({ message: 'Изображение баннера обязательно' })
  image: any;
}
