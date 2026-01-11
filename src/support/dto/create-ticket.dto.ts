import { TicketPriority, TicketTheme } from '@prisma/client';
import { IsEnum, IsNotEmpty, IsString, MaxLength } from 'class-validator';

export class CreateTicketDto {
  @IsEnum(TicketTheme, { message: 'Некорректная тема обращения' })
  @IsNotEmpty({ message: 'Тема обращения обязательна' })
  theme: TicketTheme;

  @IsString({ message: 'Заголовок должен быть строкой' })
  @IsNotEmpty({ message: 'Заголовок обязателен' })
  @MaxLength(200, { message: 'Заголовок не должен превышать 200 символов' })
  subject: string;

  @IsString({ message: 'Сообщение должно быть строкой' })
  @IsNotEmpty({ message: 'Сообщение обязательно' })
  @MaxLength(2000, { message: 'Сообщение не должно превышать 2000 символов' })
  message: string;

  @IsEnum(TicketPriority, { message: 'Некорректный приоритет' })
  priority: TicketPriority = TicketPriority.MEDIUM;
}
