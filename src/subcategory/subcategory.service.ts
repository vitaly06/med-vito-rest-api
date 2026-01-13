import {
  BadRequestException,
  Injectable,
  NotFoundException,
} from '@nestjs/common';

import { PrismaService } from 'src/prisma/prisma.service';

import { CreateSubcategoryDto } from './dto/create-subcategory.dto';
import { UpdateSubcategoryDto } from './dto/update-subcategory.dto';

@Injectable()
export class SubcategoryService {
  constructor(private readonly prisma: PrismaService) {}

  async createSubcategory(dto: CreateSubcategoryDto) {
    const { name, categoryId } = { ...dto };

    const checkCategory = await this.prisma.category.findUnique({
      where: { id: categoryId },
    });

    if (!checkCategory) {
      throw new BadRequestException('Выбранная категория не существует');
    }

    await this.prisma.subCategory.create({
      data: {
        name,
        categoryId,
      },
    });

    return { message: 'Подкатегория успешно создана' };
  }

  async updateSubcategory(subcategoryId: number, dto: UpdateSubcategoryDto) {
    if (dto.categoryId) {
      const checkCategory = await this.prisma.category.findUnique({
        where: { id: dto.categoryId },
      });

      if (!checkCategory) {
        throw new BadRequestException('Выбранная категория не существует');
      }
    }

    await this.prisma.subCategory.update({
      where: { id: subcategoryId },
      data: {
        ...dto,
      },
    });

    return { message: 'Подкатегория успешно обновлена' };
  }

  async deleteSubcategory(id: number) {
    const checkSubcategory = await this.prisma.subCategory.findUnique({
      where: { id },
    });

    if (!checkSubcategory) {
      throw new BadRequestException('Подкатегория для удаления не найдена');
    }

    await this.prisma.subCategory.delete({
      where: { id },
    });

    return { message: 'Подкатегория успешно удалена' };
  }

  async findAll() {
    const subcategories = await this.prisma.subCategory.findMany({
      select: {
        id: true,
        name: true,
        categoryId: true,
      },
    });

    return subcategories;
  }

  async findById(id: number) {
    const subcategory = await this.prisma.subCategory.findUnique({
      where: { id },
    });

    if (!subcategory) {
      throw new NotFoundException('Подкатегория не найдена');
    }

    return subcategory;
  }
}
