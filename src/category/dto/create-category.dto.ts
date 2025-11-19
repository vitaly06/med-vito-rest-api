import { ApiProperty } from '@nestjs/swagger';
import { IsNotEmpty, IsString } from 'class-validator';

export class CreateCategoryDto {
  @ApiProperty({
    example: 'Автомобили',
  })
  @IsString({ message: 'Название категории должно быть строкой' })
  @IsNotEmpty({ message: 'Название категории обязательно для заполнения' })
  name: string;
}
