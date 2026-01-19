import {
  Body,
  Controller,
  Delete,
  Get,
  Ip,
  Param,
  ParseIntPipe,
  Post,
  Put,
  Query,
  Req,
  UploadedFile,
  UseGuards,
  UseInterceptors,
} from '@nestjs/common';
import { FileInterceptor } from '@nestjs/platform-express';
import {
  ApiBearerAuth,
  ApiBody,
  ApiConsumes,
  ApiOperation,
  ApiParam,
  ApiQuery,
  ApiResponse,
  ApiTags,
} from '@nestjs/swagger';
import { Request } from 'express';

import { AdminSessionAuthGuard } from 'src/auth/guards/admin-session-auth.guard';
import { SessionAuthGuard } from 'src/auth/guards/session-auth.guard';
import { OptionalSessionAuthGuard } from 'src/auth/guards/optional-session-auth.guard';

import { BannerService } from './banner.service';
import { CreateBannerDto } from './dto/create-banner.dto';
import { UpdateBannerDto } from './dto/update-banner.dto';
import { BannerPlace } from './entities/banner-place.enum';

@ApiTags('Banners')
@Controller('banner')
export class BannerController {
  constructor(private readonly bannerService: BannerService) {}

  @Post()
  @UseGuards(SessionAuthGuard)
  @UseInterceptors(FileInterceptor('image'))
  @ApiBearerAuth()
  @ApiOperation({
    summary: 'Создание баннера',
    description:
      'Создаёт новый баннер с загрузкой изображения в S3 и указанием места размещения',
  })
  @ApiConsumes('multipart/form-data')
  @ApiBody({
    description: 'Изображение и место размещения баннера',
    type: CreateBannerDto,
  })
  @ApiResponse({
    status: 201,
    description: 'Баннер успешно создан',
    schema: {
      example: {
        id: 1,
        photoUrl:
          'https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/banners/abc123.jpg',
        place: 'PRODUCT_FEED',
        name: 'Google Browser',
        navigateToUrl: 'https://google.com',
        userId: 1,
        createdAt: '2025-12-27T10:00:00.000Z',
        updatedAt: '2025-12-27T10:00:00.000Z',
      },
    },
  })
  @ApiResponse({ status: 400, description: 'Некорректные данные' })
  @ApiResponse({ status: 401, description: 'Не авторизован' })
  @ApiResponse({ status: 403, description: 'Доступ запрещён' })
  async create(
    @Req() request: Request & { user: any },
    @UploadedFile() file: Express.Multer.File,
    @Body('place') place: BannerPlace,
    @Body('navigateToUrl') navigateToUrl: string,
    @Body('name') name: string,
  ) {
    const userId = request.user.id;
    return this.bannerService.create(file, place, navigateToUrl, name, userId);
  }

  @Get('random')
  @ApiOperation({
    summary: 'Получить случайные баннеры',
    description:
      'Возвращает 5 случайных баннеров из БД (или меньше, если их меньше 5)',
  })
  @ApiResponse({
    status: 200,
    description: 'Список случайных баннеров получен',
    schema: {
      example: [
        {
          id: 1,
          photoUrl:
            'https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/banners/abc123.jpg',
          place: 'PRODUCT_FEED',
          name: 'Yandex Browser',
          navigateToUrl: 'https://yandex.ru',
          userId: 7106521,
          createdAt: '2025-12-27T10:00:00.000Z',
          updatedAt: '2025-12-27T10:00:00.000Z',
        },
      ],
    },
  })
  async findRandom() {
    return this.bannerService.findRandom(5);
  }

  @Get()
  @ApiOperation({
    summary: 'Получить баннеры',
    description:
      'Возвращает список всех баннеров или баннеры по месту размещения',
  })
  @ApiQuery({
    name: 'place',
    enum: BannerPlace,
    required: false,
    description: 'Фильтр по месту размещения баннера',
    example: BannerPlace.PRODUCT_FEED,
  })
  @ApiResponse({
    status: 200,
    description: 'Список баннеров получен',
    schema: {
      example: [
        {
          id: 1,
          photoUrl:
            'https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/banners/abc123.jpg',
          place: 'PRODUCT_FEED',
          createdAt: '2025-12-27T10:00:00.000Z',
          updatedAt: '2025-12-27T10:00:00.000Z',
        },
        {
          id: 2,
          photoUrl:
            'https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/banners/def456.jpg',
          place: 'PROFILE',
          createdAt: '2025-12-27T11:00:00.000Z',
          updatedAt: '2025-12-27T11:00:00.000Z',
        },
      ],
    },
  })
  async findAll(@Query('place') place?: BannerPlace) {
    if (place) {
      return this.bannerService.findByPlace(place);
    }
    return this.bannerService.findAll();
  }

  @Get(':id')
  @ApiOperation({
    summary: 'Получить баннер по ID',
    description: 'Возвращает баннер с указанным ID',
  })
  @ApiParam({
    name: 'id',
    type: 'number',
    description: 'ID баннера',
    example: 1,
  })
  @ApiResponse({
    status: 200,
    description: 'Баннер найден',
    schema: {
      example: {
        id: 1,
        photoUrl:
          'https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/banners/abc123.jpg',
        place: 'PRODUCT_FEED',
        createdAt: '2025-12-27T10:00:00.000Z',
        updatedAt: '2025-12-27T10:00:00.000Z',
      },
    },
  })
  @ApiResponse({ status: 404, description: 'Баннер не найден' })
  async findOne(@Param('id', ParseIntPipe) id: number) {
    return this.bannerService.findOne(id);
  }

