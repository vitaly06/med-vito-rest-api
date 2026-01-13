import { Module } from '@nestjs/common';

import { AuthModule } from 'src/auth/auth.module';

import { PromotionController } from './promotion.controller';
import { PromotionService } from './promotion.service';

@Module({
  imports: [AuthModule],
  controllers: [PromotionController],
  providers: [PromotionService],
})
export class PromotionModule {}
