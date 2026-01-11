import {
  BadRequestException,
  Injectable,
  NotFoundException,
} from '@nestjs/common';

import { generateSlug, makeUniqueSlug } from '@/common/utils/slug.utils';
import { PrismaService } from '@/prisma/prisma.service';

import { CreateCategoryDto } from './dto/create-category.dto';
import { UpdateCategoryDto } from './dto/update-category.dto';

@Injectable()
export class CategoryService {
  constructor(private readonly prisma: PrismaService) {}

  async createCategory(dto: CreateCategoryDto) {
    // Генерируем slug из названия или используем предоставленный
    let slug = dto.slug || generateSlug(dto.name);

    // Проверяем уникальность slug
    const existingCategory = await this.prisma.category.findUnique({
      where: { slug },
    });

    if (existingCategory) {
      if (dto.slug) {
        // Если slug был предоставлен явно, выбрасываем ошибку
        throw new BadRequestException(
          `Категория с slug '${slug}' уже существует`,
        );
      }
      // Иначе генерируем уникальный slug
      const existingSlugs = (await this.prisma.category.findMany()).map(
        (c) => c.slug,
      );
      slug = makeUniqueSlug(slug, existingSlugs);
    }

    const category = await this.prisma.category.create({
      data: {
        name: dto.name,
        slug,
      },
    });

    return {
      message: 'Категория успешно создана',
      category,
    };
  }

  async updateCategory(categoryId: number, dto: UpdateCategoryDto) {
    const checkCategory = await this.prisma.category.findUnique({
      where: { id: categoryId },
    });

    if (!checkCategory) {
      throw new NotFoundException('Категория с таким id не найдена');
    }

    // Генерируем новый slug если изменилось название и slug не указан явно
    let slug = dto.slug;
    if (!slug && dto.name !== checkCategory.name) {
      slug = generateSlug(dto.name);

      // Проверяем уникальность нового slug
      const existingCategory = await this.prisma.category.findFirst({
        where: {
          slug,
          NOT: { id: categoryId },
        },
      });

      if (existingCategory) {
        const existingSlugs = (
          await this.prisma.category.findMany({
            where: { NOT: { id: categoryId } },
          })
        ).map((c) => c.slug);
        slug = makeUniqueSlug(slug, existingSlugs);
      }
    } else if (slug) {
      // Если slug указан явно, проверяем его уникальность
      const existingCategory = await this.prisma.category.findFirst({
        where: {
          slug,
          NOT: { id: categoryId },
        },
      });

      if (existingCategory) {
        throw new BadRequestException(
          `Категория с slug '${slug}' уже существует`,
        );
      }
    }

    const category = await this.prisma.category.update({
      where: { id: categoryId },
      data: {
        name: dto.name,
        ...(slug && { slug }),
      },
    });

    return {
      message: 'Категория успешно обновлена',
      category,
    };
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
    const categories = await this.prisma.category.findMany({
      select: {
        id: true,
        name: true,
        slug: true,
        subCategories: {
          select: {
            id: true,
            name: true,
            slug: true,
            subcategoryTypes: {
              select: {
                id: true,
                name: true,
                slug: true,
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
      orderBy: {
        name: 'asc',
      },
    });

    return categories;
  }

  async findById(id: number) {
    const category = await this.prisma.category.findUnique({
      where: { id },
      include: {
        subCategories: {
          include: {
            subcategoryTypes: {
              include: {
                fields: true,
              },
            },
          },
        },
      },
    });

    if (!category) {
      throw new NotFoundException('Категория не найдена');
    }

    return category;
  }

  /**
   * Получить категорию по slug
   */
  async findBySlug(slug: string) {
    const category = await this.prisma.category.findUnique({
      where: { slug },
      include: {
        subCategories: {
          include: {
            subcategoryTypes: {
              include: {
                fields: true,
              },
            },
          },
        },
      },
    });

    if (!category) {
      throw new NotFoundException(`Категория '${slug}' не найдена`);
    }

    return category;
  }

  /**
   * Получить категорию по цепочке slug'ов
   * Например: elektronika/telefony/smartfony
   */
  async findBySlugPath(slugPath: string) {
    const slugs = slugPath.split('/').filter((s) => s.length > 0);

    if (slugs.length === 0) {
      throw new BadRequestException('Некорректный путь категории');
    }

    // Первый slug - это категория
    const categorySlug = slugs[0];
    const category = await this.prisma.category.findUnique({
      where: { slug: categorySlug },
      include: {
        subCategories: {
          include: {
            subcategoryTypes: {
              include: {
                fields: true,
              },
            },
          },
        },
      },
    });

    if (!category) {
      throw new NotFoundException(`Категория '${categorySlug}' не найдена`);
    }

    // Если путь состоит только из категории
    if (slugs.length === 1) {
      return {
        type: 'category',
        data: category,
      };
    }

    // Второй slug - это подкатегория
    const subCategorySlug = slugs[1];
    const subCategory = category.subCategories.find(
      (sc) => sc.slug === subCategorySlug,
    );

    if (!subCategory) {
      throw new NotFoundException(
        `Подкатегория '${subCategorySlug}' не найдена в категории '${categorySlug}'`,
      );
    }

    // Если путь состоит из категории и подкатегории
    if (slugs.length === 2) {
      return {
        type: 'subcategory',
        data: await this.prisma.subCategory.findUnique({
          where: { id: subCategory.id },
          include: {
            category: true,
            subcategoryTypes: {
              include: {
                fields: true,
              },
            },
          },
        }),
      };
    }

    // Третий slug - это тип подкатегории
    const typeSlug = slugs[2];
    const subCategoryType = subCategory.subcategoryTypes.find(
      (st) => st.slug === typeSlug,
    );

    if (!subCategoryType) {
      throw new NotFoundException(
        `Тип '${typeSlug}' не найден в подкатегории '${subCategorySlug}'`,
      );
    }

    return {
      type: 'subcategory-type',
      data: await this.prisma.subcategotyType.findUnique({
        where: { id: subCategoryType.id },
        include: {
          subcategory: {
            include: {
              category: true,
            },
          },
          fields: true,
        },
      }),
    };
  }

  /**
   * Получить полный slug path для товара по его ID категорий
   */
  async getSlugPathForProduct(
    categoryId: number,
    subCategoryId?: number,
    typeId?: number,
  ): Promise<string> {
    const parts: string[] = [];

    // Получаем категорию
    const category = await this.prisma.category.findUnique({
      where: { id: categoryId },
      select: { slug: true },
    });
    if (category) {
      parts.push(category.slug);
    }

    // Получаем подкатегорию
    if (subCategoryId) {
      const subCategory = await this.prisma.subCategory.findUnique({
        where: { id: subCategoryId },
        select: { slug: true },
      });
      if (subCategory) {
        parts.push(subCategory.slug);
      }
    }

    // Получаем тип
    if (typeId) {
      const type = await this.prisma.subcategotyType.findUnique({
        where: { id: typeId },
        select: { slug: true },
      });
      if (type) {
        parts.push(type.slug);
      }
    }

    return parts.join('/');
  }
}
