import {
  BadRequestException,
  Injectable,
  NotFoundException,
} from '@nestjs/common';

import { PrismaService } from 'src/prisma/prisma.service';

import { SendReviewDto } from './dto/send-review.dto';
import { ModerateState } from './enum/moderate-state.enum';

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
      where: { reviewedUserId: userId, moderateState: 'APPROVED' },
      select: {
        rating: true,
        text: true,
        createdAt: true,
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
        const createdAt = new Date(review.createdAt);
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

  async moderateReview(reviewId: number, status: ModerateState) {
    if (!['APPROVED', 'DENIDED'].includes(status)) {
      throw new BadRequestException(
        'Неверный статус модерации. Доступные статуты: APPROVED, DENIDED',
      );
    }

    const checkReview = await this.prisma.review.findUnique({
      where: { id: reviewId },
    });

    if (!checkReview) {
      throw new NotFoundException('Отзыв для модерации не найден');
    }

    // Обновляем товар
    await this.prisma.review.update({
      where: { id: reviewId },
      data: {
        moderateState: status,
      },
    });

    if (status == 'APPROVED') {
      return { message: 'Отзыв успешно опубликован' };
    }
    return { message: 'Отзыв успешно отклонён' };
  }

  async allReviewsToModerate() {
    const reviews = await this.prisma.review.findMany({
      where: { moderateState: 'MODERATE' },
      select: {
        id: true,
        rating: true,
        text: true,
        createdAt: true,
        reviewedBy: {
          select: {
            id: true,
            fullName: true,
          },
        },
        reviewedUser: {
          select: {
            id: true,
            fullName: true,
          },
        },
      },
      orderBy: {
        createdAt: 'desc',
      },
    });

    return reviews;
  }
}
