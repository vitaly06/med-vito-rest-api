import { PrismaClient } from '@prisma/client';

/**
 * Генерирует уникальный семизначный ID
 * @param prisma - экземпляр PrismaClient
 * @param modelName - имя модели ('user' или 'product')
 * @returns уникальный семизначный ID
 */
export async function generateUniqueId(
  prisma: PrismaClient,
  modelName: 'user' | 'product',
): Promise<number> {
  const MIN_ID = 1000000; // Минимальное семизначное число
  const MAX_ID = 9999999; // Максимальное семизначное число
  const MAX_ATTEMPTS = 100; // Максимальное количество попыток

  for (let attempt = 0; attempt < MAX_ATTEMPTS; attempt++) {
    // Генерируем случайное семизначное число
    const randomId = Math.floor(Math.random() * (MAX_ID - MIN_ID + 1)) + MIN_ID;

    // Проверяем уникальность ID
    const exists = await checkIdExists(prisma, modelName, randomId);

    if (!exists) {
      return randomId;
    }
  }

  throw new Error(
    `Не удалось сгенерировать уникальный ID для ${modelName} после ${MAX_ATTEMPTS} попыток`,
  );
}

/**
 * Проверяет существование ID в базе данных
 * @param prisma - экземпляр PrismaClient
 * @param modelName - имя модели ('user' или 'product')
 * @param id - ID для проверки
 * @returns true если ID существует, false если нет
 */
async function checkIdExists(
  prisma: PrismaClient,
  modelName: 'user' | 'product',
  id: number,
): Promise<boolean> {
  try {
    if (modelName === 'user') {
      const user = await prisma.user.findUnique({
        where: { id },
        select: { id: true },
      });
      return !!user;
    } else if (modelName === 'product') {
      const product = await prisma.product.findUnique({
        where: { id },
        select: { id: true },
      });
      return !!product;
    }
    return false;
  } catch (error) {
    // Если таблица не существует или другая ошибка, считаем что ID не существует
    return false;
  }
}
