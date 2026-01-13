import { CACHE_MANAGER } from '@nestjs/cache-manager';
import {
  BadRequestException,
  ForbiddenException,
  Inject,
  Injectable,
  NotFoundException,
} from '@nestjs/common';

import * as cacheManager from 'cache-manager';

import { TicketStatus } from '@/generated/prisma/enums';

import { PrismaService } from '../prisma/prisma.service';
import {
  CreateTicketDto,
  GetTicketsQueryDto,
  SendSupportMessageDto,
  UpdateTicketDto,
} from './dto';

@Injectable()
export class SupportService {
  constructor(
    private prisma: PrismaService,
    @Inject(CACHE_MANAGER) private cacheManager: cacheManager.Cache,
  ) {}

  async getSessionData(sessionId: string) {
    const sessionDataStr = await this.cacheManager.get<string>(
      `session:${sessionId}`,
    );
    if (!sessionDataStr) {
      return null;
    }
    return JSON.parse(sessionDataStr);
  }

  async getUserById(userId: number) {
    return await this.prisma.user.findUnique({
      where: { id: userId },
      include: { role: true },
    });
  }

  // Создание нового тикета
  async createTicket(userId: number, createTicketDto: CreateTicketDto) {
    const { message, ...ticketData } = createTicketDto;

    // Создаем тикет и первое сообщение в одной транзакции
    const result = await this.prisma.$transaction(async (prisma) => {
      // Создаем тикет
      const ticket = await prisma.supportTicket.create({
        data: {
          ...ticketData,
          userId,
        },
        include: {
          user: {
            select: {
              id: true,
              fullName: true,
              email: true,
            },
          },
          moderator: {
            select: {
              id: true,
              fullName: true,
              email: true,
            },
          },
        },
      });

      // Создаем первое сообщение
      await prisma.supportMessage.create({
        data: {
          ticketId: ticket.id,
          authorId: userId,
          text: message,
        },
      });

      return ticket;
    });

    return result;
  }

  // Получение тикетов пользователя
  async getUserTickets(userId: number, query: GetTicketsQueryDto) {
    const { page = 1, limit = 10, status, priority, theme, search } = query;
    const skip = (page - 1) * limit;

    const where: any = { userId };

    if (status) where.status = status;
    if (priority) where.priority = priority;
    if (theme) where.theme = theme;
    if (search) {
      where.OR = [
        { subject: { contains: search, mode: 'insensitive' } },
        {
          messages: {
            some: { text: { contains: search, mode: 'insensitive' } },
          },
        },
      ];
    }

    const [tickets, total] = await Promise.all([
      this.prisma.supportTicket.findMany({
        where,
        include: {
          user: {
            select: {
              id: true,
              fullName: true,
              email: true,
            },
          },
          moderator: {
            select: {
              id: true,
              fullName: true,
              email: true,
            },
          },
          messages: {
            orderBy: { sentAt: 'desc' },
            take: 1,
            include: {
              author: {
                select: {
                  id: true,
                  fullName: true,
                },
              },
            },
          },
        },
        orderBy: { updatedAt: 'desc' },
        skip,
        take: limit,
      }),
      this.prisma.supportTicket.count({ where }),
    ]);

    return {
      tickets: tickets.map((ticket) => ({
        ...ticket,
        lastMessage: ticket.messages[0] || null,
        messages: undefined,
      })),
      pagination: {
        page,
        limit,
        total,
        totalPages: Math.ceil(total / limit),
      },
    };
  }

  // Получение всех тикетов (для модераторов)
  async getAllTickets(query: GetTicketsQueryDto) {
    const { page = 1, limit = 10, status, priority, theme, search } = query;
    const skip = (page - 1) * limit;

    const where: any = {};

    if (status) where.status = status;
    if (priority) where.priority = priority;
    if (theme) where.theme = theme;
    if (search) {
      where.OR = [
        { subject: { contains: search, mode: 'insensitive' } },
        { user: { fullName: { contains: search, mode: 'insensitive' } } },
        {
          messages: {
            some: { text: { contains: search, mode: 'insensitive' } },
          },
        },
      ];
    }

    const [tickets, total] = await Promise.all([
      this.prisma.supportTicket.findMany({
        where,
        include: {
          user: {
            select: {
              id: true,
              fullName: true,
              email: true,
            },
          },
          moderator: {
            select: {
              id: true,
              fullName: true,
              email: true,
            },
          },
          messages: {
            orderBy: { sentAt: 'desc' },
            take: 1,
            include: {
              author: {
                select: {
                  id: true,
                  fullName: true,
                },
              },
            },
          },
        },
        orderBy: { updatedAt: 'desc' },
        skip,
        take: limit,
      }),
      this.prisma.supportTicket.count({ where }),
    ]);

    return {
      tickets: tickets.map((ticket) => ({
        ...ticket,
        lastMessage: ticket.messages[0] || null,
        messages: undefined,
      })),
      pagination: {
        page,
        limit,
        total,
        totalPages: Math.ceil(total / limit),
      },
    };
  }

