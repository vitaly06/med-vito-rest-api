import {
  Body,
  Controller,
  Get,
  Param,
  Post,
  Put,
  Query,
  Req,
  UseGuards,
} from '@nestjs/common';
import { ApiOperation, ApiQuery, ApiTags } from '@nestjs/swagger';

import { Request } from 'express';
import { SessionAuthGuard } from 'src/auth/guards/session-auth.guard';

import { SendReviewDto } from './dto/send-review.dto';
import { ReviewService } from './review.service';
import { AdminSessionAuthGuard } from '@/auth/guards/admin-session-auth.guard';
import { ModerateState } from './enum/moderate-state.enum';

@Controller('review')
export class ReviewController {
  constructor(private readonly reviewService: ReviewService) {}

  @ApiOperation({
    summary: 'Оставить отзыв пользователю',
  })
  @UseGuards(SessionAuthGuard)
  @Post('send-review')
  async sendReview(
    @Body() dto: SendReviewDto,
    @Req() req: Request & { user: any },
  ) {
    return await this.reviewService.sendReview(dto, req.user.id);
  }

  @ApiOperation({
    summary: 'Возвращает отзывы пользователя',
  })
  @Get('user-reviews/:id')
  async getMyReviews(@Param('id') id: string) {
    return await this.reviewService.getUserReviews(+id);
  }

  @ApiTags('Модерация отзывов')
  @ApiOperation({ summary: 'Модерация отзыва (одобрить/отклонить)' })
  @UseGuards(AdminSessionAuthGuard)
  @ApiQuery({
    name: 'status',
    enum: ModerateState,
    required: true,
    description: 'Статус модерации: APPROVED или DENIDED',
  })
  @Put('moderate-review/:id')
  async moderateReview(
    @Param('id') id: string,
    @Query('status') status: ModerateState,
  ) {
    return await this.reviewService.moderateReview(+id, status);
  }

  @ApiTags('Модерация отзывов')
  @ApiOperation({
    summary: 'Все отзывы для модерации',
  })
  @UseGuards(AdminSessionAuthGuard)
  @Get('/all-reviews-to-moderate')
  async allReviewsToModerate() {
    return await this.reviewService.allReviewsToModerate();
  }
}
