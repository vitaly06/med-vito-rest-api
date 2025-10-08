import { IsNotEmpty, IsString, MinLength } from 'class-validator';

export class SignInDto {
  @IsString({ message: 'Данные должны быть строкой' })
  @IsNotEmpty({ message: 'Данные обязательны для заполнения' })
  login: string;
  @IsString({ message: 'Пароль должен быть строкой' })
  @IsNotEmpty({ message: 'Пароль обязателен для заполнения' })
  @MinLength(6, { message: 'Длина пароля не менее 6 символов' })
  password: string;
}