  // Получение конкретного тикета с сообщениями
  async getTicket(ticketId: number, userId: number, isModerator = false) {
    const ticket = await this.prisma.supportTicket.findUnique({
      where: { id: ticketId },
      include: {
        user: {
          select: {
            id: true,
            fullName: true,
            email: true,
          },
        },
        moderator: {
          select: {
            id: true,
            fullName: true,
            email: true,
          },
        },
        messages: {
          include: {
            author: {
              select: {
                id: true,
                fullName: true,
                email: true,
                role: {
                  select: {
                    name: true,
                  },
                },
              },
            },
          },
          orderBy: { sentAt: 'asc' },
        },
      },
    });

    if (!ticket) {
      throw new NotFoundException('Тикет не найден');
    }

    // Проверяем доступ: пользователь может видеть только свои тикеты, модератор - все
    if (!isModerator && ticket.userId !== userId) {
      throw new ForbiddenException('Нет доступа к этому тикету');
    }

    return ticket;
  }

  // Отправка сообщения в тикет
  async sendMessage(
    ticketId: number,
    userId: number,
    sendMessageDto: SendSupportMessageDto,
    isModerator = false,
  ) {
    const ticket = await this.prisma.supportTicket.findUnique({
      where: { id: ticketId },
    });

    if (!ticket) {
      throw new NotFoundException('Тикет не найден');
    }

    // Проверяем доступ
    if (!isModerator && ticket.userId !== userId) {
      throw new ForbiddenException('Нет доступа к этому тикету');
    }

    // Проверяем, что тикет не закрыт
    if (ticket.status === TicketStatus.CLOSED) {
      throw new BadRequestException(
        'Нельзя отправлять сообщения в закрытый тикет',
      );
    }

    // Создаем сообщение и обновляем тикет в транзакции
    const result = await this.prisma.$transaction(async (prisma) => {
      // Создаем сообщение
      const message = await prisma.supportMessage.create({
        data: {
          ticketId,
          authorId: userId,
          text: sendMessageDto.text,
        },
        include: {
          author: {
            select: {
              id: true,
              fullName: true,
              email: true,
              role: {
                select: {
                  name: true,
                },
              },
            },
          },
        },
      });

      // Обновляем статус тикета и назначаем модератора если нужно
      const updateData: any = { updatedAt: new Date() };

      if (isModerator) {
        // Если отвечает модератор, меняем статус на "в работе" и назначаем себя
        updateData.status = TicketStatus.IN_PROGRESS;
        updateData.moderatorId = userId;
      } else if (ticket.status === TicketStatus.RESOLVED) {
        // Если пользователь отвечает на решенный тикет, открываем его снова
        updateData.status = TicketStatus.OPEN;
      }

      await prisma.supportTicket.update({
        where: { id: ticketId },
        data: updateData,
      });

      return message;
    });

    return result;
  }

  // Обновление тикета (для модераторов)
  async updateTicket(
    ticketId: number,
    moderatorId: number,
    updateTicketDto: UpdateTicketDto,
  ) {
    const ticket = await this.prisma.supportTicket.findUnique({
      where: { id: ticketId },
    });

    if (!ticket) {
      throw new NotFoundException('Тикет не найден');
    }

    const updatedTicket = await this.prisma.supportTicket.update({
      where: { id: ticketId },
      data: {
        ...updateTicketDto,
        moderatorId, // Назначаем модератора при обновлении
        updatedAt: new Date(),
      },
      include: {
        user: {
          select: {
            id: true,
            fullName: true,
            email: true,
          },
        },
        moderator: {
          select: {
            id: true,
            fullName: true,
            email: true,
          },
        },
      },
    });

    return updatedTicket;
  }

  // Назначение тикета модератору
  async assignTicket(ticketId: number, moderatorId: number) {
    const ticket = await this.prisma.supportTicket.findUnique({
      where: { id: ticketId },
    });

    if (!ticket) {
      throw new NotFoundException('Тикет не найден');
    }

    const updatedTicket = await this.prisma.supportTicket.update({
      where: { id: ticketId },
      data: {
        moderatorId,
        status: TicketStatus.IN_PROGRESS,
        updatedAt: new Date(),
      },
      include: {
        user: {
          select: {
            id: true,
            fullName: true,
            email: true,
          },
        },
        moderator: {
          select: {
            id: true,
            fullName: true,
            email: true,
          },
        },
      },
    });

    return updatedTicket;
  }

  // Получение статистики тикетов
  async getTicketStats() {
    const [
      totalTickets,
      openTickets,
      inProgressTickets,
      resolvedTickets,
      closedTickets,
      ticketsByTheme,
      ticketsByPriority,
    ] = await Promise.all([
      this.prisma.supportTicket.count(),
      this.prisma.supportTicket.count({ where: { status: TicketStatus.OPEN } }),
      this.prisma.supportTicket.count({
        where: { status: TicketStatus.IN_PROGRESS },
      }),
      this.prisma.supportTicket.count({
        where: { status: TicketStatus.RESOLVED },
      }),
      this.prisma.supportTicket.count({
        where: { status: TicketStatus.CLOSED },
      }),
      this.prisma.supportTicket.groupBy({
        by: ['theme'],
        _count: { theme: true },
      }),
      this.prisma.supportTicket.groupBy({
        by: ['priority'],
        _count: { priority: true },
      }),
    ]);

    return {
      total: totalTickets,
      byStatus: {
        open: openTickets,
        inProgress: inProgressTickets,
        resolved: resolvedTickets,
        closed: closedTickets,
      },
      byTheme: ticketsByTheme.map((item) => ({
        theme: item.theme,
        count: item._count.theme,
      })),
      byPriority: ticketsByPriority.map((item) => ({
        priority: item.priority,
        count: item._count.priority,
      })),
    };
  }
}