  @Put(':id')
  @UseGuards(AdminSessionAuthGuard)
  @UseInterceptors(FileInterceptor('image'))
  @ApiBearerAuth()
  @ApiOperation({
    summary: 'Обновить баннер',
    description: 'Обновляет изображение и/или место размещения баннера',
  })
  @ApiConsumes('multipart/form-data')
  @ApiParam({
    name: 'id',
    type: 'number',
    description: 'ID баннера',
    example: 1,
  })
  @ApiBody({
    description: 'Новое изображение и/или место размещения баннера',
    type: UpdateBannerDto,
  })
  @ApiResponse({
    status: 200,
    description: 'Баннер успешно обновлён',
    schema: {
      example: {
        id: 1,
        photoUrl:
          'https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/banners/new123.jpg',
        place: 'FAVORITES',
        createdAt: '2025-12-27T10:00:00.000Z',
        updatedAt: '2025-12-27T12:00:00.000Z',
      },
    },
  })
  @ApiResponse({ status: 400, description: 'Некорректные данные' })
  @ApiResponse({ status: 401, description: 'Не авторизован' })
  @ApiResponse({ status: 403, description: 'Доступ запрещён' })
  @ApiResponse({ status: 404, description: 'Баннер не найден' })
  async update(
    @Param('id', ParseIntPipe) id: number,
    @UploadedFile() file?: Express.Multer.File,
    @Body('place') place?: BannerPlace,
    @Body('navigateToUrl') navigateToUrl?: string,
    @Body('name') name?: string,
  ) {
    return this.bannerService.update(id, file, place, navigateToUrl, name);
  }

  @Delete(':id')
  @UseGuards(AdminSessionAuthGuard)
  @ApiBearerAuth()
  @ApiOperation({
    summary: 'Удалить баннер',
    description: 'Удаляет баннер и его изображение из S3',
  })
  @ApiParam({
    name: 'id',
    type: 'number',
    description: 'ID баннера',
    example: 1,
  })
  @ApiResponse({
    status: 200,
    description: 'Баннер успешно удалён',
    schema: {
      example: {
        message: 'Баннер успешно удалён',
      },
    },
  })
  @ApiResponse({ status: 401, description: 'Не авторизован' })
  @ApiResponse({ status: 403, description: 'Доступ запрещён' })
  @ApiResponse({ status: 404, description: 'Баннер не найден' })
  async remove(@Param('id', ParseIntPipe) id: number) {
    return this.bannerService.remove(id);
  }

  @Post(':id/view')
  @UseGuards(OptionalSessionAuthGuard)
  @ApiOperation({
    summary: 'Зарегистрировать просмотр баннера',
    description:
      'Регистрирует просмотр баннера для аналитики. Можно вызывать без авторизации',
  })
  @ApiParam({
    name: 'id',
    type: 'number',
    description: 'ID баннера',
    example: 1,
  })
  @ApiResponse({
    status: 201,
    description: 'Просмотр зарегистрирован',
    schema: {
      example: {
        id: 1,
        bannerId: 1,
        userId: 1,
        ipAddress: '192.168.1.1',
        viewedAt: '2025-12-27T10:00:00.000Z',
      },
    },
  })
  @ApiResponse({ status: 404, description: 'Баннер не найден' })
  async registerView(
    @Param('id', ParseIntPipe) bannerId: number,
    @Req() request: Request & { user: any },
    @Ip() ipAddress: string,
  ) {
    const userId = request.user?.id;
    return this.bannerService.registerView(bannerId, userId, ipAddress);
  }

  @Get(':id/stats')
  @UseGuards(SessionAuthGuard)
  @ApiBearerAuth()
  @ApiOperation({
    summary: 'Получить статистику по баннеру',
    description:
      'Возвращает статистику просмотров для конкретного баннера (только для владельца)',
  })
  @ApiParam({
    name: 'id',
    type: 'number',
    description: 'ID баннера',
    example: 1,
  })
  @ApiResponse({
    status: 200,
    description: 'Статистика по баннеру',
    schema: {
      example: {
        bannerId: 1,
        totalViews: 150,
        uniqueUsers: 85,
        viewsByDay: [
          { date: '2025-01-15', views: 25 },
          { date: '2025-01-16', views: 30 },
        ],
      },
    },
  })
  @ApiResponse({ status: 401, description: 'Не авторизован' })
  @ApiResponse({ status: 403, description: 'Доступ запрещён' })
  @ApiResponse({ status: 404, description: 'Баннер не найден' })
  async getBannerStats(@Param('id', ParseIntPipe) bannerId: number) {
    return this.bannerService.getBannerStats(bannerId);
  }

  @Get('my-stats/all')
  @UseGuards(SessionAuthGuard)
  @ApiBearerAuth()
  @ApiOperation({
    summary: 'Получить статистику по всем баннерам пользователя',
    description:
      'Возвращает статистику просмотров для всех баннеров текущего пользователя',
  })
  @ApiResponse({
    status: 200,
    description: 'Статистика по баннерам пользователя',
    schema: {
      example: [
        {
          id: 1,
          name: 'Google Browser',
          photoUrl:
            'https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/banners/abc123.jpg',
          place: 'PRODUCT_FEED',
          navigateToUrl: 'https://google.com',
          totalViews: 150,
        },
        {
          id: 2,
          name: 'Yandex Search',
          photoUrl:
            'https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/banners/def456.jpg',
          place: 'PROFILE',
          navigateToUrl: 'https://yandex.ru',
          totalViews: 85,
        },
      ],
    },
  })
  @ApiResponse({ status: 401, description: 'Не авторизован' })
  async getUserBannersStats(@Req() request: Request & { user: any }) {
    const userId = request.user.id;
    return this.bannerService.getUserBannersStats(userId);
  }
}
