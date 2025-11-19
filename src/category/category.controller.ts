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
import { CategoryService } from './category.service';
import { ApiOperation } from '@nestjs/swagger';
import { CreateCategoryDto } from './dto/create-category.dto';
import { AdminJwtAuthGuard } from 'src/auth/guards/admin-jwt-auth.guard';
import { UpdateCategoryDto } from './dto/update-category.dto';

@Controller('category')
export class CategoryController {
  constructor(private readonly categoryService: CategoryService) {}

  @ApiOperation({
    summary: 'Создание категории',
  })
  @UseGuards(AdminJwtAuthGuard)
  @Post('create-category')
  async createCategory(@Body() dto: CreateCategoryDto) {
    return await this.categoryService.createCategory(dto);
  }

  @ApiOperation({
    summary: 'Обновление категории',
  })
  @UseGuards(AdminJwtAuthGuard)
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
  @UseGuards(AdminJwtAuthGuard)
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
}
