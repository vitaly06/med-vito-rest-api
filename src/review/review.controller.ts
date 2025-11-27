import { Body, Controller, Get, Post, Req, UseGuards } from '@nestjs/common';
import { ReviewService } from './review.service';
import { SendReviewDto } from './dto/send-review.dto';
import { SessionAuthGuard } from 'src/auth/guards/session-auth.guard';
import { ApiOperation } from '@nestjs/swagger';
import { Request } from 'express';

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
    summary: 'Мои отзывы',
  })
  @UseGuards(SessionAuthGuard)
  @Get('my-reviews')
  async getMyReviews(@Req() req: Request & { user: any }) {
    return await this.reviewService.getMyReviews(req.user.id);
  }
}
