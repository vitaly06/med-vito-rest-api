import { IsNotEmpty, IsString, MaxLength } from 'class-validator';

export class SendSupportMessageDto {
  @IsString({ message: 'Сообщение должно быть строкой' })
  @IsNotEmpty({ message: 'Сообщение обязательно' })
  @MaxLength(2000, { message: 'Сообщение не должно превышать 2000 символов' })
  text: string;
}
