import { ApiProperty } from '@nestjs/swagger';

import {
  IsEmail,
  IsEnum,
  IsNumber,
  IsOptional,
  IsString,
  Min,
} from 'class-validator';
import { ProfileType } from 'prisma/generated/enums';

export class UpdateUserDto {
  @ApiProperty({
    description: 'ФИО пользователя',
    example: 'Иванов Иван Иванович',
    required: false,
  })
  @IsOptional()
  @IsString({ message: 'ФИО должно быть строкой' })
  fullName?: string;

  @ApiProperty({
    description: 'Email пользователя',
    example: 'user@example.com',
    required: false,
  })
  @IsOptional()
  @IsEmail({}, { message: 'Некорректный email' })
  email?: string;

  @ApiProperty({
    description: 'Номер телефона',
    example: '+79001234567',
    required: false,
  })
  @IsOptional()
  @IsString({ message: 'Номер телефона должен быть строкой' })
  phoneNumber?: string;

  @ApiProperty({
    description: 'Тип профиля',
    enum: ProfileType,
    example: ProfileType.INDIVIDUAL,
    required: false,
  })
  @IsOptional()
  @IsEnum(ProfileType, {
    message: 'Неверный тип профиля. Доступные: INDIVIDUAL, OOO, IP',
  })
  profileType?: ProfileType;

  @ApiProperty({
    description: 'Бонусный баланс пользователя',
    example: 500.0,
    required: false,
  })
  @IsOptional()
  @IsNumber({}, { message: 'Бонусный баланс должен быть числом' })
  @Min(0, { message: 'Бонусный баланс не может быть отрицательным' })
  bonusBalance?: number;
}
