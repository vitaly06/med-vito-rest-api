import { Body, Controller, Get, Post, Req, UseGuards } from '@nestjs/common';
import {
  ApiBearerAuth,
  ApiBody,
  ApiOperation,
  ApiResponse,
  ApiTags,
} from '@nestjs/swagger';

import { SessionAuthGuard } from 'src/auth/guards/session-auth.guard';

import { AddPromotionDto } from './add-promotion.dto';
import { PromotionService } from './promotion.service';

@ApiTags('Promotion')
@Controller('promotion')
export class PromotionController {
  constructor(private readonly promotionService: PromotionService) {}

  @ApiOperation({
    summary: 'Все типы продвижения',
    description: 'Возвращает список всех доступных тарифов продвижения',
  })
  @ApiResponse({
    status: 200,
    description: 'Список тарифов продвижения',
    schema: {
      example: [
        {
          id: 1,
          name: 'Стандарт',
          pricePerDay: 50,
        },
        {
          id: 2,
          name: 'Люкс',
          pricePerDay: 100,
        },
      ],
    },
  })
  @Get('all-promotions')
  async allPromotions() {
    return await this.promotionService.allPromotions();
  }

  @ApiOperation({
    summary: 'Подключение продвижения к товару',
    description:
      'Активирует продвижение для товара. Деньги списываются сначала с основного баланса, затем с бонусного.',
  })
  @ApiBearerAuth()
  @ApiBody({
    type: AddPromotionDto,
    description: 'Данные для активации продвижения',
  })
  @ApiResponse({
    status: 201,
    description: 'Продвижение успешно активировано',
    schema: {
      example: {
        message: 'Продвижение успешно активировано',
        promotion: {
          id: 1,
          days: 7,
          totalPrice: 350,
          startDate: '2026-01-22T10:00:00.000Z',
          endDate: '2026-01-29T10:00:00.000Z',
        },
      },
    },
  })
  @ApiResponse({
    status: 400,
    description: 'Недостаточно средств на балансе',
  })
  @ApiResponse({
    status: 401,
    description: 'Не авторизован',
  })
  @ApiResponse({
    status: 403,
    description: 'Нельзя продвигать чужой товар',
  })
  @ApiResponse({
    status: 404,
    description: 'Товар или тип продвижения не найден',
  })
  @UseGuards(SessionAuthGuard)
  @Post('add-promotion')
  async addPromotion(
    @Body() dto: AddPromotionDto,
    @Req() req: Request & { user: any },
  ) {
    return await this.promotionService.addPromotion(dto, req.user.id);
  }
}
