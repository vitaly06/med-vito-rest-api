import { Module } from '@nestjs/common';

import { AuthModule } from '@/auth/auth.module';
import { PrismaModule } from '@/prisma/prisma.module';
import { S3Module } from '@/s3/s3.module';

import { BannerController } from './banner.controller';
import { BannerService } from './banner.service';

@Module({
  imports: [PrismaModule, S3Module, AuthModule],
  controllers: [BannerController],
  providers: [BannerService],
  exports: [BannerService],
})
export class BannerModule {}
