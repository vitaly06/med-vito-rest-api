import { PrismaClient } from '@prisma/client';

const prisma = new PrismaClient();

async function fillProductTypes() {
  console.log('🚀 Заполнение typeId для товаров...\n');

  // Получаем все товары с typeId = null
  const productsWithoutType = await prisma.product.findMany({
    where: {
      typeId: null,
    },
    include: {
      fieldValues: {
        include: {
          field: true,
        },
      },
    },
  });

  console.log(`📦 Найдено товаров без типа: ${productsWithoutType.length}\n`);

  let updated = 0;
  let skipped = 0;

  for (const product of productsWithoutType) {
    if (!product.subCategoryId) {
      console.log(
        `⚠️  Товар ${product.id} (${product.name}) - нет subCategoryId, пропускаем`,
      );
      skipped++;
      continue;
    }

    // Получаем все типы для подкатегории товара
    const availableTypes = await prisma.subcategotyType.findMany({
      where: {
        subcategoryId: product.subCategoryId,
      },
      include: {
        fields: true,
      },
    });

    if (availableTypes.length === 0) {
      console.log(
        `⚠️  Товар ${product.id} (${product.name}) - нет типов для подкатегории ${product.subCategoryId}, пропускаем`,
      );
      skipped++;
      continue;
    }

    let selectedType: (typeof availableTypes)[0] | null = null;

    // Если у товара есть характеристики, ищем подходящий тип
    if (product.fieldValues.length > 0) {
      const productFieldIds = product.fieldValues.map((fv) => fv.field.typeId);

      // Ищем тип, у которого есть поля, совпадающие с характеристиками товара
      for (const type of availableTypes) {
        const typeFieldIds = type.fields.map((f) => f.typeId);

        // Проверяем, что все поля товара принадлежат этому типу
        const allFieldsMatch = productFieldIds.every((fieldTypeId) =>
          typeFieldIds.includes(fieldTypeId),
        );

        if (allFieldsMatch) {
          selectedType = type;
          break;
        }
      }

      // Если не нашли точное совпадение, берем тип, у которого хотя бы одно поле совпадает
      if (!selectedType) {
        for (const type of availableTypes) {
          const typeFieldIds = type.fields.map((f) => f.typeId);
          const hasAnyMatch = productFieldIds.some((fieldTypeId) =>
            typeFieldIds.includes(fieldTypeId),
          );

          if (hasAnyMatch) {
            selectedType = type;
            break;
          }
        }
      }
    }

    // Если не нашли подходящий тип или нет характеристик, берем первый доступный
    if (!selectedType) {
      selectedType = availableTypes[0];
    }

    // Обновляем товар
    await prisma.product.update({
      where: { id: product.id },
      data: { typeId: selectedType.id },
    });

    console.log(
      `✓ Товар ${product.id} (${product.name}) → Тип: ${selectedType.name} (ID: ${selectedType.id})`,
    );
    updated++;
  }

  console.log(`\n✅ Обработка завершена!`);
  console.log(`   Обновлено: ${updated}`);
  console.log(`   Пропущено: ${skipped}`);
}

fillProductTypes()
  .catch((e) => {
    console.error('❌ Ошибка:', e);
    process.exit(1);
  })
  .finally(async () => {
    await prisma.$disconnect();
  });
