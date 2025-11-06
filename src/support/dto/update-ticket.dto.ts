import { IsEnum, IsOptional } from 'class-validator';
import { TicketStatus, TicketPriority } from '@prisma/client';

export class UpdateTicketDto {
  @IsOptional()
  @IsEnum(TicketStatus, { message: 'Некорректный статус тикета' })
  status?: TicketStatus;

  @IsOptional()
  @IsEnum(TicketPriority, { message: 'Некорректный приоритет' })
  priority?: TicketPriority;
}
