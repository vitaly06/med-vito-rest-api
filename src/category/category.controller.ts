import {
  Body,
  Controller,
  Delete,
  Get,
  Param,
  Post,
  Put,
  UseGuards,
} from '@nestjs/common';
import { ApiOperation, ApiParam } from '@nestjs/swagger';

import { AdminSessionAuthGuard } from '@/auth/guards/admin-session-auth.guard';

import { CategoryService } from './category.service';
import { CreateCategoryDto } from './dto/create-category.dto';
import { UpdateCategoryDto } from './dto/update-category.dto';

@Controller('category')
export class CategoryController {
  constructor(private readonly categoryService: CategoryService) {}

  @ApiOperation({
    summary: 'Создание категории',
  })
  @UseGuards(AdminSessionAuthGuard)
  @Post('create-category')
  async createCategory(@Body() dto: CreateCategoryDto) {
    return await this.categoryService.createCategory(dto);
  }

  @ApiOperation({
    summary: 'Обновление категории',
  })
  @UseGuards(AdminSessionAuthGuard)
  @Put('update-category/:id')
  async updateCategory(
    @Param('id') id: string,
    @Body() dto: UpdateCategoryDto,
  ) {
    return await this.categoryService.updateCategory(+id, dto);
  }

  @ApiOperation({
    summary: 'Удаление категории',
  })
  @UseGuards(AdminSessionAuthGuard)
  @Delete('delete-category/:id')
  async deleteCategory(@Param('id') id: string) {
    return await this.categoryService.deleteCategory(+id);
  }

  @ApiOperation({
    summary: 'Все категории',
  })
  @Get('find-all')
  async findAllCategories() {
    return await this.categoryService.findAllCategories();
  }

  @ApiOperation({
    summary: 'Категория по id',
  })
  @Get('find-by-id/:id')
  async findById(@Param('id') id: string) {
    return await this.categoryService.findById(+id);
  }

  @ApiOperation({
    summary: 'Получить категорию по slug',
    description: 'Возвращает категорию с подкатегориями по её slug',
  })
  @ApiParam({
    name: 'slug',
    description: 'Slug категории',
    example: 'elektronika',
  })
  @Get('slug/:slug')
  async findBySlug(@Param('slug') slug: string) {
    return await this.categoryService.findBySlug(slug);
  }

  @ApiOperation({
    summary: 'Получить категорию по цепочке slug (slugPath)',
    description:
      'Возвращает категорию, подкатегорию или тип по полному пути. ' +
      'Примеры: path/elektronika, path/elektronika/telefony, path/elektronika/telefony/smartfony',
  })
  @ApiParam({
    name: 'slugPath',
    description: 'Полный путь категории через /',
    example: 'elektronika/telefony/smartfony',
  })
  @Get('path/*slugPath')
  async findBySlugPath(@Param('slugPath') slugPath: string) {
    return await this.categoryService.findBySlugPath(slugPath);
  }
}
