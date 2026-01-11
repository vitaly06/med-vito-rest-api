import { Injectable, NotFoundException } from '@nestjs/common';

import { PrismaService } from '@/prisma/prisma.service';

import { CreateTypeFieldDto } from './dto/create-type-field.dto';
import { UpdateTypeFieldDto } from './dto/update-type-field.dto';

@Injectable()
export class TypeFieldService {
  constructor(private readonly prisma: PrismaService) {}

  async create(dto: CreateTypeFieldDto) {
    const checkType = await this.prisma.subcategotyType.findUnique({
      where: { id: dto.typeId },
    });

    if (!checkType) {
      throw new NotFoundException('Тип подкатегории не найден');
    }

    await this.prisma.typeField.create({
      data: {
        name: dto.name,
        typeId: dto.typeId,
      },
    });

    return { message: 'Характеристика успешно добавлена' };
  }

  async update(id: number, dto: UpdateTypeFieldDto) {
    if (dto.typeId) {
      const checkType = await this.prisma.subcategotyType.findUnique({
        where: { id: dto.typeId },
      });

      if (!checkType) {
        throw new NotFoundException('Тип подкатегории не найден');
      }
    }

    await this.prisma.typeField.update({
      where: { id },
      data: {
        ...dto,
      },
    });

    return { message: 'Характеристика успешно обновлена' };
  }

  async delete(id: number) {
    const checkTypeField = await this.prisma.typeField.findUnique({
      where: { id },
    });

    if (!checkTypeField) {
      throw new NotFoundException('Характеристика для удаления не найдена');
    }

    await this.prisma.typeField.delete({
      where: { id },
    });

    return { message: 'Характеристика успешно удалена' };
  }

  async findAll() {
    const typeFields = await this.prisma.typeField.findMany({
      select: {
        id: true,
        name: true,
        typeId: true,
      },
    });

    return typeFields;
  }

  async findById(id: number) {
    const typeField = await this.prisma.typeField.findUnique({
      where: { id },
      select: {
        id: true,
        name: true,
        typeId: true,
      },
    });

    if (!typeField) {
      throw new NotFoundException('Характеристика не найдена');
    }

    return typeField;
  }
}
