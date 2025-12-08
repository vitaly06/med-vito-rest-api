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
import { TypeFieldService } from './type-field.service';
import { ApiOperation } from '@nestjs/swagger';
import { AdminSessionAuthGuard } from 'src/auth/guards/admin-session-auth.guard';
import { CreateTypeFieldDto } from './dto/create-type-field.dto';
import { UpdateTypeFieldDto } from './dto/update-type-field.dto';

@Controller('type-field')
export class TypeFieldController {
  constructor(private readonly typeFieldService: TypeFieldService) {}

  @ApiOperation({
    summary: 'Создание характеристики',
  })
  @UseGuards(AdminSessionAuthGuard)
  @Post('create')
  async create(@Body() dto: CreateTypeFieldDto) {
    return await this.typeFieldService.create(dto);
  }

  @ApiOperation({
    summary: 'Обновление характеристики',
  })
  @UseGuards(AdminSessionAuthGuard)
  @Patch(':id')
  async update(@Param('id') id: string, dto: UpdateTypeFieldDto) {
    return await this.typeFieldService.update(+id, dto);
  }

  @ApiOperation({
    summary: 'Удаление характеристики',
  })
  @UseGuards(AdminSessionAuthGuard)
  @Delete(':id')
  async delete(@Param('id') id: string) {
    return await this.typeFieldService.delete(+id);
  }

  @ApiOperation({
    summary: 'Все характеристики',
  })
  @Get('find-all')
  async findAll() {
    return await this.typeFieldService.findAll();
  }

  @ApiOperation({
    summary: 'Получение характеристики по id',
  })
  @Get('find-by-id/:id')
  async findById(@Param('id') id: string) {
    return await this.typeFieldService.findById(+id);
  }
}
