import {
  Body,
  Controller,
  Delete,
  Get,
  Param,
  Patch,
  Post,
  Put,
  Query,
  Req,
  UploadedFiles,
  UseGuards,
  UseInterceptors,
} from '@nestjs/common';
import { FilesInterceptor } from '@nestjs/platform-express';
import {
  ApiBearerAuth,
  ApiBody,
  ApiConsumes,
  ApiOperation,
  ApiQuery,
  ApiResponse,
  ApiTags,
} from '@nestjs/swagger';

import type { Request } from 'express';

import { AdminSessionAuthGuard } from '@/auth/guards/admin-session-auth.guard';
import { OptionalSessionAuthGuard } from '@/auth/guards/optional-session-auth.guard';
import { SessionAuthGuard } from '@/auth/guards/session-auth.guard';

import { createProductDto } from './dto/create-product.dto';
import { SearchProductsDto } from './dto/search-products.dto';
import { UpdateProductDto } from './dto/update-product.dto';
import { ModerateState } from './enum/moderate-state.enum';
import { ProductService } from './product.service';

@ApiTags('Products')
@Controller('product')
export class ProductController {
  constructor(private readonly productService: ProductService) {}

  @Post('create')
  @UseGuards(SessionAuthGuard)
  @UseInterceptors(FilesInterceptor('images', 8))
  @ApiOperation({
    summary: 'Создание нового продукта',
    description:
      'Создает новый продукт с возможностью загрузки до 8 изображений',
  })
  @ApiBearerAuth()
  @ApiConsumes('multipart/form-data')
  @ApiBody({
    description: 'Данные продукта и изображения',
    schema: {
      type: 'object',
      properties: {
        name: {
          type: 'string',
          description: 'Название продукта',
          example: 'iPhone 15 Pro',
        },
        price: {
          type: 'number',
          description: 'Цена в рублях',
          example: 120000,
        },
        state: {
          type: 'string',
          enum: ['NEW', 'USED'],
          description: 'Состояние товара',
          example: 'NEW',
        },
        description: {
          type: 'string',
          description: 'Описание товара',
          example: 'Новый iPhone 15 Pro в отличном состоянии',
        },
        address: {
          type: 'string',
          description: 'Адрес продавца',
          example: 'г. Москва, ул. Тверская, д. 1',
        },
        categoryId: {
          type: 'number',
          description: 'ID категории',
          example: 1,
        },
        subcategoryId: {
          type: 'number',
          description: 'ID подкатегории',
          example: 1,
        },
        typeId: {
          type: 'number',
          description: 'ID типа подкатегории (например: Измерительные приборы)',
          example: 1,
        },
        fieldValues: {
          type: 'string',
          description:
            'Дополнительные поля в формате JSON объекта (fieldId: значение). Например: {"1":"Тонометр","2":"Omron"}',
          example: '{"1":"Тонометр","2":"Omron"}',
        },
        images: {
          type: 'array',
          items: {
            type: 'string',
            format: 'binary',
          },
          description: 'Изображения продукта (до 8 файлов)',
        },
        videoUrl: {
          type: 'string',
          description: 'Ссылка на видео',
          example: 'https://www.youtube.com/watch?v=8u6xaEYGLa0',
        },
      },
      required: ['name', 'price', 'state', 'categoryId', 'subcategoryId'],
    },
  })
  @ApiResponse({
    status: 201,
    description: 'Продукт успешно создан',
    schema: {
      type: 'object',
      properties: {
        message: {
          type: 'string',
          example: 'Продукт успешно создан',
        },
        product: {
          type: 'object',
          properties: {
            id: { type: 'number', example: 1 },
            name: { type: 'string', example: 'iPhone 15 Pro' },
            price: { type: 'number', example: 120000 },
            state: { type: 'string', example: 'NEW' },
            images: {
              type: 'array',
              items: { type: 'string' },
              example: [
                '/uploads/product/image1.jpg',
                '/uploads/product/image2.jpg',
              ],
            },
            category: {
              type: 'object',
              properties: {
                id: { type: 'number', example: 1 },
                name: { type: 'string', example: 'Электроника' },
              },
            },
            subCategory: {
              type: 'object',
              properties: {
                id: { type: 'number', example: 1 },
                name: { type: 'string', example: 'Смартфоны' },
              },
            },
            user: {
              type: 'object',
              properties: {
                id: { type: 'number', example: 1 },
                fullName: { type: 'string', example: 'Иван Иванов' },
                email: { type: 'string', example: 'ivan@example.com' },
                rating: { type: 'number', example: 4.5 },
              },
            },
          },
        },
      },
    },
  })
  @ApiResponse({
    status: 400,
    description: 'Ошибка валидации данных',
  })
  @ApiResponse({
    status: 401,
    description: 'Не авторизован',
  })
  async createProduct(
    @Body() dto: createProductDto,
    @UploadedFiles() images: Express.Multer.File[],
    @Req() req: Request & { user: any },
  ) {
    return await this.productService.createProduct(dto, images, req.user.id);
  }

