import { ApiProperty } from '@nestjs/swagger';

export class UpdateBannerDto {
  @ApiProperty({
    description: 'Изображение баннера',
    type: 'string',
    format: 'binary',
    required: false,
  })
  image?: any;
}
