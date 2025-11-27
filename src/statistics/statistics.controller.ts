import { Controller, Get, Query, UseGuards, Req } from '@nestjs/common';
import {
  ApiTags,
  ApiOperation,
  ApiResponse,
  ApiBearerAuth,
  ApiQuery,
} from '@nestjs/swagger';
import { Request } from 'express';
import { StatisticsService } from './statistics.service';
import { SessionAuthGuard } from 'src/auth/guards/session-auth.guard';
import { StatisticsQueryDto, PeriodEnum } from './dto/statistics-query.dto';

@ApiTags('Statistics')
@ApiBearerAuth()
@Controller('statistics')
export class StatisticsController {
  constructor(private readonly statisticsService: StatisticsService) {}

  @ApiOperation({
    summary: 'Получение статистики пользователя',
    description: `
      Получение детальной статистики по товарам пользователя.
      Включает количество просмотров, показов телефона и добавлений в избранное.
      Все параметры фильтрации являются опциональными.
    `,
  })
  @ApiQuery({
    name: 'period',
    required: false,
    enum: PeriodEnum,
    description: 'Период для анализа статистики',
    example: 'month',
  })
  @ApiQuery({
    name: 'categoryId',
    required: false,
    type: Number,
    description: 'ID категории товаров для фильтрации',
    example: 1,
  })
  @ApiQuery({
    name: 'region',
    required: false,
    type: String,
    description: 'Фильтр по региону',
    example: 'Москва',
  })
  @ApiQuery({
    name: 'productId',
    required: false,
    type: Number,
    description: 'ID конкретного товара для анализа',
    example: 1,
  })
  @ApiResponse({
    status: 200,
    description: 'Статистика успешно получена',
    schema: {
      type: 'object',
      properties: {
        period: {
          type: 'string',
          example: 'month',
          description: 'Период анализа',
        },
        totalViews: { type: 'number', example: 150 },
        totalPhoneViews: { type: 'number', example: 25 },
        totalFavorites: { type: 'number', example: 40 },
      },
    },
  })
  @ApiResponse({
    status: 401,
    description: 'Неавторизован',
  })
  @UseGuards(SessionAuthGuard)
  @Get('analytic')
  async getUserStatistics(
    @Query() query: StatisticsQueryDto,
    @Req() req: Request & { user: any },
  ) {
    return await this.statisticsService.getUserStatistics(req.user.id, query);
  }

  @ApiOperation({
    summary: 'Аналитика для каждого товара',
  })
  @UseGuards(SessionAuthGuard)
  @Get('products-analytic')
  async getProductsAnalytic(@Req() req: Request & { user: any }) {
    return await this.statisticsService.getProductsAnalytic(req.user.id);
  }
}
