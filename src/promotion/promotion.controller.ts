import { Body, Controller, Get, Post, Req, UseGuards } from '@nestjs/common';
import { ApiOperation } from '@nestjs/swagger';

import { SessionAuthGuard } from '@/auth/guards/session-auth.guard';

import { AddPromotionDto } from './add-promotion.dto';
import { PromotionService } from './promotion.service';

@Controller('promotion')
export class PromotionController {
  constructor(private readonly promotionService: PromotionService) {}

  @ApiOperation({
    summary: 'Все типы продвижения',
  })
  @Get('all-promotions')
  async allPromotions() {
    return await this.promotionService.allPromotions();
  }

  @ApiOperation({
    summary: 'Подключение продвижения к товару',
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
