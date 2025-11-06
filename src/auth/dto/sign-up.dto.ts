import { ApiProperty } from '@nestjs/swagger';
import { ProfileType } from '@prisma/client';
import {
  IsEmail,
  IsEnum,
  IsNotEmpty,
  IsString,
  MinLength,
} from 'class-validator';

export class SignUpDto {
  @ApiProperty({
    example: 'Садиков Виталий Дмитриевич',
  })
  @IsString({ message: 'ФИО должно быть строкой' })
  @IsNotEmpty({ message: 'ФИО обязательно для заполнения' })
  fullName: string;
  @ApiProperty({
    example: 'vitaly.sadikov1@yandex.ru',
  })
  @IsString({ message: 'Email должен быть строкой' })
  @IsNotEmpty({ message: 'Email обязателен для заполнения' })
  @IsEmail({}, { message: 'Неверный формат почты' })
  email: string;
  @ApiProperty({
    example: '+79510341677',
  })
  @IsString({ message: 'Номер телефона должен быть строкой' })
  @IsNotEmpty({ message: 'Номер телефона обязателен для заполнения' })
  phoneNumber: string;
  @ApiProperty({
    example: '123456',
  })
  @IsString({ message: 'Пароль должен быть строкой' })
  @IsNotEmpty({ message: 'Пароль обязателен для заполнения' })
  @MinLength(6, { message: 'Длина пароля не менее 6 символов' })
  password: string;
  // @ApiProperty({
  //   example: 'INDIVIDUAL',
  // })
  // @IsEnum(ProfileType, { message: 'Неверный тип профиля' })
  // profileType: ProfileType;
}
