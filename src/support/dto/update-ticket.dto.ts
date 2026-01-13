import { IsEnum, IsOptional } from 'class-validator';

import { TicketPriority, TicketStatus } from '@/generated/prisma/enums';

export class UpdateTicketDto {
  @IsOptional()
  @IsEnum(TicketStatus, { message: 'Некорректный статус тикета' })
  status?: TicketStatus;

  @IsOptional()
  @IsEnum(TicketPriority, { message: 'Некорректный приоритет' })
  priority?: TicketPriority;
}
