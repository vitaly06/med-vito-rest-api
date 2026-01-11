import {
  BadRequestException,
  Body,
  Controller,
  Delete,
  Post,
  UploadedFile,
  UploadedFiles,
  UseInterceptors,
} from '@nestjs/common';
import { FileInterceptor, FilesInterceptor } from '@nestjs/platform-express';
import { ApiBody, ApiConsumes, ApiOperation, ApiTags } from '@nestjs/swagger';

import { S3Service } from './s3.service';

@Controller('s3')
export class S3Controller {
  constructor(private readonly s3Service: S3Service) {}

  // @Post('upload')
  // @UseInterceptors(FileInterceptor('file'))
  // @ApiOperation({ summary: 'Загрузка одного файла в S3' })
  // @ApiConsumes('multipart/form-data')
  // @ApiBody({
  //   schema: {
  //     type: 'object',
  //     properties: {
  //       file: {
  //         type: 'string',
  //         format: 'binary',
  //       },
  //       folder: {
  //         type: 'string',
  //         description: 'Папка в бакете (опционально)',
  //         example: 'products',
  //       },
  //     },
  //   },
  // })
  // async uploadFile(
  //   @UploadedFile() file: Express.Multer.File,
  //   @Body('folder') folder?: string,
  // ) {
  //   if (!file) {
  //     throw new BadRequestException('Файл не предоставлен');
  //   }

  //   const fileUrl = await this.s3Service.uploadFile(file, folder || 'uploads');

  //   return {
  //     message: 'Файл успешно загружен',
  //     url: fileUrl,
  //   };
  // }

  // @Post('upload-multiple')
  // @UseInterceptors(FilesInterceptor('files', 10))
  // @ApiOperation({ summary: 'Загрузка нескольких файлов в S3 (до 10)' })
  // @ApiConsumes('multipart/form-data')
  // @ApiBody({
  //   schema: {
  //     type: 'object',
  //     properties: {
  //       files: {
  //         type: 'array',
  //         items: {
  //           type: 'string',
  //           format: 'binary',
  //         },
  //       },
  //       folder: {
  //         type: 'string',
  //         description: 'Папка в бакете (опционально)',
  //         example: 'products',
  //       },
  //     },
  //   },
  // })
  // async uploadFiles(
  //   @UploadedFiles() files: Express.Multer.File[],
  //   @Body('folder') folder?: string,
  // ) {
  //   if (!files || files.length === 0) {
  //     throw new BadRequestException('Файлы не предоставлены');
  //   }

  //   const fileUrls = await this.s3Service.uploadFiles(
  //     files,
  //     folder || 'uploads',
  //   );

  //   return {
  //     message: `Успешно загружено ${files.length} файлов`,
  //     urls: fileUrls,
  //   };
  // }

  // @Delete('delete')
  // @ApiOperation({ summary: 'Удаление файла из S3' })
  // @ApiBody({
  //   schema: {
  //     type: 'object',
  //     properties: {
  //       fileUrl: {
  //         type: 'string',
  //         description: 'URL или ключ файла для удаления',
  //         example:
  //           'https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/file.jpg',
  //       },
  //     },
  //   },
  // })
  // async deleteFile(@Body('fileUrl') fileUrl: string) {
  //   if (!fileUrl) {
  //     throw new BadRequestException('URL файла не предоставлен');
  //   }

  //   await this.s3Service.deleteFile(fileUrl);

  //   return {
  //     message: 'Файл успешно удалён',
  //   };
  // }

  // @Delete('delete-multiple')
  // @ApiOperation({ summary: 'Удаление нескольких файлов из S3' })
  // @ApiBody({
  //   schema: {
  //     type: 'object',
  //     properties: {
  //       fileUrls: {
  //         type: 'array',
  //         items: { type: 'string' },
  //         description: 'Массив URL или ключей файлов для удаления',
  //         example: [
  //           'https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/file1.jpg',
  //           'https://c15b4d655f70-medvito-data.s3.ru1.storage.beget.cloud/products/file2.jpg',
  //         ],
  //       },
  //     },
  //   },
  // })
  // async deleteFiles(@Body('fileUrls') fileUrls: string[]) {
  //   if (!fileUrls || fileUrls.length === 0) {
  //     throw new BadRequestException('URL файлов не предоставлены');
  //   }

  //   await this.s3Service.deleteFiles(fileUrls);

  //   return {
  //     message: `Успешно удалено ${fileUrls.length} файлов`,
  //   };
  // }
}
