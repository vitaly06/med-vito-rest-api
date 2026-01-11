import { Transform } from 'class-transformer';
import { IsBoolean, IsEnum, IsOptional, IsString } from 'class-validator';

import { ProfileType } from '@/auth/enum/profile-type.enum';

export class UpdateSettingsDto {
  @IsOptional()
  @IsString({ message: 'ФИО должно быть строкой' })
  fullName?: string;
  @IsOptional()
  @IsString({ message: 'Номер телефона должен быть строкой' })
  phoneNumber?: string;
  @IsOptional()
  @Transform(({ value }) => value === 'true' || value === true)
  @IsBoolean({
    message: 'Поле "Отвечаете ли вы на звонки?" должно быть true/false',
  })
  isAnswersCall?: boolean;

  @IsOptional()
  @IsEnum(ProfileType, {
    message: 'Неверный тип профиля. Доступные: INDIVIDUAL, OOO, IP',
  })
  profileType?: ProfileType;
}
