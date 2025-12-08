import { Controller, UseGuards } from '@nestjs/common';
import { SubcategoryTypeService } from './subcategory-type.service';
import { ApiOperation } from '@nestjs/swagger';
import { AdminSessionAuthGuard } from 'src/auth/guards/admin-session-auth.guard';

@Controller('subcategory-type')
export class SubcategoryTypeController {
  constructor(
    private readonly subcategoryTypeService: SubcategoryTypeService,
  ) {}

  @ApiOperation({
    summary: 'Создание типа подкатегории',
  })
  @UseGuards(AdminSessionAuthGuard)
  async createsubcategoryType() {}
}
