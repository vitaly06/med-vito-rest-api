import {
  Injectable,
  BadRequestException,
  ForbiddenException,
  NotFoundException,
} from '@nestjs/common';
import { PrismaService } from 'src/prisma/prisma.service';
import { createProductDto } from './dto/create-product.dto';
import { UpdateProductDto } from './dto/update-product.dto';
import { ConfigService } from '@nestjs/config';
import { UserService } from 'src/user/user.service';
import { ModerateState } from './enum/moderate-state.enum';
import { S3Service } from 'src/s3/s3.service';
import { ChatService } from 'src/chat/chat.service';
import { ChatGateway } from 'src/chat/gateway';
import { createOkseiProductDto } from './dto/oksei-create-product.dto';

@Injectable()
export class ProductService {
  baseUrl: string;
  constructor(
    private readonly prisma: PrismaService,
    private readonly configService: ConfigService,
    private readonly userService: UserService,
    private readonly s3Service: S3Service,
    private readonly chatService: ChatService,
    private readonly chatGateway: ChatGateway,
  ) {
    this.baseUrl = this.configService.get<string>(
      'BASE_URL',
      'http://localhost:3000',
    );
  }

  async createProduct(
    dto: createProductDto,
    files: Express.Multer.File[],
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

      // Загружаем изображения в S3
      const imageUrls =
        files && files.length > 0
          ? await this.s3Service.uploadFiles(files, 'products')
          : [];

      const product = await this.prisma.product.create({
        data: {
          name: dto.name,
          price: +dto.price,
          state: dto.state,
          description: dto.description,
          address: dto.address,
          images: imageUrls,
          categoryId: +dto.categoryId,
          moderateState: 'MODERATE',
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
          images: product.images, // URL уже полные из S3
        },
      };
    } catch (error) {
      throw new BadRequestException(
        `Ошибка при создании продукта: ${error.message}`,
      );
    }
  }

  async createOkseiProduct(
    dto: createOkseiProductDto,
    file: Express.Multer.File,
  ) {
    // Проверяем наличие файла
    if (!file) {
      throw new BadRequestException('Изображение обязательно для загрузки');
    }

    // Загружаем изображение в S3
    const imageUrl = await this.s3Service.uploadFile(file, 'products');

    const product = await this.prisma.okseiProduct.create({
      data: {
        name: dto.name,
        price: +dto.price,
        description: dto.description || '',
        image: imageUrl,
      },
    });

    return {
      message: 'Продукт успешно создан',
      product: {
        ...product,
        image: product.image, // URL уже полные из S3
      },
    };
  }

  async deleteProduct(productId, userId: number) {
    const checkProduct = await this.prisma.product.findUnique({
      where: { id: productId },
      select: {
        userId: true,
        images: true,
      },
    });

    if (!checkProduct) {
      throw new BadRequestException('Товар для удаления не найден');
    }

    if (checkProduct.userId != userId) {
      throw new ForbiddenException('Вы не можете удалить чужой товар');
    }

    // Удаляем изображения из S3
    if (checkProduct.images && checkProduct.images.length > 0) {
      await this.s3Service.deleteFiles(checkProduct.images);
    }

    await this.prisma.product.delete({
      where: { id: productId },
    });

    return { message: 'Товар успешно удалён' };
  }

  async findAllForOksei() {
    const products = await this.prisma.okseiProduct.findMany({
      select: {
        id: true,
        name: true,
        description: true,
        image: true,
      },
    });

    return products;
  }

