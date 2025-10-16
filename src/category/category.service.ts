import { Injectable } from '@nestjs/common';
import { PrismaService } from 'src/prisma/prisma.service';

@Injectable()
export class CategoryService {
  constructor(private readonly prisma: PrismaService) {}

  async findAllCategories() {
    return await this.prisma.category.findMany();
  }

  async findAllSubcategories() {
    return await this.prisma.subCategory.findMany();
  }
}
