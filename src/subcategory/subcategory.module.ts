import { Module } from '@nestjs/common';
import { SubcategoryService } from './subcategory.service';
import { SubcategoryController } from './subcategory.controller';
import { AuthModule } from 'src/auth/auth.module';

@Module({
  imports: [AuthModule],
  controllers: [SubcategoryController],
  providers: [SubcategoryService],
})
export class SubcategoryModule {}
