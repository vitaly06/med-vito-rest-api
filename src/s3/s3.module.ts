import { Module } from '@nestjs/common';
import { MulterModule } from '@nestjs/platform-express';

import { memoryStorage } from 'multer';

import { S3Controller } from './s3.controller';
import { S3Service } from './s3.service';

@Module({
  imports: [
    MulterModule.register({
      storage: memoryStorage(), // Используем память для S3
    }),
  ],
  controllers: [S3Controller],
  providers: [S3Service],
  exports: [S3Service],
})
export class S3Module {}
