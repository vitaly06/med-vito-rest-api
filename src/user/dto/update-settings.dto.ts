import { Transform } from 'class-transformer';
import { IsBoolean, IsOptional, IsString } from 'class-validator';

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
}
