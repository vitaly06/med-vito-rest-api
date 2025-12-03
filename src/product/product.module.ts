import { Module } from '@nestjs/common';
import { ProductService } from './product.service';
import { ProductController } from './product.controller';
import { MulterModule } from '@nestjs/platform-express';
import { AuthModule } from 'src/auth/auth.module';
import { PrismaModule } from 'src/prisma/prisma.module';
import { OptionalSessionAuthGuard } from 'src/auth/guards/optional-session-auth.guard';
import { UserModule } from 'src/user/user.module';
import { S3Module } from 'src/s3/s3.module';

@Module({
  imports: [
    MulterModule.register({
      storage: 'memory', // Используем память вместо диска для S3
    }),
    AuthModule,
    PrismaModule,
    UserModule,
    S3Module,
  ],
  controllers: [ProductController],
  providers: [ProductService, OptionalSessionAuthGuard],
})
export class ProductModule {}
