import {
  BadRequestException,
  Injectable,
  NotFoundException,
} from '@nestjs/common';
import { PrismaService } from 'src/prisma/prisma.service';
import { CreateCategoryDto } from './dto/create-category.dto';
import { UpdateCategoryDto } from './dto/update-category.dto';

@Injectable()
export class CategoryService {
  constructor(private readonly prisma: PrismaService) {}

  async createCategory(dto: CreateCategoryDto) {
    const checkCategory = await this.prisma.category.findFirst({
      where: { name: dto.name },
    });

    if (checkCategory) {
      throw new BadRequestException('Данная категория уже существует');
    }

    await this.prisma.category.create({
      data: {
        name: dto.name,
      },
    });

    return { message: 'Категория успешно создана' };
  }

  async updateCategory(categoryId: number, dto: UpdateCategoryDto) {
    const checkCategory = await this.prisma.category.findUnique({
      where: { id: categoryId },
    });

    if (!checkCategory) {
      throw new NotFoundException('Категория с таким id не найдена');
    }

    await this.prisma.category.update({
      where: { id: categoryId },
      data: {
        name: dto.name,
      },
    });

    return { message: 'Категория успешно обновлена' };
  }

  async deleteCategory(categoryId: number) {
    const checkCategory = await this.prisma.category.findUnique({
      where: { id: categoryId },
    });

    if (!checkCategory) {
      throw new NotFoundException('Категория для удаления не найдена');
    }

    await this.prisma.category.delete({
      where: { id: categoryId },
    });

    return { message: 'Категория успешно удалена' };
  }

  async findAllCategories() {
    const categoires = await this.prisma.category.findMany({
      select: {
        id: true,
        name: true,
        subCategories: {
          select: {
            id: true,
            name: true,
            subcategoryTypes: {
              select: {
                id: true,
                name: true,
                fields: {
                  select: {
                    id: true,
                    name: true,
                  },
                },
              },
            },
          },
        },
      },
    });

    return categoires;
  }

  async findById(id: number) {
    const checkCategory = await this.prisma.category.findUnique({
      where: { id },
    });

    if (!checkCategory) {
      throw new NotFoundException('Категория не найдена');
    }

    return checkCategory;
  }
}