  async updateProduct(
    productId: number,
    dto: UpdateProductDto,
    files: Express.Multer.File[],
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

      // Обработка новых изображений - загружаем в S3
      if (files && files.length > 0) {
        const newImageUrls = await this.s3Service.uploadFiles(
          files,
          'products',
        );
        // Добавляем новые изображения к существующим
        updateData.images = [...existingProduct.images, ...newImageUrls];
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
          images: updatedProduct!.images, // URL уже полные из S3
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
        categorySlug,
        subCategoryId,
        subCategorySlug,
        typeId,
        typeSlug,
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

      // Фильтр по категории (приоритет у slug)
      if (categorySlug) {
        const category = await this.prisma.category.findUnique({
          where: { slug: categorySlug },
        });
        if (category) {
          whereConditions.categoryId = category.id;
        }
      } else if (categoryId) {
        whereConditions.categoryId = parseInt(categoryId);
      }

      // Фильтр по подкатегории (приоритет у slug)
      if (subCategorySlug) {
        const subCategory = await this.prisma.subCategory.findFirst({
          where: { slug: subCategorySlug },
        });
        if (subCategory) {
          whereConditions.subCategoryId = subCategory.id;
        }
      } else if (subCategoryId) {
        whereConditions.subCategoryId = parseInt(subCategoryId);
      }

      // Фильтр по типу подкатегории (приоритет у slug)
      if (typeSlug) {
        const type = await this.prisma.subcategotyType.findFirst({
          where: { slug: typeSlug },
        });
        if (type) {
          whereConditions.typeId = type.id;
        }
      } else if (typeId) {
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
        where: {
          AND: [
            whereConditions,
            { moderateState: 'APPROVED' },
            { isHide: false },
          ],
        },
        select: {
          id: true,
          images: true,
          name: true,
          address: true,
          createdAt: true,
          price: true,
          userId: true,
          videoUrl: true,
          category: {
            select: {
              id: true,
              name: true,
              slug: true,
            },
          },
          subCategory: {
            select: {
              id: true,
              name: true,
              slug: true,
            },
          },
          type: {
            select: {
              id: true,
              name: true,
              slug: true,
            },
          },
          promotions: {
            where: {
              isActive: true,
              isPaid: true,
              endDate: {
                gte: new Date(), // Только активные продвижения
              },
            },
            include: {
              promotion: {
                select: {
                  pricePerDay: true,
                },
              },
            },
            orderBy: {
              promotion: {
                pricePerDay: 'desc', // Самый дорогой тариф
              },
            },
            take: 1, // Берем только самое дорогое продвижение
          },
        },
        orderBy,
        skip: (page - 1) * limit,
        take: limit,
      });

      // Сортируем товары по продвижению
      const sortedProducts = products.sort((a, b) => {
        const aPromo = a.promotions[0]?.promotion?.pricePerDay || 0;
        const bPromo = b.promotions[0]?.promotion?.pricePerDay || 0;
        return bPromo - aPromo; // Сначала более дорогие продвижения
      });

      // Форматируем результат как в обычном findAll
      console.log(sortedProducts.length);
      if (userId) {
        const result = await Promise.all(
          sortedProducts.map(async (product) => ({
            ...product,
            isFavorited: await this.isProductInUserFavorites(
              product.id,
              userId,
            ),
            hasPromotion: product.promotions.length > 0,
            promotionLevel: product.promotions[0]?.promotion?.pricePerDay || 0,
          })),
        );
        return this.formatProductsResponse(result);
      }

      return this.formatProductsResponse(
        sortedProducts.map((product) => ({
          ...product,
          isFavorited: false,
          hasPromotion: product.promotions.length > 0,
          promotionLevel: product.promotions[0]?.promotion?.pricePerDay || 0,
        })),
      );
    }

    // Если параметров поиска нет - возвращаем все товары

    // Получаем все товары
    const allProducts = await this.prisma.product.findMany({
      select: {
        id: true,
        images: true,
        name: true,
        address: true,
        createdAt: true,
        isHide: true,
        price: true,
        userId: true,
        videoUrl: true,
        category: {
          select: {
            id: true,
            name: true,
            slug: true,
          },
        },
        subCategory: {
          select: {
            id: true,
            name: true,
            slug: true,
          },
        },
        type: {
          select: {
            id: true,
            name: true,
            slug: true,
          },
        },
        promotions: {
          where: {
            isActive: true,
            isPaid: true,
            endDate: {
              gte: new Date(), // Только активные продвижения
            },
          },
          include: {
            promotion: {
              select: {
                pricePerDay: true,
              },
            },
          },
          orderBy: {
            promotion: {
              pricePerDay: 'desc', // Самый дорогой тариф
            },
          },
          take: 1, // Берем только самое дорогое продвижение
        },
      },
      where: {
        isHide: false,
        moderateState: 'APPROVED',
      },
      orderBy: { createdAt: 'desc' },
    });

    // Сортируем товары: сначала с продвижением (дорогие → дешевые), потом без продвижения
    const sortedProducts = allProducts.sort((a, b) => {
      const aPromo = a.promotions[0]?.promotion?.pricePerDay || 0;
      const bPromo = b.promotions[0]?.promotion?.pricePerDay || 0;

      // Если у обоих есть продвижение или у обоих нет - сортируем по цене продвижения
      if ((aPromo > 0 && bPromo > 0) || (aPromo === 0 && bPromo === 0)) {
        return bPromo - aPromo;
      }

      // Товары с продвижением всегда выше товаров без продвижения
      return bPromo > 0 ? 1 : -1;
    });

    // Форматируем оба списка
    if (userId) {
      const allProductsWithFavorites = await Promise.all(
        sortedProducts.map(async (product) => ({
          ...product,
          isFavorited: await this.isProductInUserFavorites(product.id, userId),
          hasPromotion: product.promotions.length > 0,
          promotionLevel: product.promotions[0]?.promotion?.pricePerDay || 0,
        })),
      );

      return this.formatProductsResponse(allProductsWithFavorites);
    }

