import { Injectable, BadRequestException } from '@nestjs/common';
import { PrismaService } from 'src/prisma/prisma.service';
import { createProductDto } from './dto/create-product.dto';
import { ConfigService } from '@nestjs/config';
import { UserService } from 'src/user/user.service';

@Injectable()
export class ProductService {
  baseUrl: string;
  constructor(
    private readonly prisma: PrismaService,
    private readonly configService: ConfigService,
    private readonly userService: UserService,
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

      // Валидация fieldValues - проверяем что все поля существуют
      if (dto.fieldValues) {
        const fieldIds = Object.keys(dto.fieldValues).map((id) => parseInt(id));

        if (fieldIds.length > 0) {
          const fields = await this.prisma.typeField.findMany({
            where: {
              id: { in: fieldIds },
            },
          });

          if (fields.length !== fieldIds.length) {
            const foundIds = fields.map((f) => f.id);
            const missingIds = fieldIds.filter((id) => !foundIds.includes(id));
            throw new BadRequestException(
              `Поля с ID ${missingIds.join(', ')} не найдены. Проверьте fieldValues.`,
            );
          }
        }
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
          description: dto.description,
          address: dto.address,
          images: imagePaths,
          categoryId: +dto.categoryId,
          subCategoryId: +dto.subcategoryId,
          typeId: dto.typeId ? +dto.typeId : null,
          userId: userId,
          // Создаем связанные значения полей, если они переданы
          fieldValues: dto.fieldValues
            ? {
                create: Object.entries(dto.fieldValues).map(
                  ([fieldId, value]) => ({
                    fieldId: parseInt(fieldId),
                    value: value,
                  }),
                ),
              }
            : undefined,
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
          fieldValues: {
            include: {
              field: true,
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

  async findAll(userId?: number) {
    let products = await this.prisma.product.findMany({
      select: {
        id: true,
        images: true,
        name: true,
        address: true,
        createdAt: true,
        price: true,
        userId: true,
      },
    });

    if (userId) {
      const result = await Promise.all(
        products.map(async (product) => {
          if (product.userId != userId) {
            return {
              ...product,
              isFavorited: await this.isProductInUserFavorites(
                product.id,
                userId,
              ),
            };
          }
          return {
            ...product,
            isFavorited: false,
          };
        }),
      );
      return this.formatProductsResponse(result);
    }

    return this.formatProductsResponse(
      products.map((product) => {
        return {
          ...product,
          isFavorited: false,
        };
      }),
    );
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
        images: product.images
          ? product.images.map((img) => `${this.baseUrl}${img}`)
          : [],
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

    if (checkProduct.userId == userId) {
      throw new BadRequestException('Нельзя добавить свой товар в избранное');
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

    // Добавляем товар в избранное
    await this.prisma.user.update({
      where: { id: userId },
      data: {
        favorites: {
          connect: { id },
        },
      },
    });

    // Записываем статистику добавления в избранное
    try {
      await this.prisma.favoriteAction.create({
        data: {
          userId,
          productId: id,
        },
      });
    } catch (error) {
      // Игнорируем ошибки записи статистики, чтобы не блокировать основной функционал
      console.log('Error recording favorite action:', error.message);
    }

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

  async getProductCard(id: number, userId?: number) {
    const product = await this.prisma.product.findUnique({
      where: {
        id,
      },
      select: {
        id: true,
        name: true,
        description: true,
        price: true,
        images: true,
        address: true,
        userId: true, // Добавляем для проверки, что пользователь не смотрит свой товар
        fieldValues: {
          include: {
            field: true,
          },
        },
      },
    });

    if (!product) {
      throw new BadRequestException('Товар не найден');
    }

    // Если пользователь авторизован и это не его товар - записываем просмотр
    if (userId && userId !== product.userId) {
      try {
        await this.prisma.productView.upsert({
          where: {
            viewedById_productId: {
              viewedById: userId,
              productId: id,
            },
          },
          update: {
            viewedAt: new Date(),
          },
          create: {
            viewedById: userId,
            productId: id,
          },
        });
      } catch (error) {
        // Игнорируем ошибки записи просмотров, чтобы не блокировать основной функционал
        console.log('Error recording product view:', error.message);
      }
    }

    return {
      ...product,
      images: product?.images.map((img) => `${this.baseUrl}${img}`),
      isFavorited: userId
        ? await this.isProductInUserFavorites(product.id, userId)
        : false,
      seller: await this.userService.getUserInfo(product.userId),
    };
  }

  async searchProducts(searchDto: any, userId?: number) {
    const {
      search,
      categoryId,
      subCategoryId,
      minPrice,
      maxPrice,
      state,
      region,
      sortBy = 'relevance',
      page = 1,
      limit = 20,
    } = searchDto;

    // Строим условия для фильтрации
    const whereConditions: any = {};

    // Поиск по тексту (название, описание)
    if (search) {
      whereConditions.OR = [
        {
          name: {
            contains: search,
            mode: 'insensitive',
          },
        },
        {
          description: {
            contains: search,
            mode: 'insensitive',
          },
        },
      ];
    }

    // Фильтр по категории
    if (categoryId) {
      whereConditions.categoryId = categoryId;
    }

    // Фильтр по подкатегории
    if (subCategoryId) {
      whereConditions.subCategoryId = subCategoryId;
    }

    // Фильтр по цене
    if (minPrice || maxPrice) {
      whereConditions.price = {};
      if (minPrice) whereConditions.price.gte = minPrice;
      if (maxPrice) whereConditions.price.lte = maxPrice;
    }

    // Фильтр по состоянию
    if (state) {
      whereConditions.state = state;
    }

    // Фильтр по региону
    if (region) {
      whereConditions.address = {
        contains: region,
        mode: 'insensitive',
      };
    }

    // Настройка сортировки
    let orderBy: any = {};
    switch (sortBy) {
      case 'price_asc':
        orderBy = { price: 'asc' };
        break;
      case 'price_desc':
        orderBy = { price: 'desc' };
        break;
      case 'date_desc':
        orderBy = { createdAt: 'desc' };
        break;
      case 'date_asc':
        orderBy = { createdAt: 'asc' };
        break;
      case 'relevance':
      default:
        // При поиске по тексту - сортировка по релевантности (сначала точные совпадения в названии)
        if (search) {
          orderBy = [
            {
              name: {
                _relevance: {
                  search: search,
                  sort: 'desc',
                },
              },
            },
            { createdAt: 'desc' },
          ];
        } else {
          orderBy = { createdAt: 'desc' };
        }
        break;
    }

    // Выполняем поиск с пагинацией
    const [products, total] = await Promise.all([
      this.prisma.product.findMany({
        where: whereConditions,
        include: {
          category: {
            select: {
              id: true,
              name: true,
            },
          },
          subCategory: {
            select: {
              id: true,
              name: true,
            },
          },
          user: {
            select: {
              id: true,
              fullName: true,
              rating: true,
            },
          },
          _count: {
            select: {
              views: true,
              favoriteActions: true,
            },
          },
        },
        orderBy,
        skip: (page - 1) * limit,
        take: limit,
      }),
      this.prisma.product.count({
        where: whereConditions,
      }),
    ]);

    // Форматируем результат
    const formattedProducts = await Promise.all(
      products.map(async (product) => ({
        id: product.id,
        name: product.name,
        price: product.price,
        state: product.state,
        description: product.description,
        address: product.address,
        images: product.images.map((img) => `${this.baseUrl}${img}`),
        category: product.category,
        subCategory: product.subCategory,
        seller: {
          id: product.user.id,
          fullName: product.user.fullName,
          rating: product.user.rating,
        },
        stats: {
          views: product._count.views,
          favorites: product._count.favoriteActions,
        },
        createdAt: product.createdAt,
        // Отмечаем, добавлен ли в избранное текущим пользователем (если авторизован)
        isInFavorites: userId
          ? await this.isProductInUserFavorites(product.id, userId)
          : false,
      })),
    );

    return {
      products: formattedProducts,
      pagination: {
        page,
        limit,
        total,
        pages: Math.ceil(total / limit),
      },
      filters: {
        search: search || null,
        categoryId: categoryId || null,
        subCategoryId: subCategoryId || null,
        minPrice: minPrice || null,
        maxPrice: maxPrice || null,
        state: state || null,
        region: region || null,
        sortBy,
      },
    };
  }

  // Вспомогательный метод для проверки, добавлен ли товар в избранное
  private async isProductInUserFavorites(
    productId: number,
    userId: number,
  ): Promise<boolean> {
    const favorite = await this.prisma.user.findFirst({
      where: {
        id: userId,
        favorites: {
          some: { id: productId },
        },
      },
    });

    return !!favorite;
  }
}
