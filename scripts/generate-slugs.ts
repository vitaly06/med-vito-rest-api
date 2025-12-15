import { PrismaClient } from '@prisma/client';

const prisma = new PrismaClient();

function transliterate(text: string): string {
  const ru2en: Record<string, string> = {
    а: 'a',
    б: 'b',
    в: 'v',
    г: 'g',
    д: 'd',
    е: 'e',
    ё: 'e',
    ж: 'zh',
    з: 'z',
    и: 'i',
    й: 'y',
    к: 'k',
    л: 'l',
    м: 'm',
    н: 'n',
    о: 'o',
    п: 'p',
    р: 'r',
    с: 's',
    т: 't',
    у: 'u',
    ф: 'f',
    х: 'h',
    ц: 'ts',
    ч: 'ch',
    ш: 'sh',
    щ: 'sch',
    ъ: '',
    ы: 'y',
    ь: '',
    э: 'e',
    ю: 'yu',
    я: 'ya',
    А: 'A',
    Б: 'B',
    В: 'V',
    Г: 'G',
    Д: 'D',
    Е: 'E',
    Ё: 'E',
    Ж: 'Zh',
    З: 'Z',
    И: 'I',
    Й: 'Y',
    К: 'K',
    Л: 'L',
    М: 'M',
    Н: 'N',
    О: 'O',
    П: 'P',
    Р: 'R',
    С: 'S',
    Т: 'T',
    У: 'U',
    Ф: 'F',
    Х: 'H',
    Ц: 'Ts',
    Ч: 'Ch',
    Ш: 'Sh',
    Щ: 'Sch',
    Ъ: '',
    Ы: 'Y',
    Ь: '',
    Э: 'E',
    Ю: 'Yu',
    Я: 'Ya',
  };

  return text
    .split('')
    .map((char) => ru2en[char] || char)
    .join('');
}

function generateSlug(name: string): string {
  return transliterate(name)
    .toLowerCase()
    .replace(/[^a-z0-9]+/g, '-')
    .replace(/^-+|-+$/g, '');
}

async function generateSlugs() {
  console.log('🚀 Генерация slug для существующих данных...\n');

  // 1. Генерация slug для категорий
  console.log('📁 Обработка категорий...');
  const categories = await prisma.category.findMany({
    where: {
      slug: '',
    },
  });

  for (const category of categories) {
    let slug = generateSlug(category.name);
    let counter = 1;

    // Проверка уникальности
    while (
      await prisma.category.findFirst({
        where: { slug, id: { not: category.id } },
      })
    ) {
      counter++;
      slug = `${generateSlug(category.name)}-${counter}`;
    }

    await prisma.category.update({
      where: { id: category.id },
      data: { slug },
    });
    console.log(`  ✓ ${category.name} → ${slug}`);
  }

  // 2. Генерация slug для подкатегорий
  console.log('\n📂 Обработка подкатегорий...');
  const subcategories = await prisma.subCategory.findMany({
    where: {
      slug: '',
    },
  });

  for (const subcategory of subcategories) {
    let slug = generateSlug(subcategory.name);
    let counter = 1;

    // Проверка уникальности в рамках родительской категории
    while (
      await prisma.subCategory.findFirst({
        where: {
          slug,
          categoryId: subcategory.categoryId,
          id: { not: subcategory.id },
        },
      })
    ) {
      counter++;
      slug = `${generateSlug(subcategory.name)}-${counter}`;
    }

    await prisma.subCategory.update({
      where: { id: subcategory.id },
      data: { slug },
    });
    console.log(`  ✓ ${subcategory.name} → ${slug}`);
  }

  // 3. Генерация slug для типов подкатегорий
  console.log('\n📄 Обработка типов подкатегорий...');
  const types = await prisma.subcategotyType.findMany({
    where: {
      slug: '',
    },
  });

  for (const type of types) {
    let slug = generateSlug(type.name);
    let counter = 1;

    // Проверка уникальности в рамках родительской подкатегории
    while (
      await prisma.subcategotyType.findFirst({
        where: {
          slug,
          subcategoryId: type.subcategoryId,
          id: { not: type.id },
        },
      })
    ) {
      counter++;
      slug = `${generateSlug(type.name)}-${counter}`;
    }

    await prisma.subcategotyType.update({
      where: { id: type.id },
      data: { slug },
    });
    console.log(`  ✓ ${type.name} → ${slug}`);
  }

  console.log('\n✅ Генерация slug завершена!');
}

generateSlugs()
  .catch((e) => {
    console.error('❌ Ошибка:', e);
    process.exit(1);
  })
  .finally(async () => {
    await prisma.$disconnect();
  });
