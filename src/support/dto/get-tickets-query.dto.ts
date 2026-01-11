import { TicketPriority, TicketStatus, TicketTheme } from '@prisma/client';
import { Transform } from 'class-transformer';
import { IsEnum, IsOptional, IsString } from 'class-validator';

export class GetTicketsQueryDto {
  @IsOptional()
  @IsEnum(TicketStatus, { message: 'Некорректный статус тикета' })
  status?: TicketStatus;

  @IsOptional()
  @IsEnum(TicketPriority, { message: 'Некорректный приоритет' })
  priority?: TicketPriority;

  @IsOptional()
  @IsEnum(TicketTheme, { message: 'Некорректная тема' })
  theme?: TicketTheme;

  @IsOptional()
  @Transform(({ value }) => parseInt(value))
  page?: number = 1;

  @IsOptional()
  @Transform(({ value }) => parseInt(value))
  limit?: number = 10;

  @IsOptional()
  @IsString({ message: 'Поиск должен быть строкой' })
  search?: string;
}
