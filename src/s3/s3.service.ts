import { BadRequestException, Injectable } from '@nestjs/common';
import { ConfigService } from '@nestjs/config';

import {
  DeleteObjectCommand,
  DeleteObjectsCommand,
  GetObjectCommand,
  PutObjectCommand,
  S3Client,
} from '@aws-sdk/client-s3';
import { v4 as uuidv4 } from 'uuid';

@Injectable()
export class S3Service {
  private s3Client: S3Client;
  private bucketName: string;

  constructor(private readonly configService: ConfigService) {
    this.bucketName = this.configService.get<string>(
      'S3_BUCKET_NAME',
      'not found',
    );

    const accessKey = this.configService.get<string>('S3_ACCESS_KEY');
    const secretKey = this.configService.get<string>('S3_SECRET_KEY');
    const endpoint = this.configService.get<string>(
      'S3_ENDPOINT',
      'https://s3.ru1.storage.beget.cloud',
    );
    const region = this.configService.get<string>('S3_REGION', 'ru1');

    console.log('S3 Configuration:');
    console.log('Bucket:', this.bucketName);
    console.log('Endpoint:', endpoint);
    console.log('Region:', region);
    console.log(
      'Access Key:',
      accessKey ? `${accessKey.substring(0, 4)}...` : 'NOT FOUND',
    );
    console.log(
      'Secret Key:',
      secretKey ? `${secretKey.substring(0, 4)}...` : 'NOT FOUND',
    );

    if (!accessKey || !secretKey) {
      throw new Error(
        'S3 credentials not found in environment variables. Check your .env file.',
      );
    }

    this.s3Client = new S3Client({
      region,
      endpoint,
      credentials: {
        accessKeyId: accessKey,
        secretAccessKey: secretKey,
      },
      forcePathStyle: true,
    });
  }

  /**
   * Загрузка файла в S3
   * @param file - файл для загрузки
   * @param folder - папка в бакете (например: 'products', 'users')
   * @returns URL загруженного файла
   */
  async uploadFile(
    file: Express.Multer.File,
    folder: string = 'uploads',
  ): Promise<string> {
    try {
      const fileExtension = file.originalname.split('.').pop();
      const fileName = `${folder}/${uuidv4()}.${fileExtension}`;

      const command = new PutObjectCommand({
        Bucket: this.bucketName,
        Key: fileName,
        Body: file.buffer,
        ContentType: file.mimetype,
        // ACL убран - используем "Сделать все материалы публичными" из панели Beget
      });

      await this.s3Client.send(command);

      // Path-style URL для совместимости с Beget
      return `https://s3.ru1.storage.beget.cloud/${this.bucketName}/${fileName}`;
    } catch (error) {
      throw new BadRequestException(
        `Ошибка загрузки файла в S3: ${error.message}`,
      );
    }
  }

  /**
   * Загрузка нескольких файлов
   * @param files - массив файлов
   * @param folder - папка в бакете
   * @returns массив URL загруженных файлов
   */
  async uploadFiles(
    files: Express.Multer.File[],
    folder: string = 'uploads',
  ): Promise<string[]> {
    try {
      const uploadPromises = files.map((file) => this.uploadFile(file, folder));
      return await Promise.all(uploadPromises);
    } catch (error) {
      throw new BadRequestException(
        `Ошибка загрузки файлов в S3: ${error.message}`,
      );
    }
  }

  /**
   * Удаление файла из S3
   * @param fileUrl - полный URL файла или путь (key)
   */
  async deleteFile(fileUrl: string): Promise<void> {
    try {
      // Извлекаем ключ файла из URL
      const key = this.extractKeyFromUrl(fileUrl);

      const command = new DeleteObjectCommand({
        Bucket: this.bucketName,
        Key: key,
      });

      await this.s3Client.send(command);
    } catch (error) {
      console.error(`Ошибка удаления файла из S3: ${error.message}`);
      // Не выбрасываем исключение, чтобы не блокировать основной функционал
    }
  }

  /**
   * Удаление нескольких файлов
   * @param fileUrls - массив URL файлов
   */
  async deleteFiles(fileUrls: string[]): Promise<void> {
    try {
      const keys = fileUrls.map((url) => this.extractKeyFromUrl(url));

      const command = new DeleteObjectsCommand({
        Bucket: this.bucketName,
        Delete: {
          Objects: keys.map((key) => ({ Key: key })),
        },
      });

      await this.s3Client.send(command);
    } catch (error) {
      console.error(`Ошибка удаления файлов из S3: ${error.message}`);
    }
  }

  /**
   * Получение файла из S3 (для приватных файлов)
   * @param key - ключ файла в S3
   * @returns поток данных файла
   */
  async getFile(key: string) {
    try {
      const command = new GetObjectCommand({
        Bucket: this.bucketName,
        Key: key,
      });

      const response = await this.s3Client.send(command);
      return response.Body;
    } catch (error) {
      throw new BadRequestException(
        `Ошибка получения файла из S3: ${error.message}`,
      );
    }
  }

  /**
   * Извлечение ключа файла из полного URL
   * @param fileUrl - URL файла
   * @returns ключ файла (путь в бакете)
   */
  private extractKeyFromUrl(fileUrl: string): string {
    // Если это уже ключ (не URL), возвращаем как есть
    if (!fileUrl.includes('http')) {
      return fileUrl;
    }

    // Извлекаем путь из URL
    // Например: https://bucket.s3.region.beget.cloud/folder/file.jpg -> folder/file.jpg
    const urlParts = fileUrl.split('.cloud/');
    return urlParts.length > 1 ? urlParts[1] : fileUrl;
  }
}
