import { BadRequestException, Injectable } from '@nestjs/common';

import { PrismaService } from 'src/prisma/prisma.service';

import { SendReviewDto } from './dto/send-review.dto';

@Injectable()
export class ReviewService {
  constructor(private readonly prisma: PrismaService) {}

  async sendReview(dto: SendReviewDto, userId: number) {
    const checkUser = await this.prisma.user.findUnique({
      where: { id: dto.reviewedUserId },
    });

    if (!checkUser) {
      throw new BadRequestException('Продавец с таким id не найден');
    }

    if (userId == dto.reviewedUserId) {
      throw new BadRequestException('Нельзя оставить отзыв самому себе');
    }

    await this.prisma.review.create({
      data: {
        ...dto,
        reviewedById: userId,
      },
    });

    return { message: 'Отзыв успешно оставлен' };
  }

  async getUserReviews(userId: number) {
    const reviews = await this.prisma.review.findMany({
      where: { reviewedUserId: userId },
      select: {
        rating: true,
        text: true,
        reviewedAt: true,
        reviewedBy: {
          select: {
            id: true,
            fullName: true,
          },
        },
      },
    });
    // Итоговый рейтинг пользователя
    const totalRating =
      reviews.reduce((sum, elem) => sum + elem.rating, 0) / reviews.length;

    return {
      totalRating: Math.round(totalRating * 100) / 100 || 0,
      reviewsCount: reviews.length,
      reviews: reviews.map((review) => {
        const createdAt = new Date(review.reviewedAt);
        const day = createdAt.getDate().toString().padStart(2, '0');
        const month = (createdAt.getMonth() + 1).toString().padStart(2, '0');
        const year = createdAt.getFullYear().toString().slice(-2);
        return {
          rating: review.rating,
          text: review.text,
          date: `${day}.${month}.${year}`,
          fullName: review.reviewedBy.fullName,
        };
      }),
    };
  }
}
