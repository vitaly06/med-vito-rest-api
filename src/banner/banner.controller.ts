import {
  Body,
  Controller,
  Delete,
  Get,
  Param,
  ParseIntPipe,
  Post,
  Put,
  Query,
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

import { AdminSessionAuthGuard } from 'src/auth/guards/admin-session-auth.guard';

import { BannerService } from './banner.service';
import { CreateBannerDto } from './dto/create-banner.dto';
import { UpdateBannerDto } from './dto/update-banner.dto';
import { BannerPlace } from './entities/banner-place.enum';

@ApiTags('Banners')
@Controller('banner')
export class BannerController {
  constructor(private readonly bannerService: BannerService) {}

  @Post()
  @UseGuards(AdminSessionAuthGuard)
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
        createdAt: '2025-12-27T10:00:00.000Z',
        updatedAt: '2025-12-27T10:00:00.000Z',
      },
    },
  })
  @ApiResponse({ status: 400, description: 'Некорректные данные' })
  @ApiResponse({ status: 401, description: 'Не авторизован' })
  @ApiResponse({ status: 403, description: 'Доступ запрещён' })
  async create(
    @UploadedFile() file: Express.Multer.File,
    @Body('place') place: BannerPlace,
  ) {
    return this.bannerService.create(file, place);
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
  ) {
    return this.bannerService.update(id, file, place);
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
}
