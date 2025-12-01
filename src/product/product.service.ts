import {
  Injectable,
  BadRequestException,
  ForbiddenException,
} from '@nestjs/common';
import { PrismaService } from 'src/prisma/prisma.service';
import { createProductDto } from './dto/create-product.dto';
import { UpdateProductDto } from './dto/update-product.dto';
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
        console.log('fieldValues type:', typeof dto.fieldValues);
        console.log('fieldValues value:', dto.fieldValues);
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
          videoUrl: dto.videoUrl || null,
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

  async deleteProduct(productId, userId: number) {
    const checkProduct = await this.prisma.product.findUnique({
      where: { id: productId },
      select: {
        userId: true,
      },
    });

    if (!checkProduct) {
      throw new BadRequestException('Товар для удаления не найден');
    }

    if (checkProduct.userId != userId) {
      throw new ForbiddenException('Вы не можете удалить чужой товар');
    }

    await this.prisma.product.delete({
      where: { id: productId },
    });

    return { message: 'Товар успешно удалён' };
  }

  async updateProduct(
    productId: number,
    dto: UpdateProductDto,
    fileNames: string[],
    userId: number,
  ) {
    try {
      // Проверяем, что товар существует и принадлежит пользователю
      const existingProduct = await this.prisma.product.findUnique({
        where: { id: productId },
        include: {
          fieldValues: true,
          type: {
            include: {
              fields: true,
            },
          },
        },
      });

      if (!existingProduct) {
        throw new BadRequestException('Товар не найден');
      }

      if (existingProduct.userId !== userId) {
        throw new ForbiddenException('Вы не можете редактировать чужой товар');
      }

      // Валидация fieldValues - проверяем что все поля существуют и принадлежат типу товара
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

          // Проверяем, что все поля принадлежат типу этого товара
          if (existingProduct.typeId) {
            const invalidFields = fields.filter(
              (f) => f.typeId !== existingProduct.typeId,
            );
            if (invalidFields.length > 0) {
              throw new BadRequestException(
                `Поля ${invalidFields.map((f) => f.name).join(', ')} не принадлежат типу этого товара`,
              );
            }
          }
        }
      }

      // Подготовка данных для обновления
      const updateData: any = {};

      if (dto.name !== undefined) updateData.name = dto.name;
      if (dto.price !== undefined) updateData.price = +dto.price;
      if (dto.state !== undefined) updateData.state = dto.state;
      if (dto.description !== undefined)
        updateData.description = dto.description;
      if (dto.address !== undefined) updateData.address = dto.address;
      if (dto.videoUrl !== undefined) updateData.videoUrl = dto.videoUrl;

      // Обработка новых изображений
      if (fileNames && fileNames.length > 0) {
        const newImagePaths = fileNames.map(
          (fileName) => `/uploads/product/${fileName}`,
        );
        // Добавляем новые изображения к существующим
        updateData.images = [...existingProduct.images, ...newImagePaths];
      }

      // Обновляем товар
      const product = await this.prisma.product.update({
        where: { id: productId },
        data: updateData,
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

      // Обработка fieldValues - добавляем или обновляем
      if (dto.fieldValues) {
        for (const [fieldId, value] of Object.entries(dto.fieldValues)) {
          await this.prisma.productFieldValue.upsert({
            where: {
              fieldId_productId: {
                fieldId: parseInt(fieldId),
                productId: productId,
              },
            },
            update: {
              value: value,
            },
            create: {
              fieldId: parseInt(fieldId),
              productId: productId,
              value: value,
            },
          });
        }
      }

      // Получаем обновленный товар с новыми fieldValues
      const updatedProduct = await this.prisma.product.findUnique({
        where: { id: productId },
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
        message: 'Товар успешно обновлён',
        product: {
          ...updatedProduct!,
          images: updatedProduct!.images.map(
            (path) => `${this.baseUrl}${path}`,
          ),
        },
      };
    } catch (error) {
      throw new BadRequestException(
        `Ошибка при обновлении товара: ${error.message}`,
      );
    }
  }

  async findAll(userId?: number, searchDto?: any) {
    // Если есть параметры поиска, используем расширенный поиск
    if (
      searchDto &&
      Object.keys(searchDto).length > 0 &&
      Object.values(searchDto).some((val) => val !== undefined && val !== '')
    ) {
      const {
        search,
        categoryId,
        subCategoryId,
        typeId,
        minPrice,
        maxPrice,
        state,
        region,
        profileType,
        fieldValues,
        sortBy = 'date_desc',
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
        whereConditions.categoryId = parseInt(categoryId);
      }

      // Фильтр по подкатегории
      if (subCategoryId) {
        whereConditions.subCategoryId = parseInt(subCategoryId);
      }

      // Фильтр по типу подкатегории
      if (typeId) {
        whereConditions.typeId = parseInt(typeId);
      }

      // Фильтр по цене
      if (minPrice || maxPrice) {
        whereConditions.price = {};
        if (minPrice) whereConditions.price.gte = parseInt(minPrice);
        if (maxPrice) whereConditions.price.lte = parseInt(maxPrice);
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

      // Фильтр по типу продавца
      if (profileType) {
        whereConditions.user = {
          profileType: profileType,
        };
      }

      // Фильтр по характеристикам товара (fieldValues)
      if (fieldValues && Object.keys(fieldValues).length > 0) {
        whereConditions.fieldValues = {
          some: {
            OR: Object.entries(fieldValues).map(([fieldId, value]) => ({
              fieldId: parseInt(fieldId),
              value: {
                contains: value as string,
                mode: 'insensitive',
              },
            })),
          },
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
          orderBy = { createdAt: 'desc' };
          break;
      }

      // Выполняем поиск
      const products = await this.prisma.product.findMany({
        where: whereConditions,
        select: {
          id: true,
          images: true,
          name: true,
          address: true,
          createdAt: true,
          price: true,
          userId: true,
          videoUrl: true,
        },
        orderBy,
        skip: (page - 1) * limit,
        take: limit,
      });

      // Форматируем результат как в обычном findAll
      if (userId) {
        const result = await Promise.all(
          products.map(async (product) => ({
            ...product,
            isFavorited: await this.isProductInUserFavorites(
              product.id,
              userId,
            ),
          })),
        );
        return this.formatProductsResponse(result);
      }

      return this.formatProductsResponse(
        products.map((product) => ({
          ...product,
          isFavorited: false,
        })),
      );
    }

    // Если параметров поиска нет - возвращаем все товары + 5 случайных товаров из разных подкатегорий

    // Получаем все товары
    const allProducts = await this.prisma.product.findMany({
      select: {
        id: true,
        images: true,
        name: true,
        address: true,
        createdAt: true,
        price: true,
        userId: true,
        videoUrl: true,
      },
      orderBy: { createdAt: 'desc' },
    });

    // Форматируем оба списка
    if (userId) {
      const allProductsWithFavorites = await Promise.all(
        allProducts.map(async (product) => ({
          ...product,
          isFavorited: await this.isProductInUserFavorites(product.id, userId),
        })),
      );

      return this.formatProductsResponse(allProductsWithFavorites);
    }

    return this.formatProductsResponse(
      allProducts.map((product) => ({
        ...product,
        isFavorited: false,
      })),
    );
  }

  async getRandomProducts(userId?: number) {
    // Получаем 5 случайных подкатегорий, у которых есть товары
    const randomSubCategories = await this.prisma.$queryRaw<
      Array<{ id: number; name: string }>
    >`
      SELECT sc.id, sc.name
      FROM "SubCategory" sc
      WHERE EXISTS (
        SELECT 1 FROM "Product" p WHERE p."subCategoryId" = sc.id
      )
      ORDER BY RANDOM()
      LIMIT 5
    `;

    const randomProductsList: any[] = [];

    // Для каждой подкатегории берём один случайный товар
    for (const subCat of randomSubCategories) {
      const products = await this.prisma.product.findMany({
        where: { subCategoryId: subCat.id },
        select: {
          id: true,
          images: true,
          name: true,
          address: true,
          createdAt: true,
          price: true,
          userId: true,
          videoUrl: true,
        },
      });

      if (products.length > 0) {
        // Берём случайный товар из этой подкатегории
        const randomIndex = Math.floor(Math.random() * products.length);
        randomProductsList.push({
          ...products[randomIndex],
          subCategoryName: subCat.name,
        });
      }
    }

    if (userId) {
      const randomProductsWithFavorites = await Promise.all(
        randomProductsList.map(async (product) => ({
          ...product,
          isFavorited: await this.isProductInUserFavorites(product.id, userId),
        })),
      );
      return this.formatProductsResponse(randomProductsWithFavorites);
    }

    return this.formatProductsResponse(
      randomProductsList.map((product) => ({
        ...product,
        isFavorited: false,
      })),
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
        videoUrl: true,
      },
    });

    const result = Promise.all(
      products.map(async (product) => {
        return {
          ...product,
          isFavorited: await this.isProductInUserFavorites(product.id, id),
        };
      }),
    );

    return this.formatProductsResponse(await result);
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
            videoUrl: true,
          },
        },
      },
    });

    return this.formatProductsResponse(
      checkUser?.favorites.map((product) => {
        return {
          ...product,
          isFavorited: true,
        };
      }) || [],
    );
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
        userId: true,
        videoUrl: true,
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
        type: {
          select: {
            id: true,
            name: true,
          },
        },
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
      fieldValues: product.fieldValues.map((field) => {
        return {
          id: field.id,
          [field.field.name]: field.value,
        };
      }),
      images: product?.images.map((img) => `${this.baseUrl}${img}`),
      isFavorited: userId
        ? await this.isProductInUserFavorites(product.id, userId)
        : false,
      seller: await this.userService.getUserInfo(product.userId),
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
