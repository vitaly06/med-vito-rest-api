import { PrismaClient } from '@prisma/client';
import { S3Client, PutObjectCommand } from '@aws-sdk/client-s3';
import * as fs from 'fs';
import * as path from 'path';
import { v4 as uuidv4 } from 'uuid';

// Конфигурация S3
const S3_CONFIG = {
  endpoint: 'https://s3.ru1.storage.beget.cloud',
  region: 'ru1',
  bucketName: 'c15b4d655f70-medvito-data',
  credentials: {
    accessKeyId: 'I6I3KOJ2YO3TN08TDJAI',
    secretAccessKey: '5up6F9kLNHRGmPIczdqAVZgBNgKhFpAGJ1JnCJUY',
  },
};

const prisma = new PrismaClient();
const s3Client = new S3Client({
  region: S3_CONFIG.region,
  endpoint: S3_CONFIG.endpoint,
  credentials: S3_CONFIG.credentials,
  forcePathStyle: true,
});

interface MigrationStats {
  totalProducts: number;
  totalImages: number;
  uploadedImages: number;
  failedImages: number;
  updatedProducts: number;
  totalUsers: number;
  uploadedUserPhotos: number;
  failedUserPhotos: number;
  updatedUsers: number;
  errors: Array<{
    type: 'product' | 'user';
    id: number;
    path: string;
    error: string;
  }>;
}

/**
 * Загружает файл из локальной папки в S3
 */
async function uploadFileToS3(
  localPath: string,
  folder: string = 'products',
): Promise<string> {
  try {
    // Читаем файл
    const fileBuffer = fs.readFileSync(localPath);
    const fileName = path.basename(localPath);
    const fileExtension = path.extname(fileName);

    // Генерируем уникальное имя для файла в S3
    const s3FileName = `${folder}/${uuidv4()}${fileExtension}`;

    // Определяем Content-Type
    const contentType = getContentType(fileExtension);

    // Загружаем в S3
    const command = new PutObjectCommand({
      Bucket: S3_CONFIG.bucketName,
      Key: s3FileName,
      Body: fileBuffer,
      ContentType: contentType,
      ACL: 'public-read',
    });

    await s3Client.send(command);

    // Возвращаем публичный URL
    return `https://${S3_CONFIG.bucketName}.s3.ru1.storage.beget.cloud/${s3FileName}`;
  } catch (error) {
    throw new Error(`Failed to upload file: ${error.message}`);
  }
}

/**
 * Определяет Content-Type по расширению файла
 */
function getContentType(extension: string): string {
  const mimeTypes: { [key: string]: string } = {
    '.jpg': 'image/jpeg',
    '.jpeg': 'image/jpeg',
    '.png': 'image/png',
    '.gif': 'image/gif',
    '.webp': 'image/webp',
    '.jfif': 'image/jpeg',
  };
  return mimeTypes[extension.toLowerCase()] || 'application/octet-stream';
}

/**
 * Преобразует локальный путь в полный путь к файлу
 */
function getLocalFilePath(imagePath: string): string {
  // Убираем начальный слэш, если есть
  const cleanPath = imagePath.startsWith('/') ? imagePath.slice(1) : imagePath;
  return path.join(process.cwd(), cleanPath);
}

/**
 * Мигрирует изображения одного товара
 */
async function migrateProductImages(
  productId: number,
  images: string[],
  stats: MigrationStats,
): Promise<string[]> {
  const newImageUrls: string[] = [];

  for (const imagePath of images) {
    try {
      // Пропускаем, если это уже S3 URL
      if (imagePath.includes('s3.ru1.storage.beget.cloud')) {
        console.log(`  ⏭️  Skipping (already S3): ${imagePath}`);
        newImageUrls.push(imagePath);
        continue;
      }

      // Получаем локальный путь к файлу
      const localFilePath = getLocalFilePath(imagePath);

      // Проверяем существование файла
      if (!fs.existsSync(localFilePath)) {
        console.log(`  ⚠️  File not found: ${localFilePath}`);
        stats.errors.push({
          type: 'product',
          id: productId,
          path: imagePath,
          error: 'File not found',
        });
        stats.failedImages++;
        continue;
      }

      // Загружаем в S3
      console.log(`  ⬆️  Uploading: ${imagePath}`);
      const s3Url = await uploadFileToS3(localFilePath, 'products');
      newImageUrls.push(s3Url);
      stats.uploadedImages++;
      console.log(`  ✅ Uploaded to: ${s3Url}`);
    } catch (error) {
      console.error(`  ❌ Error uploading ${imagePath}:`, error.message);
      stats.errors.push({
        type: 'product',
        id: productId,
        path: imagePath,
        error: error.message,
      });
      stats.failedImages++;
    }
  }

  return newImageUrls;
}

/**
 * Мигрирует фото одного пользователя
 */