  @ApiOperation({
    summary: 'Удаление продукта',
  })
  @UseGuards(SessionAuthGuard)
  @Delete(':id')
  async deleteProduct(
    @Param('id') id: string,
    @Req() req: Request & { user: any },
  ) {
    return await this.productService.deleteProduct(+id, req.user.id);
  }

  @Patch(':id')
  @UseGuards(SessionAuthGuard)
  @UseInterceptors(FilesInterceptor('images', 8))
  @ApiOperation({
    summary: 'Обновление товара',
    description:
      'Обновляет данные товара. Можно добавить новые характеристики (fieldValues) к существующим или обновить их значения.',
  })
  @ApiBearerAuth()
  @ApiConsumes('multipart/form-data')
  @ApiBody({
    description: 'Данные для обновления товара',
    schema: {
      type: 'object',
      properties: {
        name: {
          type: 'string',
          description: 'Название продукта',
          example: 'iPhone 15 Pro Max',
        },
        price: {
          type: 'number',
          description: 'Цена в рублях',
          example: 130000,
        },
        state: {
          type: 'string',
          enum: ['NEW', 'USED'],
          description: 'Состояние товара',
          example: 'NEW',
        },
        description: {
          type: 'string',
          description: 'Описание товара',
          example: 'Обновленное описание товара',
        },
        address: {
          type: 'string',
          description: 'Адрес продавца',
          example: 'г. Москва, ул. Пушкина, д. 10',
        },
        videoUrl: {
          type: 'string',
          description: 'Ссылка на видео',
          example: 'https://www.youtube.com/watch?v=newvideo',
        },
        fieldValues: {
          type: 'string',
          description:
            'Дополнительные поля для добавления/обновления в формате JSON. Существующие поля будут обновлены, новые - добавлены. Например: {"3":"Новая характеристика"}',
          example: '{"1":"XXL","2":"Красный"}',
        },
        images: {
          type: 'array',
          items: {
            type: 'string',
            format: 'binary',
          },
          description:
            'Дополнительные изображения (будут добавлены к существующим)',
        },
      },
    },
  })
  @ApiResponse({
    status: 200,
    description: 'Товар успешно обновлён',
  })
  @ApiResponse({
    status: 400,
    description: 'Ошибка валидации данных',
  })
  @ApiResponse({
    status: 401,
    description: 'Не авторизован',
  })
  @ApiResponse({
    status: 403,
    description: 'Нет прав для редактирования этого товара',
  })
  async updateProduct(
    @Param('id') id: string,
    @Body() dto: UpdateProductDto,
    @UploadedFiles() images: Express.Multer.File[],
    @Req() req: Request & { user: any },
  ) {
    return await this.productService.updateProduct(
      +id,
      dto,
      images,
      req.user.id,
    );
  }

  @ApiOperation({
    summary: 'Получение доступных фильтров для товаров',
    description:
      'Возвращает все доступные значения фильтров (размер, цвет и т.д.) для заданной категории/подкатегории/типа',
  })
  @ApiQuery({
    name: 'categorySlug',
    required: false,
    description: 'Slug категории',
  })
  @ApiQuery({
    name: 'subCategorySlug',
    required: false,
    description: 'Slug подкатегории',
  })
  @ApiQuery({
    name: 'typeSlug',
    required: false,
    description: 'Slug типа подкатегории',
  })
  @ApiTags('Поиск')
  @Get('available-filters')
  async getAvailableFilters(
    @Query('categorySlug') categorySlug?: string,
    @Query('subCategorySlug') subCategorySlug?: string,
    @Query('typeSlug') typeSlug?: string,
  ) {
    return await this.productService.getAvailableFilters(
      categorySlug,
      subCategorySlug,
      typeSlug,
    );
  }

