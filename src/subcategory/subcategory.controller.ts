import {
  Body,
  Controller,
  Delete,
  Get,
  Param,
  Patch,
  Post,
  UseGuards,
} from '@nestjs/common';
import { SubcategoryService } from './subcategory.service';
import { ApiOperation } from '@nestjs/swagger';
import { AdminJwtAuthGuard } from 'src/auth/guards/admin-jwt-auth.guard';
import { CreateSubcategoryDto } from './dto/create-subcategory.dto';
import { UpdateSubcategoryDto } from './dto/update-subcategory.dto';

@Controller('subcategory')
export class SubcategoryController {
  constructor(private readonly subcategoryService: SubcategoryService) {}

  @ApiOperation({
    summary: 'Создание подкатегории',
  })
  @UseGuards(AdminJwtAuthGuard)
  @Post('create-subcategory')
  async createSubcategory(@Body() dto: CreateSubcategoryDto) {
    return await this.subcategoryService.createSubcategory(dto);
  }

  @ApiOperation({
    summary: 'Обновление подкатегории',
  })
  @UseGuards(AdminJwtAuthGuard)
  @Patch('update-subcategory/:id')
  async updateCategory(
    @Param('id') id: string,
    @Body() dto: UpdateSubcategoryDto,
  ) {
    return await this.subcategoryService.updateSubcategory(+id, dto);
  }

  @ApiOperation({
    summary: 'Удаление подкатегории',
  })
  @UseGuards(AdminJwtAuthGuard)
  @Delete('delete-subcategory/:id')
  async deleteSubcategory(@Param('id') id: string) {
    return await this.subcategoryService.deleteSubcategory(+id);
  }

  @ApiOperation({
    summary: 'Все подкатегории',
  })
  @Get('find-all')
  async findAll() {
    return await this.subcategoryService.findAll();
  }

  @ApiOperation({
    summary: 'Подкатегория по id',
  })
  @Get('find-by-id/:id')
  async findById(@Param('id') id: string) {
    return await this.subcategoryService.findById(+id);
  }
}
