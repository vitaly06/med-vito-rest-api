import { ApiProperty } from '@nestjs/swagger';

import { IsNotEmpty, IsOptional, IsString, Matches } from 'class-validator';

export class CreateCategoryDto {
  @ApiProperty({
    example: 'Автомобили',
  })
  @IsString({ message: 'Название категории должно быть строкой' })
  @IsNotEmpty({ message: 'Название категории обязательно для заполнения' })
  name: string;

  @ApiProperty({
    example: 'avtomobili',
    description:
      'URL-friendly идентификатор (если не указан, генерируется автоматически)',
    required: false,
  })
  @IsOptional()
  @IsString({ message: 'Slug должен быть строкой' })
  @Matches(/^[a-z0-9]+(?:-[a-z0-9]+)*$/, {
    message:
      'Slug должен содержать только латинские буквы в нижнем регистре, цифры и дефисы',
  })
  slug?: string;
}
