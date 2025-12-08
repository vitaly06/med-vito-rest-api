import { Module } from '@nestjs/common';
import { SubcategoryTypeService } from './subcategory-type.service';
import { SubcategoryTypeController } from './subcategory-type.controller';
import { AuthModule } from 'src/auth/auth.module';

@Module({
  imports: [AuthModule],
  controllers: [SubcategoryTypeController],
  providers: [SubcategoryTypeService],
})
export class SubcategoryTypeModule {}
