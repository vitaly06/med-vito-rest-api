import { Module } from '@nestjs/common';

import { AuthModule } from 'src/auth/auth.module';

import { SubcategoryTypeController } from './subcategory-type.controller';
import { SubcategoryTypeService } from './subcategory-type.service';

@Module({
  imports: [AuthModule],
  controllers: [SubcategoryTypeController],
  providers: [SubcategoryTypeService],
})
export class SubcategoryTypeModule {}
