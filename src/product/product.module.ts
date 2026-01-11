import { Module } from '@nestjs/common';
import { MulterModule } from '@nestjs/platform-express';

import { memoryStorage } from 'multer';
import { AuthModule } from 'src/auth/auth.module';
import { OptionalSessionAuthGuard } from 'src/auth/guards/optional-session-auth.guard';
import { ChatModule } from 'src/chat/chat.module';
import { PrismaModule } from 'src/prisma/prisma.module';
import { S3Module } from 'src/s3/s3.module';
import { UserModule } from 'src/user/user.module';

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
