import { Injectable } from '@nestjs/common';
import { ConfigService } from '@nestjs/config';

import { PrismaService } from '@/prisma/prisma.service';

export interface StatisticsQueryDto {
  period?: 'day' | 'week' | 'month' | 'quarter' | 'half-year' | 'year';
  categoryId?: number;
  region?: string;
  productId?: number;
}

@Injectable()
export class StatisticsService {
  baseUrl: string;

  constructor(
    private readonly prisma: PrismaService,
    private readonly configService: ConfigService,
  ) {
    this.baseUrl = this.configService.get<string>(
      'BASE_URL',
      'http://localhost:3000',
    );
  }

  // Получить статистику для пользователя
  async getUserStatistics(userId: number, query: StatisticsQueryDto) {
    const dateFilter = this.getDateFilter(query.period);

    // Базовые условия для товаров пользователя
    const baseProductFilter = {
      userId,
      ...(query.categoryId && {
        subCategory: {
          categoryId: query.categoryId,
        },
      }),
      ...(query.region && {
        address: {
          contains: query.region,
          mode: 'insensitive' as const,
        },
      }),
      ...(query.productId && { id: query.productId }),
    };

    // Получаем статистику просмотров товаров
    const viewsStats = await this.getProductViewsStats(
      baseProductFilter,
      dateFilter,
    );

    // Получаем статистику показов телефонов
    const phoneStats = await this.getPhoneViewsStats(userId, dateFilter);

    // Получаем статистику добавлений в избранное
    const favoriteStats = await this.getFavoriteStats(
      baseProductFilter,
      dateFilter,
    );

    return {
      period: query.period || 'all-time',
      totalViews: viewsStats.total,
      totalPhoneViews: phoneStats.total,
      totalFavorites: favoriteStats.total,
    };
  }

  // Получить статистику просмотров товаров
  private async getProductViewsStats(productFilter: any, dateFilter: any) {
    const views = await this.prisma.productView.groupBy({
      by: ['productId'],
      where: {
        product: productFilter,
        ...dateFilter,
      },
      _count: {
        id: true,
      },
    });

    const total = views.reduce((sum, view) => sum + view._count.id, 0);

    const breakdown = await Promise.all(
      views.map(async (view) => {
        const product = await this.prisma.product.findUnique({
          where: { id: view.productId },
          select: {
            id: true,
            name: true,
            images: true,
          },
        });

        return {
          productId: view.productId,
          productName: product?.name || 'Неизвестный товар',
          productImage: product?.images[0]
            ? `${this.baseUrl}${product.images[0]}`
            : null,
          count: view._count.id,
        };
      }),
    );

    return { total, breakdown };
  }

  // Получить статистику показов телефонов
  private async getPhoneViewsStats(userId: number, dateFilter: any) {
    const phoneViews = await this.prisma.phoneNumberView.findMany({
      where: {
        viewedUserId: userId,
        ...dateFilter,
      },
      include: {
        viewedBy: {
          select: {
            id: true,
            fullName: true,
            email: true,
          },
        },
      },
    });

    const total = phoneViews.length;

    // Группируем по дням для графика
    const breakdown = phoneViews.reduce((acc: any[], view) => {
      const date = view.viewedAt.toISOString().split('T')[0];
      const existing = acc.find((item) => item.date === date);

      if (existing) {
        existing.count++;
      } else {
        acc.push({
          date,
          count: 1,
        });
      }

      return acc;
    }, []);

    return { total, breakdown };
  }

  // Получить статистику добавлений в избранное
  private async getFavoriteStats(productFilter: any, dateFilter: any) {
    // Создаем фильтр с правильным полем для FavoriteAction (addedAt вместо viewedAt)
    const favoriteFilter = dateFilter.viewedAt
      ? { addedAt: dateFilter.viewedAt }
      : {};

    const favorites = await this.prisma.favoriteAction.groupBy({
      by: ['productId'],
      where: {
        product: productFilter,
        ...favoriteFilter,
      },
      _count: {
        id: true,
      },
    });

    const total = favorites.reduce((sum, fav) => sum + fav._count.id, 0);

    const breakdown = await Promise.all(
      favorites.map(async (favorite) => {
        const product = await this.prisma.product.findUnique({
          where: { id: favorite.productId },
          select: {
            id: true,
            name: true,
            images: true,
          },
        });

        return {
          productId: favorite.productId,
          productName: product?.name || 'Неизвестный товар',
          productImage: product?.images[0]
            ? `${this.baseUrl}${product.images[0]}`
            : null,
          count: favorite._count.id,
        };
      }),
    );

    return { total, breakdown };
  }

  // Получить детальную информацию по товарам
  private async getProductsDetails(productFilter: any) {
    const products = await this.prisma.product.findMany({
      where: productFilter,
      select: {
        id: true,
        name: true,
        price: true,
        images: true,
        address: true,
        createdAt: true,
        subCategory: {
          select: {
            name: true,
            category: {
              select: {
                name: true,
              },
            },
          },
        },
        _count: {
          select: {
            views: true,
            favoriteActions: true,
          },
        },
      },
      orderBy: {
        createdAt: 'desc',
      },
    });

    return products.map((product) => ({
      id: product.id,
      name: product.name,
      price: product.price,
      image: product.images[0] ? `${this.baseUrl}${product.images[0]}` : null,
      address: product.address,
      category: product.subCategory?.category?.name || null,
      subCategory: product.subCategory?.name || null,
      createdAt: product.createdAt,
      stats: {
        views: product._count.views,
        favorites: product._count.favoriteActions,
      },
    }));
  }

  // Вспомогательный метод для создания фильтра по датам
  private getDateFilter(period?: string) {
    if (!period) return {};

    const now = new Date();
    let startDate: Date;

    switch (period) {
      case 'day':
        startDate = new Date(now.getFullYear(), now.getMonth(), now.getDate());
        break;
      case 'week':
        startDate = new Date(now.getTime() - 7 * 24 * 60 * 60 * 1000);
        break;
      case 'month':
        startDate = new Date(now.getFullYear(), now.getMonth(), 1);
        break;
      case 'quarter':
        const quarterStart = Math.floor(now.getMonth() / 3) * 3;
        startDate = new Date(now.getFullYear(), quarterStart, 1);
        break;
      case 'half-year':
        const halfYearStart = now.getMonth() < 6 ? 0 : 6;
        startDate = new Date(now.getFullYear(), halfYearStart, 1);
        break;
      case 'year':
        startDate = new Date(now.getFullYear(), 0, 1);
        break;
      default:
        return {};
    }

    return {
      viewedAt: {
        gte: startDate,
      },
    };
  }

  async getProductsAnalytic(userId: number) {
    const products = await this.prisma.product.findMany({
      where: { userId },
      select: {
        id: true,
        images: true,
        name: true,
        price: true,
        views: true,
        favoritedBy: true,
        user: {
          select: {
            phoneNumberViews: true,
          },
        },
      },
    });

    return products.map((product) => {
      return {
        id: product.id,
        image: product.images[0],
        name: product.name,
        price: product.price,
        views: product.views.length,
        favoritedBy: product.favoritedBy.length,
        phoneNumberViews: product.user.phoneNumberViews.length,
      };
    });
  }
}
