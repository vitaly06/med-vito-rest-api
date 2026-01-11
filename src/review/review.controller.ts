import {
  Body,
  Controller,
  Get,
  Param,
  Post,
  Req,
  UseGuards,
} from '@nestjs/common';
import { ApiOperation } from '@nestjs/swagger';

import { Request } from 'express';
import { SessionAuthGuard } from 'src/auth/guards/session-auth.guard';

import { SendReviewDto } from './dto/send-review.dto';
import { ReviewService } from './review.service';

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
}
