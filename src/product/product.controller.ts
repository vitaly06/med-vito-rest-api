import {
  Body,
  Controller,
  Post,
  UploadedFiles,
  UseGuards,
  UseInterceptors,
  Req,
  Get,
  Param,
  Delete,
  Query,
} from '@nestjs/common';
import { ProductService } from './product.service';
import { JwtAuthGuard } from 'src/auth/guards/jwt-auth.guard';
import { OptionalJwtAuthGuard } from 'src/auth/guards/optional-jwt-auth.guard';
import { createProductDto } from './dto/create-product.dto';
import { SearchProductsDto } from './dto/search-products.dto';
import { FilesInterceptor } from '@nestjs/platform-express';
import type { Request } from 'express';
import {
  ApiTags,
  ApiOperation,
  ApiConsumes,
  ApiBody,
  ApiBearerAuth,
  ApiResponse,
} from '@nestjs/swagger';
import { request } from 'https';

@ApiTags('Products')
@Controller('product')
export class ProductController {
  constructor(private readonly productService: ProductService) {}

  @Post('create')
  @UseGuards(JwtAuthGuard)
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
    const fileNames = images ? images.map((file) => file.filename) : [];

    return await this.productService.createProduct(dto, fileNames, req.user.id);
  }

  @ApiOperation({
    summary: 'Удаление продукта',
  })
  @UseGuards(JwtAuthGuard)
  @Delete(':id')
  async deleteProduct(
    @Param('id') id: string,
    @Req() req: Request & { user: any },
  ) {
    return await this.productService.deleteProduct(+id, req.user.id);
  }

  @ApiOperation({
    summary: 'Возвращает все товары для главной страницы',
  })
  @UseGuards(OptionalJwtAuthGuard)
  @Get('all-products')
  async findAll(@Req() req: Request & { user: any }) {
    return await this.productService.findAll(req?.user?.id);
  }

  @ApiOperation({
    summary: 'Получение товаров текущего пользователя',
    description:
      'Возвращает список всех товаров, созданных текущим пользователем',
  })
  @Get('my-products')
  @UseGuards(JwtAuthGuard)
  async getMyProducts(@Req() req: Request & { user: any }) {
    return await this.productService.getProductsByUserId(req.user.id);
  }

  @ApiTags('Избранное')
  @ApiOperation({
    summary: 'Добавление товара в избранное',
  })
  @Post('add-to-favorites/:id')
  @UseGuards(JwtAuthGuard)
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
  @UseGuards(JwtAuthGuard)
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
  @UseGuards(JwtAuthGuard)
  async getMyFavorites(@Req() req: Request & { user: any }) {
    return await this.productService.getMyFavorites(req.user.id);
  }

  @ApiOperation({
    summary: 'Получение данных товара для карточки по id',
  })
  @UseGuards(OptionalJwtAuthGuard)
  @Get('product-card/:id')
  async getProductCard(
    @Param('id') id: string,
    @Req() req?: Request & { user: any },
  ) {
    console.log(req);
    return await this.productService.getProductCard(+id, req?.user?.id);
  }

  // @ApiOperation({
  //   summary: 'Получение статистики просмотров товаров пользователя',
  // })
  // @UseGuards(JwtAuthGuard)
  // @Get('view-stats')
  // async getProductViewStats(
  //   @Req() req: Request & { user: any },
  //   @Query('page') page: string = '1',
  //   @Query('limit') limit: string = '20',
  // ) {
  //   return await this.productService.getProductViewStats(
  //     req.user.id,
  //     +page,
  //     +limit,
  //   );
  // }

  // @ApiOperation({
  //   summary:
  //     'Получение статистики добавлений в избранное для товаров пользователя',
  // })
  // @UseGuards(JwtAuthGuard)
  // @Get('favorite-stats')
  // async getFavoriteStats(
  //   @Req() req: Request & { user: any },
  //   @Query('page') page: string = '1',
  //   @Query('limit') limit: string = '20',
  // ) {
  //   return await this.productService.getFavoriteStats(
  //     req.user.id,
  //     +page,
  //     +limit,
  //   );
  // }

  @Get('search')
  @UseGuards(OptionalJwtAuthGuard)
  @ApiOperation({
    summary: 'Поиск товаров',
    description: 'Поиск товаров с фильтрацией по различным параметрам',
  })
  @ApiResponse({
    status: 200,
    description: 'Результаты поиска товаров',
    schema: {
      type: 'object',
      properties: {
        products: {
          type: 'array',
          items: {
            type: 'object',
            properties: {
              id: { type: 'string' },
              name: { type: 'string' },
              price: { type: 'number' },
              state: { type: 'string' },
              brand: { type: 'string' },
              model: { type: 'string' },
              description: { type: 'string' },
              address: { type: 'string' },
              images: {
                type: 'array',
                items: { type: 'string' },
              },
              category: {
                type: 'object',
                properties: {
                  id: { type: 'string' },
                  name: { type: 'string' },
                },
              },
              subCategory: {
                type: 'object',
                properties: {
                  id: { type: 'string' },
                  name: { type: 'string' },
                },
              },
              seller: {
                type: 'object',
                properties: {
                  id: { type: 'string' },
                  fullName: { type: 'string' },
                  rating: { type: 'number' },
                },
              },
              stats: {
                type: 'object',
                properties: {
                  views: { type: 'number' },
                  favorites: { type: 'number' },
                },
              },
              createdAt: { type: 'string', format: 'date-time' },
              isInFavorites: { type: 'boolean' },
            },
          },
        },
        pagination: {
          type: 'object',
          properties: {
            currentPage: { type: 'number' },
            totalPages: { type: 'number' },
            totalItems: { type: 'number' },
            hasNextPage: { type: 'boolean' },
            hasPrevPage: { type: 'boolean' },
          },
        },
      },
    },
  })
  async searchProducts(
    @Query() searchDto: SearchProductsDto,
    @Req() req: Request & { user: any },
  ) {
    return await this.productService.searchProducts(searchDto, req?.user?.id);
  }
}
