import { Controller, Get } from '@nestjs/common';
import { CategoryService } from './category.service';
import { ApiOperation } from '@nestjs/swagger';

@Controller('category')
export class CategoryController {
  constructor(private readonly categoryService: CategoryService) {}

  @ApiOperation({
    summary: 'Получение всех категорий',
  })
  @Get('all-categories')
  async findAllCategories() {
    return await this.categoryService.findAllCategories();
  }

  @ApiOperation({
    summary: 'Получение всех подкатегорий',
  })
  @Get('all-subcategories')
  async findAllSubcategories() {
    return await this.categoryService.findAllSubcategories();
  }
}
