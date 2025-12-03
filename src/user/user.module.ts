import { Module } from '@nestjs/common';
import { UserService } from './user.service';
import { UserController } from './user.controller';
import { AuthModule } from 'src/auth/auth.module';
import { PrismaModule } from 'src/prisma/prisma.module';
import { MulterModule } from '@nestjs/platform-express';
import { memoryStorage } from 'multer';
import { S3Module } from 'src/s3/s3.module';

@Module({
  imports: [
    AuthModule,
    PrismaModule,
    S3Module,
    MulterModule.register({
      storage: memoryStorage(), // Используем память для S3
    }),
  ],
  controllers: [UserController],
  providers: [UserService],
  exports: [UserService],
})
export class UserModule {}
