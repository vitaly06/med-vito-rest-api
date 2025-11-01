import { ApiProperty } from '@nestjs/swagger';
import { IsEmail, IsNotEmpty, IsString } from 'class-validator';

export class ForgotPasswordDto {
  @ApiProperty({
    description: 'Почта для восстановления пароля',
    example: 'vitaly.sadikov1@yandex.ru',
  })
  @IsString({ message: 'Почта должна быть строкой' })
  @IsNotEmpty({ message: 'Почта обязательная для заполнения' })
  @IsEmail({}, { message: 'Неверный формат почты' })
  email: string;
}
