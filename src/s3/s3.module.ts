import { Module } from '@nestjs/common';
import { S3Service } from './s3.service';
import { MulterModule } from '@nestjs/platform-express';
import { memoryStorage } from 'multer';

@Module({
  imports: [
    MulterModule.register({
      storage: memoryStorage(), // Используем память для S3
    }),
  ],
  providers: [S3Service],
  exports: [S3Service],
})
export class S3Module {}
