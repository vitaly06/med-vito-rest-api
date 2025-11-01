import { ApiProperty } from '@nestjs/swagger';
import { Type } from 'class-transformer';
import {
  IsInt,
  IsNotEmpty,
  IsNumber,
  IsPositive,
  IsString,
  MinLength,
} from 'class-validator';

export class ChangePasswordDto {
  @ApiProperty({
    example: 1,
  })
  @Type(() => Number)
  @IsNumber({}, { message: 'Id пользователя должен быть числом' })
  @IsInt({ message: 'Id пользователя должен быть целым числом' })
  @IsNotEmpty({ message: 'Id пользователя обязателен для заполнения' })
  @IsPositive({ message: 'Id пользователя должен быть положительным числом' })
  userId: number;
  @ApiProperty({
    example: '123456',
  })
  @IsString({ message: 'Пароль должен быть строкой' })
  @IsNotEmpty({ message: 'Пароль обязателен для заполнения' })
  @MinLength(6, { message: 'Длина пароля не менее 6 символов' })
  password: string;
}
