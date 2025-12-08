import { Injectable, NotFoundException } from '@nestjs/common';
import { PrismaService } from 'src/prisma/prisma.service';
import { CreateSubcategoryTypeDto } from './dto/create-subcategory-type.dto';
import { UpdateSubcategoryTypeDto } from './dto/update-subcategory-type.dto';

@Injectable()
export class SubcategoryTypeService {
  constructor(private readonly prisma: PrismaService) {}

  async create(dto: CreateSubcategoryTypeDto) {
    const { name, subcategoryId } = { ...dto };

    const checkSubcategory = await this.prisma.subCategory.findUnique({
      where: { id: subcategoryId },
    });

    if (!checkSubcategory) {
      throw new NotFoundException('Подкатегория не найдена');
    }

    await this.prisma.subcategotyType.create({
      data: {
        name,
        subcategoryId,
      },
    });

    return { message: 'Тип подкатегории успешно создан' };
  }

  async update(id: number, dto: UpdateSubcategoryTypeDto) {
    if (dto.subcategoryId) {
      const checkSubcategory = await this.prisma.subCategory.findUnique({
        where: { id },
      });

      if (!checkSubcategory) {
        throw new NotFoundException('подкатегория не найден');
      }
    }

    await this.prisma.subcategotyType.update({
      where: { id },
      data: {
        ...dto,
      },
    });

    return { message: 'Тип подкатегории успешно обновлен' };
  }

  async delete(id: number) {
    const checkSubcategoryType = await this.prisma.subcategotyType.findUnique({
      where: { id },
    });

    if (!checkSubcategoryType) {
      throw new NotFoundException('Тип подкатегории для удаления не надйен');
    }

    await this.prisma.subcategotyType.delete({
      where: { id },
    });

    return { message: 'Тип подкатегории успешно удалён' };
  }

  async findAll() {
    const subcategoryTypes = await this.prisma.subcategotyType.findMany({
      select: {
        id: true,
        name: true,
        subcategoryId: true,
      },
    });

    return subcategoryTypes;
  }

  async findById(id: number) {
    const subcategoryType = await this.prisma.subcategotyType.findUnique({
      where: { id },
    });

    if (!subcategoryType) {
      throw new NotFoundException('Тип подкатегории не найден');
    }

    return subcategoryType;
  }
}
