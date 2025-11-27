import { Module } from '@nestjs/common';
import { ProductService } from './product.service';
import { ProductController } from './product.controller';
import { diskStorage } from 'multer';
import * as path from 'path';
import { MulterModule } from '@nestjs/platform-express';
import { AuthModule } from 'src/auth/auth.module';
import { PrismaModule } from 'src/prisma/prisma.module';
import { OptionalSessionAuthGuard } from 'src/auth/guards/optional-session-auth.guard';
import { UserModule } from 'src/user/user.module';

@Module({
  imports: [
    MulterModule.register({
      storage: diskStorage({
        destination: './uploads/product',
        filename: (req, file, callback) => {
          const uniqueSuffix =
            Date.now() + '-' + Math.round(Math.random() * 1e9);
          const ext = path.extname(file.originalname);
          callback(null, file.fieldname + '-' + uniqueSuffix + ext);
        },
      }),
    }),
    AuthModule,
    PrismaModule,
    UserModule,
  ],
  controllers: [ProductController],
  providers: [ProductService, OptionalSessionAuthGuard],
})
export class ProductModule {}
