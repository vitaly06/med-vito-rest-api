import { Injectable, NotFoundException } from '@nestjs/common';
import { PrismaService } from 'src/prisma/prisma.service';
import { S3Service } from 'src/s3/s3.service';
import { BannerPlace } from './entities/banner-place.enum';

@Injectable()
export class BannerService {
  constructor(
    private readonly prisma: PrismaService,
    private readonly s3Service: S3Service,
  ) {}

  async create(file: Express.Multer.File, place: BannerPlace) {
    // Загружаем изображение в S3
    const photoUrl = await this.s3Service.uploadFile(file, 'banners');

    // Создаём баннер в БД
    const banner = await this.prisma.banner.create({
      data: {
        photoUrl,
        place,
      },
    });

    return banner;
  }

  async findAll() {
    return this.prisma.banner.findMany({
      orderBy: {
        createdAt: 'desc',
      },
    });
  }

  async findByPlace(place: BannerPlace) {
    return this.prisma.banner.findMany({
      where: {
        place,
      },
      orderBy: {
        createdAt: 'desc',
      },
    });
  }

  async findOne(id: number) {
    const banner = await this.prisma.banner.findUnique({
      where: { id },
    });

    if (!banner) {
      throw new NotFoundException(`Баннер с ID ${id} не найден`);
    }

    return banner;
  }

  async update(id: number, file?: Express.Multer.File, place?: BannerPlace) {
    // Проверяем существование баннера
    const existingBanner = await this.findOne(id);

    const updateData: any = {};

    // Если есть новое изображение
    if (file) {
      // Удаляем старое изображение из S3
      if (existingBanner.photoUrl) {
        await this.s3Service.deleteFile(existingBanner.photoUrl);
      }

      // Загружаем новое изображение
      const photoUrl = await this.s3Service.uploadFile(file, 'banners');
      updateData.photoUrl = photoUrl;
    }

    // Если передано место размещения
    if (place) {
      updateData.place = place;
    }

    // Обновляем баннер
    const banner = await this.prisma.banner.update({
      where: { id },
      data: updateData,
    });

    return banner;
  }

  async remove(id: number) {
    // Проверяем существование баннера
    const banner = await this.findOne(id);

    // Удаляем изображение из S3
    if (banner.photoUrl) {
      await this.s3Service.deleteFile(banner.photoUrl);
    }

    // Удаляем баннер из БД
    await this.prisma.banner.delete({
      where: { id },
    });

    return { message: 'Баннер успешно удалён' };
  }
}
