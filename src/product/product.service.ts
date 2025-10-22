import { Injectable, BadRequestException } from '@nestjs/common';
import { PrismaService } from 'src/prisma/prisma.service';
import { createProductDto } from './dto/create-product.dto';
import { ConfigService } from '@nestjs/config';

@Injectable()
export class ProductService {
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

  async createProduct(
    dto: createProductDto,
    fileNames: string[],
    userId: number,
  ) {
    try {
      const category = await this.prisma.category.findUnique({
        where: { id: +dto.categoryId },
        include: {
          subCategories: {
            where: { id: +dto.subcategoryId },
          },
        },
      });

      if (!category) {
        throw new BadRequestException('Категория не найдена');
      }

      if (category.subCategories.length === 0) {
        throw new BadRequestException(
          'Подкатегория не найдена или не принадлежит указанной категории',
        );
      }

      // Создаем массив путей к изображениям
      const imagePaths = fileNames.map(
        (fileName) => `/uploads/product/${fileName}`,
      );

      const product = await this.prisma.product.create({
        data: {
          name: dto.name,
          price: +dto.price,
          state: dto.state,
          brand: dto.brand,
          model: dto.model,
          description: dto.description,
          address: dto.address,
          images: imagePaths,
          categoryId: +dto.categoryId,
          subCategoryId: +dto.subcategoryId,
          userId: userId,
        },
        include: {
          category: true,
          subCategory: true,
          user: {
            select: {
              id: true,
              fullName: true,
              email: true,
              rating: true,
            },
          },
        },
      });

      return {
        message: 'Продукт успешно создан',
        product: {
          ...product,
          images: imagePaths.map((path) => `${this.baseUrl}${path}`),
        },
      };
    } catch (error) {
      throw new BadRequestException(
        `Ошибка при создании продукта: ${error.message}`,
      );
    }
  }

  async findAll() {
    const products = await this.prisma.product.findMany({
      select: {
        id: true,
        images: true,
        name: true,
        address: true,
        createdAt: true,
        price: true,
      },
    });

    return this.formatProductsResponse(products);
  }

  async getProductsByUserId(id: number) {
    const checkUser = await this.prisma.user.findUnique({
      where: { id },
    });

    if (!checkUser) {
      throw new BadRequestException('Пользователь не найден');
    }

    const products = await this.prisma.product.findMany({
      where: { userId: id },
      select: {
        id: true,
        images: true,
        name: true,
        address: true,
        createdAt: true,
        price: true,
      },
    });

    return this.formatProductsResponse(products);
  }

  private formatProductsResponse(products: any[]) {
    return products.map((product) => {
      const createdAt = new Date(product.createdAt);
      const day = createdAt.getDate().toString().padStart(2, '0');
      const month = (createdAt.getMonth() + 1).toString().padStart(2, '0');
      const year = createdAt.getFullYear().toString().slice(-2);
      const hours = createdAt.getHours().toString().padStart(2, '0');
      const minutes = createdAt.getMinutes().toString().padStart(2, '0');

      return {
        ...product,
        images: product.images.map((img) => `${this.baseUrl}${img}`),
        createdAt: `${day}.${month}.${year} в ${hours}:${minutes}`,
      };
    });
  }

  async addProductToFavorites(id: number, userId: number) {
    const checkProduct = await this.prisma.product.findUnique({
      where: { id },
    });

    if (!checkProduct) {
      throw new BadRequestException('Товар не найден');
    }

    const existingFavorite = await this.prisma.user.findFirst({
      where: {
        id: userId,
        favorites: {
          some: { id },
        },
      },
    });

    if (existingFavorite) {
      throw new BadRequestException('Товар уже добавлен в избранное');
    }

    await this.prisma.user.update({
      where: { id: userId },
      data: {
        favorites: {
          connect: { id },
        },
      },
    });

    return { message: 'Товар успешно добавлен в избранное' };
  }

  async removeProductFromFavorites(id: number, userId: number) {
    const checkProduct = await this.prisma.product.findUnique({
      where: { id },
    });

    if (!checkProduct) {
      throw new BadRequestException('Товар не найден');
    }

    const existingFavorite = await this.prisma.user.findFirst({
      where: {
        id: userId,
        favorites: {
          some: { id },
        },
      },
    });

    if (!existingFavorite) {
      throw new BadRequestException('Товар не в избранном');
    }

    await this.prisma.user.update({
      where: { id: userId },
      data: {
        favorites: {
          disconnect: { id },
        },
      },
    });

    return { message: 'Товар успешно удалён из избранного' };
  }

  async getMyFavorites(userId: number) {
    const checkUser = await this.prisma.user.findUnique({
      where: { id: userId },
      select: {
        favorites: {
          select: {
            id: true,
            images: true,
            name: true,
            address: true,
            createdAt: true,
            price: true,
          },
        },
      },
    });

    return this.formatProductsResponse(checkUser?.favorites || []);
  }

  async getProductCard(id: number) {
    const product = await this.prisma.product.findUnique({
      where: {
        id,
      },
      select: {
        id: true,
        name: true,
        price: true,
        images: true,
        address: true,
        brand: true,
        model: true,
      },
    });

    return {
      ...product,
      images: product?.images.map((img) => `${this.baseUrl}${img}`),
    };
  }
}