async function migrateUserPhoto(
  userId: number,
  photoPath: string,
  stats: MigrationStats,
): Promise<string | null> {
  try {
    // Пропускаем, если это уже S3 URL
    if (photoPath.includes('s3.ru1.storage.beget.cloud')) {
      console.log(`  ⏭️  Skipping (already S3): ${photoPath}`);
      return photoPath;
    }

    // Получаем локальный путь к файлу
    const localFilePath = getLocalFilePath(photoPath);

    // Проверяем существование файла
    if (!fs.existsSync(localFilePath)) {
      console.log(`  ⚠️  File not found: ${localFilePath}`);
      stats.errors.push({
        type: 'user',
        id: userId,
        path: photoPath,
        error: 'File not found',
      });
      stats.failedUserPhotos++;
      return null;
    }

    // Загружаем в S3
    console.log(`  ⬆️  Uploading: ${photoPath}`);
    const s3Url = await uploadFileToS3(localFilePath, 'users');
    stats.uploadedUserPhotos++;
    console.log(`  ✅ Uploaded to: ${s3Url}`);
    return s3Url;
  } catch (error) {
    console.error(`  ❌ Error uploading ${photoPath}:`, error.message);
    stats.errors.push({
      type: 'user',
      id: userId,
      path: photoPath,
      error: error.message,
    });
    stats.failedUserPhotos++;
    return null;
  }
}

/**
 * Основная функция миграции
 */
async function migrateToS3() {
  console.log('🚀 Starting migration to S3...\n');

  const stats: MigrationStats = {
    totalProducts: 0,
    totalImages: 0,
    uploadedImages: 0,
    failedImages: 0,
    updatedProducts: 0,
    totalUsers: 0,
    uploadedUserPhotos: 0,
    failedUserPhotos: 0,
    updatedUsers: 0,
    errors: [],
  };

  try {
    // ========== МИГРАЦИЯ ТОВАРОВ ==========
    console.log('📦 MIGRATING PRODUCTS\n');
    const products = await prisma.product.findMany({
      select: {
        id: true,
        name: true,
        images: true,
      },
    });

    stats.totalProducts = products.length;
    console.log(`📦 Found ${products.length} products\n`);

    // Мигрируем каждый товар
    for (const product of products) {
      if (!product.images || product.images.length === 0) {
        console.log(
          `⏭️  Product #${product.id} (${product.name}): No images\n`,
        );
        continue;
      }

      stats.totalImages += product.images.length;
      console.log(
        `📦 Product #${product.id} (${product.name}): ${product.images.length} images`,
      );

      // Мигрируем изображения
      const newImageUrls = await migrateProductImages(
        product.id,
        product.images,
        stats,
      );

      // Если есть новые URL, обновляем товар в БД
      if (newImageUrls.length > 0) {
        await prisma.product.update({
          where: { id: product.id },
          data: { images: newImageUrls },
        });
        stats.updatedProducts++;
        console.log(`  💾 Updated database with ${newImageUrls.length} URLs`);
      }

      console.log('');
    }

    // ========== МИГРАЦИЯ ПОЛЬЗОВАТЕЛЕЙ ==========
    console.log('\n👤 MIGRATING USERS\n');
    const users = await prisma.user.findMany({
      select: {
        id: true,
        fullName: true,
        photo: true,
      },
      where: {
        photo: {
          not: null,
        },
      },
    });

    stats.totalUsers = users.length;
    console.log(`👤 Found ${users.length} users with photos\n`);

    // Мигрируем каждого пользователя
    for (const user of users) {
      if (!user.photo) continue;

      console.log(`👤 User #${user.id} (${user.fullName})`);

      // Мигрируем фото
      const newPhotoUrl = await migrateUserPhoto(user.id, user.photo, stats);

      // Если есть новый URL, обновляем пользователя в БД
      if (newPhotoUrl) {
        await prisma.user.update({
          where: { id: user.id },
          data: { photo: newPhotoUrl },
        });
        stats.updatedUsers++;
        console.log(`  💾 Updated database with new photo URL`);
      }

      console.log('');
    }

    // Выводим статистику
    console.log('\n' + '='.repeat(60));
    console.log('📊 Migration Statistics:');
    console.log('='.repeat(60));
    console.log(`Products:`);
    console.log(`  Total Products:       ${stats.totalProducts}`);
    console.log(`  Total Images:         ${stats.totalImages}`);
    console.log(`  Uploaded Images:      ${stats.uploadedImages} ✅`);
    console.log(`  Failed Images:        ${stats.failedImages} ❌`);
    console.log(`  Updated Products:     ${stats.updatedProducts}`);
    console.log('');
    console.log(`Users:`);
    console.log(`  Total Users:          ${stats.totalUsers}`);
    console.log(`  Uploaded Photos:      ${stats.uploadedUserPhotos} ✅`);
    console.log(`  Failed Photos:        ${stats.failedUserPhotos} ❌`);
    console.log(`  Updated Users:        ${stats.updatedUsers}`);
    console.log('='.repeat(60));

    if (stats.errors.length > 0) {
      console.log('\n⚠️  Errors:');
      stats.errors.forEach((err) => {
        const prefix = err.type === 'product' ? 'Product' : 'User';
        console.log(`  ${prefix} #${err.id} - ${err.path}: ${err.error}`);
      });
    }

    console.log('\n✨ Migration completed!\n');
  } catch (error) {
    console.error('\n❌ Migration failed:', error);
    throw error;
  } finally {
    await prisma.$disconnect();
  }
}

// Запускаем миграцию
migrateToS3().catch((error) => {
  console.error('Fatal error:', error);
  process.exit(1);
});
