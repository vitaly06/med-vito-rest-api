import { Injectable, NotFoundException } from '@nestjs/common';
import { PrismaService } from 'src/prisma/prisma.service';
import { S3Service } from 'src/s3/s3.service';

@Injectable()
export class BannerService {
  constructor(
    private readonly prisma: PrismaService,
    private readonly s3Service: S3Service,
  ) {}

  async create(file: Express.Multer.File) {
    // Загружаем изображение в S3
    const photoUrl = await this.s3Service.uploadFile(file, 'banners');

    // Создаём баннер в БД
    const banner = await this.prisma.banner.create({
      data: {
        photoUrl,
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

  async findOne(id: number) {
    const banner = await this.prisma.banner.findUnique({
      where: { id },
    });

    if (!banner) {
      throw new NotFoundException(`Баннер с ID ${id} не найден`);
    }

    return banner;
  }

  async update(id: number, file?: Express.Multer.File) {
    // Проверяем существование баннера
    const existingBanner = await this.findOne(id);

    let updateData: any = {};

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
