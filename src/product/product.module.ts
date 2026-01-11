import { Module } from '@nestjs/common';
import { MulterModule } from '@nestjs/platform-express';

import { memoryStorage } from 'multer';

import { AuthModule } from '@/auth/auth.module';
import { OptionalSessionAuthGuard } from '@/auth/guards/optional-session-auth.guard';
import { ChatModule } from '@/chat/chat.module';
import { PrismaModule } from '@/prisma/prisma.module';
import { S3Module } from '@/s3/s3.module';
import { UserModule } from '@/user/user.module';

import { ProductController } from './product.controller';
import { ProductService } from './product.service';

@Module({
  imports: [
    MulterModule.register({
      storage: memoryStorage(), // Используем память вместо диска для S3
    }),
    AuthModule,
    PrismaModule,
    UserModule,
    S3Module,
    ChatModule,
  ],
  controllers: [ProductController],
  providers: [ProductService, OptionalSessionAuthGuard],
})
export class ProductModule {}
