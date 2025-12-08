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
import { SubcategoryTypeService } from './subcategory-type.service';
import { ApiOperation } from '@nestjs/swagger';
import { AdminSessionAuthGuard } from 'src/auth/guards/admin-session-auth.guard';
import { CreateSubcategoryTypeDto } from './dto/create-subcategory-type.dto';
import { UpdateSubcategoryTypeDto } from './dto/update-subcategory-type.dto';

@Controller('subcategory-type')
export class SubcategoryTypeController {
  constructor(
    private readonly subcategoryTypeService: SubcategoryTypeService,
  ) {}

  @ApiOperation({
    summary: 'Создание типа подкатегории',
  })
  @UseGuards(AdminSessionAuthGuard)
  @Post('create')
  async createsubcategoryType(@Body() dto: CreateSubcategoryTypeDto) {
    return await this.subcategoryTypeService.create(dto);
  }

  @ApiOperation({
    summary: 'Обновление типа подкатегории',
  })
  @UseGuards(AdminSessionAuthGuard)
  @Patch(':id')
  async updateSubcategoryType(
    @Param('id') id: string,
    @Body() dto: UpdateSubcategoryTypeDto,
  ) {
    return await this.subcategoryTypeService.update(+id, dto);
  }

  @ApiOperation({
    summary: 'Удаление типа подкатегории',
  })
  @UseGuards(AdminSessionAuthGuard)
  @Delete(':id')
  async deleteSubcategoryType(@Param('id') id: string) {
    return await this.subcategoryTypeService.delete(+id);
  }

  @ApiOperation({
    summary: 'Получение всех типов подкатегорий',
  })
  @Get('find-all')
  async findAll() {
    return await this.subcategoryTypeService.findAll();
  }

  @ApiOperation({
    summary: 'Получение типа подкатегории по id',
  })
  @Get('find-by-id/:id')
  async findById(@Param('id') id: string) {
    return await this.subcategoryTypeService.findById(+id);
  }
}
// slug