  @ApiOperation({
    summary: 'Получение всех товаров с возможностью фильтрации',
    description:
      'Возвращает все товары или фильтрует их по заданным параметрам. Если параметры не указаны, возвращаются все товары.',
  })
  @ApiTags('Поиск')
  @UseGuards(OptionalSessionAuthGuard)
  @Get('all-products')
  async findAll(
    @Req() req: Request & { user: any },
    @Query() searchDto: SearchProductsDto,
  ) {
    return await this.productService.findAll(req?.user?.id, searchDto);
  }

  @ApiOperation({
    summary: 'Возвращает 5 рандомных товаров',
  })
  @UseGuards(OptionalSessionAuthGuard)
  @Get('random-products')
  async getRandomProducts(@Req() req: Request & { user: any }) {
    return await this.productService.getRandomProducts(req?.user?.id);
  }

  @ApiOperation({
    summary: 'Получение товаров пользователя',
    description:
      'Возвращает список всех товаров, созданных текущим пользователем',
  })
  @Get('user-products/:id')
  async getMyProducts(@Param('id') id: string) {
    return await this.productService.getProductsByUserId(+id);
  }

  @ApiTags('Избранное')
  @ApiOperation({
    summary: 'Добавление товара в избранное',
  })
  @Post('add-to-favorites/:id')
  @UseGuards(SessionAuthGuard)
  async addToFavorites(
    @Param('id') id: string,
    @Req() req: Request & { user: any },
  ) {
    return await this.productService.addProductToFavorites(+id, req.user.id);
  }

  @ApiTags('Избранное')
  @ApiOperation({
    summary: 'Удаление товара из избранного',
  })
  @Delete('remove-from-favorites/:id')
  @UseGuards(SessionAuthGuard)
  async removeFromFavorites(
    @Req() req: Request & { user: any },
    @Param('id') id: string,
  ) {
    return await this.productService.removeProductFromFavorites(
      +id,
      req.user.id,
    );
  }

  @ApiTags('Избранное')
  @ApiOperation({
    summary: 'Получение всех товаров из избранного',
  })
  @Get('my-favorites')
  @UseGuards(SessionAuthGuard)
  async getMyFavorites(@Req() req: Request & { user: any }) {
    return await this.productService.getMyFavorites(req.user.id);
  }

  @ApiOperation({
    summary: 'Получение данных товара для карточки по id',
  })
  @UseGuards(OptionalSessionAuthGuard)
  @Get('product-card/:id')
  async getProductCard(
    @Param('id') id: string,
    @Req() req?: Request & { user: any },
  ) {
    return await this.productService.getProductCard(+id, req?.user?.id);
  }

  @ApiOperation({ summary: 'Смена статуса видимости товара' })
  @UseGuards(SessionAuthGuard)
  @Put('toggle-product/:id')
  async toggleProduct(
    @Param('id') id: string,
    @Req() req: Request & { user: any },
  ) {
    return await this.productService.toggleProduct(+id, req.user.id);
  }

  @ApiOperation({ summary: 'Модерация товара (одобрить/отклонить)' })
  @UseGuards(AdminSessionAuthGuard)
  @ApiQuery({
    name: 'status',
    enum: ModerateState,
    required: true,
    description: 'Статус модерации: APPROVED или DENIDED',
  })
  @ApiQuery({
    name: 'reason',
    required: false,
    description: 'Причина отказа (обязательна при status=DENIDED)',
  })
  @Put('moderate-product/:id')
  async moderateProduct(
    @Param('id') id: string,
    @Query('status') status: ModerateState,
    @Query('reason') reason?: string,
  ) {
    return await this.productService.moderateProduct(+id, status, reason);
  }

  @ApiOperation({
    summary: 'Возвращает все товары на модерации',
  })
  @UseGuards(AdminSessionAuthGuard)
  @Get('all-products-to-moderate')
  async allProductsToModerate() {
    return await this.productService.allProductsToModerate();
  }
}
