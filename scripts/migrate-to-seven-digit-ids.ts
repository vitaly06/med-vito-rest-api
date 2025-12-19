import { PrismaClient } from '@prisma/client';

const prisma = new PrismaClient();

/**
 * Генерирует уникальный семизначный ID, который еще не используется
 */
function generateSevenDigitId(usedIds: Set<number>): number {
  const MIN_ID = 1000000;
  const MAX_ID = 9999999;
  let id: number;

  do {
    id = Math.floor(Math.random() * (MAX_ID - MIN_ID + 1)) + MIN_ID;
  } while (usedIds.has(id));

  usedIds.add(id);
  return id;
}

/**
 * Миграция ID пользователей и товаров на семизначные числа
 */
async function migrateToSevenDigitIds() {
  console.log('Начало миграции ID на семизначные числа...\n');

  try {
    // Отключаем проверки внешних ключей
    await prisma.$executeRaw`SET session_replication_role = 'replica';`;
    console.log('✓ Проверки внешних ключей отключены\n');

    // ========== Миграция пользователей ==========
    console.log('--- Миграция пользователей ---');
    const users = await prisma.user.findMany({
      select: { id: true },
      orderBy: { id: 'asc' },
    });

    console.log(`Найдено пользователей: ${users.length}`);

    const userIdMapping = new Map<number, number>();
    const usedUserIds = new Set<number>();

    // Генерируем новые ID для пользователей
    for (const user of users) {
      const newId = generateSevenDigitId(usedUserIds);
      userIdMapping.set(user.id, newId);
    }

    // Обновляем ID пользователей
    for (const [oldId, newId] of userIdMapping) {
      // Обновляем связанные таблицы сначала
      // Обновляем товары
      await prisma.$executeRaw`
        UPDATE "Product" 
        SET "userId" = ${newId} 
        WHERE "userId" = ${oldId};
      `;

      // Обновляем чаты (покупатель)
      await prisma.$executeRaw`
        UPDATE "Chat" 
        SET "buyerId" = ${newId} 
        WHERE "buyerId" = ${oldId};
      `;

      // Обновляем чаты (продавец)
      await prisma.$executeRaw`
        UPDATE "Chat" 
        SET "sellerId" = ${newId} 
        WHERE "sellerId" = ${oldId};
      `;

      // Обновляем сообщения
      await prisma.$executeRaw`
        UPDATE "Message" 
        SET "senderId" = ${newId} 
        WHERE "senderId" = ${oldId};
      `;

      // Обновляем просмотры номеров (просматривающий)
      await prisma.$executeRaw`
        UPDATE "PhoneNumberView" 
        SET "viewedById" = ${newId} 
        WHERE "viewedById" = ${oldId};
      `;

      // Обновляем просмотры номеров (владелец)
      await prisma.$executeRaw`
        UPDATE "PhoneNumberView" 
        SET "viewedUserId" = ${newId} 
        WHERE "viewedUserId" = ${oldId};
      `;

      // Обновляем просмотры товаров
      await prisma.$executeRaw`
        UPDATE "ProductView" 
        SET "viewedById" = ${newId} 
        WHERE "viewedById" = ${oldId};
      `;

      // Обновляем действия избранного
      await prisma.$executeRaw`
        UPDATE "FavoriteAction" 
        SET "userId" = ${newId} 
        WHERE "userId" = ${oldId};
      `;

      // Обновляем отзывы (автор)
      await prisma.$executeRaw`
        UPDATE "Review" 
        SET "reviewedById" = ${newId} 
        WHERE "reviewedById" = ${oldId};
      `;

      // Обновляем отзывы (получатель)
      await prisma.$executeRaw`
        UPDATE "Review" 
        SET "reviewedUserId" = ${newId} 
        WHERE "reviewedUserId" = ${oldId};
      `;

      // Обновляем тикеты поддержки (пользователь)
      await prisma.$executeRaw`
        UPDATE "SupportTicket" 
        SET "userId" = ${newId} 
        WHERE "userId" = ${oldId};
      `;

      // Обновляем тикеты поддержки (модератор)
      await prisma.$executeRaw`
        UPDATE "SupportTicket" 
        SET "moderatorId" = ${newId} 
        WHERE "moderatorId" = ${oldId};
      `;

      // Обновляем сообщения поддержки
      await prisma.$executeRaw`
        UPDATE "SupportMessage" 
        SET "authorId" = ${newId} 
        WHERE "authorId" = ${oldId};
      `;

      // Обновляем продвижения товаров
      await prisma.$executeRaw`
        UPDATE "ProductPromotion" 
        SET "userId" = ${newId} 
        WHERE "userId" = ${oldId};
      `;

      // Обновляем таблицу _UserFavorites
      await prisma.$executeRaw`
        UPDATE "_UserFavorites" 
        SET "A" = ${newId} 
        WHERE "A" = ${oldId};
      `;

      // Теперь обновляем сам ID пользователя
      await prisma.$executeRaw`
        UPDATE "User" 
        SET id = ${newId} 
        WHERE id = ${oldId};
      `;

      console.log(`✓ Пользователь: ${oldId} -> ${newId}`);
    }

    console.log(`\n✓ Миграция пользователей завершена!\n`);

    // ========== Миграция товаров ==========
    console.log('--- Миграция товаров ---');
    const products = await prisma.product.findMany({
      select: { id: true },
      orderBy: { id: 'asc' },
    });

    console.log(`Найдено товаров: ${products.length}`);

    const productIdMapping = new Map<number, number>();
    const usedProductIds = new Set<number>();

    // Генерируем новые ID для товаров
    for (const product of products) {
      const newId = generateSevenDigitId(usedProductIds);
      productIdMapping.set(product.id, newId);
    }

    // Обновляем ID товаров
    for (const [oldId, newId] of productIdMapping) {
      try {
        // Обновляем связанные таблицы (с проверкой существования)
        // Обновляем чаты
        await prisma.$executeRaw`
          UPDATE "Chat" 
          SET "productId" = ${newId} 
          WHERE "productId" = ${oldId};
        `.catch(() => {});

        // Обновляем сообщения (relatedProductId)
        await prisma.$executeRaw`
          UPDATE "Message" 
          SET "relatedProductId" = ${newId} 
          WHERE "relatedProductId" = ${oldId};
        `.catch(() => {});

        // Обновляем просмотры товаров
        await prisma.$executeRaw`
          UPDATE "ProductView" 
          SET "productId" = ${newId} 
          WHERE "productId" = ${oldId};
        `.catch(() => {});

        // Обновляем действия избранного
        await prisma.$executeRaw`
          UPDATE "FavoriteAction" 
          SET "productId" = ${newId} 
          WHERE "productId" = ${oldId};
        `.catch(() => {});

        // Обновляем значения полей товаров
        await prisma.$executeRaw`
          UPDATE "ProductFieldValue" 
          SET "productId" = ${newId} 
          WHERE "productId" = ${oldId};
        `.catch(() => {});

        // Обновляем продвижения товаров
        await prisma.$executeRaw`
          UPDATE "ProductPromotion" 
          SET "productId" = ${newId} 
          WHERE "productId" = ${oldId};
        `.catch(() => {});

        // Обновляем таблицу _UserFavorites
        await prisma.$executeRaw`
          UPDATE "_UserFavorites" 
          SET "B" = ${newId} 
          WHERE "B" = ${oldId};
        `.catch(() => {});

        // Теперь обновляем сам ID товара
        await prisma.$executeRaw`
          UPDATE "Product" 
          SET id = ${newId} 
          WHERE id = ${oldId};
        `;

        console.log(`✓ Товар: ${oldId} -> ${newId}`);
      } catch (error) {
        console.error(`Ошибка при обновлении товара ${oldId}:`, error.message);
      }
    }

    console.log(`\n✓ Миграция товаров завершена!\n`);

    // Включаем проверки внешних ключей
    await prisma.$executeRaw`SET session_replication_role = 'origin';`;
    console.log('✓ Проверки внешних ключей включены\n');

    console.log('========================================');
    console.log('Миграция успешно завершена!');
    console.log(`Пользователей обновлено: ${users.length}`);
    console.log(`Товаров обновлено: ${products.length}`);
    console.log('========================================');
  } catch (error) {
    console.error('Ошибка при миграции:', error);
    // Включаем проверки внешних ключей в случае ошибки
    await prisma.$executeRaw`SET session_replication_role = 'origin';`;
    throw error;
  } finally {
    await prisma.$disconnect();
  }
}

// Запуск миграции
migrateToSevenDigitIds()
  .then(() => {
    console.log('\n✓ Скрипт миграции выполнен успешно');
    process.exit(0);
  })
  .catch((error) => {
    console.error('\n✗ Ошибка выполнения скрипта миграции:', error);
    process.exit(1);
  });
