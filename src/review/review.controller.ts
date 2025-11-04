import { Body, Controller, Get, Post, Req, UseGuards } from '@nestjs/common';
import { ReviewService } from './review.service';
import { SendReviewDto } from './dto/send-review.dto';
import { JwtAuthGuard } from 'src/auth/guards/jwt-auth.guard';
import { ApiOperation } from '@nestjs/swagger';
import { Request } from 'express';

@Controller('review')
export class ReviewController {
  constructor(private readonly reviewService: ReviewService) {}

  @ApiOperation({
    summary: 'Оставить отзыв пользователю',
  })
  @UseGuards(JwtAuthGuard)
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
  @UseGuards(JwtAuthGuard)
  @Get('my-reviews')
  async getMyReviews(@Req() req: Request & { user: any }) {
    return await this.reviewService.getMyReviews(req.user.id);
  }
}
