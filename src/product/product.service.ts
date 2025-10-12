import { Injectable, BadRequestException } from '@nestjs/common';
import { PrismaService } from 'src/prisma/prisma.service';
import { createProductDto } from './dto/create-product.dto';

@Injectable()
export class ProductService {
  constructor(private readonly prisma: PrismaService) {}

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
        product,
      };
    } catch (error) {
      throw new BadRequestException(
        `Ошибка при создании продукта: ${error.message}`,
      );
    }
  }
}
