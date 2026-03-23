import {
  Controller,
  DefaultValuePipe,
  Get,
  Param,
  ParseIntPipe,
  Query,
  UseGuards,
} from '@nestjs/common';
import {
  ApiBearerAuth,
  ApiOperation,
  ApiQuery,
  ApiTags,
} from '@nestjs/swagger';

import { AdminSessionAuthGuard } from 'src/auth/guards/admin-session-auth.guard';
import { PrismaService } from 'src/prisma/prisma.service';

enum ModerationFilter {
  ALL = 'ALL',
  DENIED = 'DENIED',
  MANUAL = 'MANUAL',
  APPROVED_AI = 'APPROVED_AI',
}

@ApiTags('Moderation (Admin)')
@ApiBearerAuth()
@UseGuards(AdminSessionAuthGuard)
@Controller('admin/moderation')
export class ModerationController {
  constructor(private readonly prisma: PrismaService) {}

  /**
   * GET /admin/moderation/products?filter=ALL|DENIED|MANUAL&page=1
   *
   * Returns products that were rejected by AI (DENIDED) or sent to manual review (AI_REVIEWED), together with the AI reason when available.
   */
  @Get('products')
  @ApiOperation({
    summary: 'Список товаров отклонённых или на ручной проверке',
    description:
      '`filter=DENIED` — только отклонённые ИИ; ' +
      '`filter=MANUAL` — только на ручной проверке после ИИ; ' +
      '`filter=APPROVED_AI` — только одобренные ИИ и видимые в ленте; ' +
      '`filter=ALL` — отклонённые/ручная проверка/одобренные ИИ. По умолчанию ALL.',
  })
  @ApiQuery({ name: 'filter', enum: ModerationFilter, required: false })
  @ApiQuery({ name: 'page', required: false, type: Number })
  async getModerationList(
    @Query('filter', new DefaultValuePipe(ModerationFilter.ALL))
    filter: ModerationFilter,
    @Query('page', new DefaultValuePipe(1), ParseIntPipe) page: number,
  ) {
    const take = 20;
    const skip = (page - 1) * take;

    const aiApprovedWhere = {
      moderateState: 'APPROVED' as const,
      isHide: false as const,
      moderationRejectionReason: 'Одобрено ИИ автоматически',
    };

    let where:
      | { moderateState: 'DENIDED' }
      | {
          moderateState: 'APPROVED';
          isHide: false;
          moderationRejectionReason: string;
        }
      | {
          OR: Array<
            | { moderateState: 'AI_REVIEWED' }
            | { moderateState: 'MODERATE' }
            | { moderateState: 'DENIDED' }
            | {
                moderateState: 'APPROVED';
                isHide: false;
                moderationRejectionReason: string;
              }
          >;
        };

    if (filter === ModerationFilter.DENIED) {
      where = { moderateState: 'DENIDED' };
    } else if (filter === ModerationFilter.APPROVED_AI) {
      where = aiApprovedWhere;
    } else if (filter === ModerationFilter.MANUAL) {
      where = {
        OR: [{ moderateState: 'AI_REVIEWED' }, { moderateState: 'MODERATE' }],
      };
    } else {
      where = {
        OR: [
          { moderateState: 'DENIDED' },
          { moderateState: 'AI_REVIEWED' },
          { moderateState: 'MODERATE' },
          aiApprovedWhere,
        ],
      };
    }

    const [items, total] = await Promise.all([
      this.prisma.product.findMany({
        where,
        skip,
        take,
        orderBy: { updatedAt: 'desc' },
        select: {
          id: true,
          name: true,
          price: true,
          images: true,
          moderateState: true,
          moderationRejectionReason: true,
          createdAt: true,
          updatedAt: true,
          category: { select: { id: true, name: true } },
          subCategory: { select: { id: true, name: true } },
          user: {
            select: {
              id: true,
              fullName: true,
              email: true,
              phoneNumber: true,
            },
          },
        },
      }),
      this.prisma.product.count({ where }),
    ]);

    return {
      items,
      total,
      page,
      pages: Math.ceil(total / take),
    };
  }

  /**
   * GET /admin/moderation/products/:id
   *
   * Full details of a single product under moderation.
   */
  @Get('products/:id')
  @ApiOperation({ summary: 'Детали конкретного товара на модерации' })
  async getModerationItem(@Param('id', ParseIntPipe) id: number) {
    return await this.prisma.product.findUniqueOrThrow({
      where: { id },
      select: {
        id: true,
        name: true,
        price: true,
        description: true,
        images: true,
        videoUrl: true,
        moderateState: true,
        moderationRejectionReason: true,
        createdAt: true,
        updatedAt: true,
        category: { select: { id: true, name: true } },
        subCategory: { select: { id: true, name: true } },
        type: { select: { id: true, name: true } },
        user: {
          select: {
            id: true,
            fullName: true,
            email: true,
            phoneNumber: true,
            profileType: true,
          },
        },
        fieldValues: {
          select: {
            value: true,
            field: { select: { id: true, name: true } },
          },
        },
      },
    });
  }
}
