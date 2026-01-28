import {
  BadRequestException,
  Injectable,
  NotFoundException,
} from '@nestjs/common';

import { PrismaService } from 'src/prisma/prisma.service';
import { S3Service } from 'src/s3/s3.service';

import { BannerPlace } from './entities/banner-place.enum';
import { ModerateState } from './enum/moderate-state.enum';

@Injectable()
export class BannerService {
  constructor(
    private readonly prisma: PrismaService,
    private readonly s3Service: S3Service,
  ) {}

  async create(
    file: Express.Multer.File,
    place: BannerPlace,
    navigateToUrl: string,
    name: string,
    userId: number,
  ) {
    // Загружаем изображение в S3
    const photoUrl = await this.s3Service.uploadFile(file, 'banners');

    // Создаём баннер в БД
    const banner = await this.prisma.banner.create({
      data: {
        name,
        photoUrl,
        place,
        navigateToUrl,
        userId,
      },
    });

    return banner;
  }

  async findAll() {
    return this.prisma.banner.findMany({
      where: {
        moderateState: 'APPROVED',
      },
      orderBy: {
        createdAt: 'desc',
      },
    });
  }

  async findByPlace(place: BannerPlace) {
    return this.prisma.banner.findMany({
      where: {
        place,
        moderateState: 'APPROVED',
      },
      orderBy: {
        createdAt: 'desc',
      },
    });
  }

  async findRandom(limit: number = 5) {
    return this.prisma.$queryRaw`
      SELECT * FROM "Banner"
      WHERE "moderateState" = 'APPROVED'
      ORDER BY RANDOM()
      LIMIT ${limit}
    `;
  }

  async findOne(id: number) {
    const banner = await this.prisma.banner.findUnique({
      where: { id },
    });

    if (!banner) {
      throw new NotFoundException(`Баннер с ID ${id} не найден`);
    }

    return banner;
  }

  async update(
    id: number,
    file?: Express.Multer.File,
    place?: BannerPlace,
    navigateToUrl?: string,
    name?: string,
  ) {
    // Проверяем существование баннера
    const existingBanner = await this.findOne(id);

    const updateData: any = {};

    // Если есть новое изображение
    if (file) {
      // Удаляем старое изображение из S3
      if (existingBanner.photoUrl) {
        await this.s3Service.deleteFile(existingBanner.photoUrl);
      }

      // Загружаем новое изображение
      const photoUrl = await this.s3Service.uploadFile(file, 'banners');
      updateData.photoUrl = photoUrl;
    }

    // Если передано место размещения
    if (place) {
      updateData.place = place;
    }

    // Если передан URL навигации
    if (navigateToUrl !== undefined) {
      updateData.navigateToUrl = navigateToUrl;
    }

    if (name) {
      updateData.name = name;
    }

    // Обновляем баннер
    const banner = await this.prisma.banner.update({
      where: { id },
      data: updateData,
    });

    return banner;
  }

  async remove(id: number) {
    // Проверяем существование баннера
    const banner = await this.findOne(id);

    // Удаляем изображение из S3
    if (banner.photoUrl) {
      await this.s3Service.deleteFile(banner.photoUrl);
    }

    // Удаляем баннер из БД
    await this.prisma.banner.delete({
      where: { id },
    });

    return { message: 'Баннер успешно удалён' };
  }

  /**
   * Зарегистрировать просмотр баннера
   */
  async registerView(bannerId: number, userId?: number, ipAddress?: string) {
    // Проверяем существование баннера
    await this.findOne(bannerId);

    // Проверяем, был ли просмотр от этого пользователя/IP за последние 24 часа
    const twentyFourHoursAgo = new Date(Date.now() - 24 * 60 * 60 * 1000);

    const orConditions: Array<{ userId?: number; ipAddress?: string }> = [];
    if (userId) orConditions.push({ userId });
    if (ipAddress) orConditions.push({ ipAddress });

    const existingView = await this.prisma.bannerView.findFirst({
      where: {
        bannerId,
        viewedAt: { gte: twentyFourHoursAgo },
        ...(orConditions.length > 0 && { OR: orConditions }),
      },
    });

    // Если уже был просмотр - не дублируем
    if (existingView) {
      return existingView;
    }

    // Создаём запись о просмотре
    return this.prisma.bannerView.create({
      data: {
        bannerId,
        userId,
        ipAddress,
      },
    });
  }

  /**
   * Получить статистику по баннеру
   */
  async getBannerStats(bannerId: number) {
    // Проверяем существование баннера
    const banner = await this.findOne(bannerId);

    // Общее количество просмотров
    const totalViews = await this.prisma.bannerView.count({
      where: { bannerId },
    });

    // Уникальные пользователи
    const uniqueUsers = await this.prisma.bannerView.groupBy({
      by: ['userId'],
      where: {
        bannerId,
        userId: { not: null },
      },
    });

    // Просмотры по дням (последние 30 дней)
    const thirtyDaysAgo = new Date();
    thirtyDaysAgo.setDate(thirtyDaysAgo.getDate() - 30);

    const viewsByDay = await this.prisma.$queryRaw`
      SELECT DATE("viewedAt") as date, COUNT(*) as views
      FROM "BannerView"
      WHERE "bannerId" = ${bannerId}
        AND "viewedAt" >= ${thirtyDaysAgo}
      GROUP BY DATE("viewedAt")
      ORDER BY date DESC
    `;

    return {
      banner,
      totalViews,
      uniqueUsers: uniqueUsers.length,
      viewsByDay,
    };
  }

  /**
   * Получить общую статистику по всем баннерам пользователя
   */
  async getUserBannersStats(userId: number) {
    const banners = await this.prisma.banner.findMany({
      where: { userId },
      include: {
        _count: {
          select: { views: true },
        },
      },
      orderBy: {
        createdAt: 'desc',
      },
    });

    return banners.map((banner) => ({
      id: banner.id,
      name: banner.name,
      place: banner.place,
      totalViews: banner._count.views,
      createdAt: banner.createdAt,
    }));
  }

  async moderateBanner(bannerId: number, status: ModerateState) {
    if (!['MODERATE', 'APPROVED', 'DENIDED'].includes(status)) {
      throw new BadRequestException(
        'Неверный статус модерации. Доступные статусы: MODERATE, APPROVED, DENIDED',
      );
    }

    const checkBanner = await this.prisma.banner.findUnique({
      where: { id: bannerId },
    });

    if (!checkBanner) {
      throw new NotFoundException('Баннер для модерации не найден');
    }

    // Обновляем баннер
    await this.prisma.banner.update({
      where: { id: bannerId },
      data: {
        moderateState: status,
      },
    });

    if (status == 'APPROVED') {
      return { message: 'Баннер успешно опубликован' };
    } else if (status == 'DENIDED') {
      return { message: 'Баннер успешно отклонён' };
    }
    return { message: 'Баннер возвращён на модерацию' };
  }

  async allBannersToModerate() {
    const banners = await this.prisma.banner.findMany({
      where: { moderateState: 'MODERATE' },
      select: {
        id: true,
        name: true,
        photoUrl: true,
        place: true,
        navigateToUrl: true,
        createdAt: true,
        user: {
          select: {
            id: true,
            fullName: true,
            photo: true,
          },
        },
      },
      orderBy: {
        createdAt: 'desc',
      },
    });

    return banners;
  }
}