    return this.formatProductsResponse(
      sortedProducts.map((product) => ({
        ...product,
        isFavorited: false,
        hasPromotion: product.promotions.length > 0,
        promotionLevel: product.promotions[0]?.promotion?.pricePerDay || 0,
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
        SELECT 1 FROM "Product" p WHERE p."subCategoryId" = sc.id AND p."isHide" = false AND p."moderateState" = 'APPROVED'
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
        isHide: true,
        moderateState: true,
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
        images: product.images || [], // URL уже полные из S3
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
        isHide: true,
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
      images: product?.images || [], // URL уже полные из S3
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

  async toggleProduct(productId: number, userId: number) {
    const checkProduct = await this.prisma.product.findUnique({
      where: { id: productId },
    });

    if (!checkProduct) {
      throw new NotFoundException('Товар не надйен');
    }

    if (checkProduct.userId != userId) {
      throw new ForbiddenException('Вы не можете редактировать не свой товар');
    }

    await this.prisma.product.update({
      where: { id: productId },
      data: {
        isHide: !checkProduct.isHide,
      },
    });

    return { message: 'Статус активности товара сменен' };
  }

  async moderateProduct(
    productId: number,
    status: ModerateState,
    rejectionReason?: string,
  ) {
    if (!['APPROVED', 'DENIDED'].includes(status)) {
      throw new BadRequestException(
        'Неверный статус модерации. Доступные статуты: APPROVED, DENIDED',
      );
    }

    const checkProduct = await this.prisma.product.findUnique({
      where: { id: productId },
      include: {
        user: {
          select: {
            id: true,
            fullName: true,
          },
        },
      },
    });

    if (!checkProduct) {
      throw new NotFoundException('Товар для модерации не найден');
    }

    // Если отказ, проверяем наличие причины
    if (status === 'DENIDED' && !rejectionReason) {
      throw new BadRequestException(
        'Необходимо указать причину отказа в модерации',
      );
    }

    // Обновляем товар
    await this.prisma.product.update({
      where: { id: productId },
      data: {
        moderateState: status,
        moderationRejectionReason:
          status === 'DENIDED' ? rejectionReason : null,
      },
    });

    // Если отказ, отправляем уведомление владельцу товара
    if (status === 'DENIDED' && rejectionReason) {
      try {
        // Находим любого администратора для отправки системного сообщения
        const adminRole = await this.prisma.role.findFirst({
          where: { name: 'admin' },
          include: {
            users: {
              take: 1,
            },
          },
        });

        if (!adminRole || adminRole.users.length === 0) {
          throw new NotFoundException(
            'No admin user found for sending moderation notification',
          );
          return; // Не отправляем уведомление, если нет админа
        }

        const adminUserId = adminRole.users[0].id;

        // Проверяем существует ли чат
        let chat = await this.prisma.chat.findFirst({
          where: {
            productId: productId,
            OR: [
              { buyerId: adminUserId, sellerId: checkProduct.userId },
              { buyerId: checkProduct.userId, sellerId: adminUserId },
            ],
          },
        });

        // Если чата нет, создаем новый
        if (!chat) {
          chat = await this.prisma.chat.create({
            data: {
              productId: productId,
              buyerId: adminUserId,
              sellerId: checkProduct.userId,
            },
          });
        }

        // Отправляем сообщение с причиной отказа
        const messageContent = `❌ Ваш товар "${checkProduct.name}" был отклонен модерацией.\n\nПричина отказа: ${rejectionReason}`;

        const message = await this.chatService.sendMessage(
          chat.id,
          adminUserId,
          messageContent,
        );

        // Отправляем WebSocket уведомление владельцу товара
        this.chatGateway.server
          .to(`user_${checkProduct.userId}`)
          .emit('newMessage', {
            id: message.id,
            content: message.content,
            senderId: message.senderId,
            sender: message.sender,
            createdAt: message.createdAt,
            timeString: message.timeString,
            chatId: chat.id,
            isModeration: true,
          });

        console.log(
          `Moderation rejection notification sent to user ${checkProduct.userId} for product ${productId}`,
        );
      } catch (error) {
        console.error(
          'Error sending moderation rejection notification:',
          error,
        );
        // Не выбрасываем ошибку, чтобы модерация все равно прошла
      }
    }
  }

  async allProductsToModerate() {
    const allProducts = await this.prisma.product.findMany({
      select: {
        id: true,
        images: true,
        name: true,
        address: true,
        createdAt: true,
        isHide: true,
        price: true,
        userId: true,
        videoUrl: true,
      },
      where: {
        moderateState: 'MODERATE',
      },
      orderBy: { createdAt: 'desc' },
    });
    return this.formatProductsResponse(allProducts);
  }
}
