import {
  BadRequestException,
  Injectable,
  NotFoundException,
} from '@nestjs/common';
import { PrismaService } from 'src/prisma/prisma.service';

@Injectable()
export class UserService {
  profileTypes = {
    IP: 'Индивидуальный предприниматель',
    OOO: 'Юридическое лицо',
    INDIVIDUAL: 'Физическое лицо',
  };
  constructor(private readonly prisma: PrismaService) {}

  async getUserInfo(id: number) {
    const checkUser = await this.prisma.user.findUnique({
      where: { id },
      select: {
        id: true,
        fullName: true,
        profileType: true,
        phoneNumber: true,
      },
    });

    const reviews = await this.prisma.review.findMany({
      where: { reviewedUserId: id },
      select: {
        rating: true,
      },
    });

    if (!checkUser) {
      throw new BadRequestException('Пользователь не найден');
    }

    return {
      ...checkUser,
      profileType: this.profileTypes[checkUser.profileType],
      rating:
        reviews.reduce((sum, review) => sum + review.rating, 0) /
        reviews.length,
      reviewsCount: reviews.length,
    };
  }

  async getPhoneNumber(userId: number, myId: number) {
    const checkUser = await this.prisma.user.findUnique({
      where: { id: userId },
    });

    if (!checkUser) {
      throw new NotFoundException('Продавец с данным id не найден');
    }

    if (!checkUser.phoneNumber) {
      throw new BadRequestException('У продавца отсутствует номер телефона');
    }

    // Проверяем, не смотрел ли уже этот пользователь номер этого продавца
    const existingView = await this.prisma.phoneNumberView.findUnique({
      where: {
        viewedById_viewedUserId: {
          viewedById: myId,
          viewedUserId: userId,
        },
      },
    });

    // Если ещё не смотрел, записываем в статистику
    if (!existingView) {
      await this.prisma.phoneNumberView.create({
        data: {
          viewedById: myId,
          viewedUserId: userId,
        },
      });
    }

    return {
      phoneNumber: checkUser.phoneNumber,
    };
  }

  // Получить статистику просмотров номеров
  // async getPhoneNumberViewStats(
  //   userId: number,
  //   page: number = 1,
  //   limit: number = 20,
  // ) {
  //   // Статистика кто смотрел номер этого пользователя
  //   const viewsReceived = await this.prisma.phoneNumberView.findMany({
  //     where: { viewedUserId: userId },
  //     include: {
  //       viewedBy: {
  //         select: {
  //           id: true,
  //           fullName: true,
  //           profileType: true,
  //         },
  //       },
  //     },
  //     orderBy: { viewedAt: 'desc' },
  //     skip: (page - 1) * limit,
  //     take: limit,
  //   });

  //   // Подсчитываем общее количество
  //   const totalViews = await this.prisma.phoneNumberView.count({
  //     where: { viewedUserId: userId },
  //   });

  //   return {
  //     views: viewsReceived.map((view) => ({
  //       id: view.id,
  //       viewedBy: {
  //         id: view.viewedBy.id,
  //         fullName: view.viewedBy.fullName,
  //         profileType: this.profileTypes[view.viewedBy.profileType],
  //       },
  //       viewedAt: view.viewedAt,
  //       userAgent: view.userAgent,
  //       ipAddress: view.ipAddress,
  //     })),
  //     pagination: {
  //       page,
  //       limit,
  //       total: totalViews,
  //       pages: Math.ceil(totalViews / limit),
  //     },
  //   };
  // }
}
