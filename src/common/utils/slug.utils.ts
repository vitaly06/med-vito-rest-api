/**
 * Утилиты для работы со slug
 */

/**
 * Таблица транслитерации кириллицы в латиницу
 */
const translitMap: Record<string, string> = {
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
};

/**
 * Генерирует slug из строки
 * @param text - Исходный текст
 * @returns slug в формате kebab-case
 */
export function generateSlug(text: string): string {
  return text
    .toLowerCase()
    .trim()
    .split('')
    .map((char) => translitMap[char] || char)
    .join('')
    .replace(/[^a-z0-9]+/g, '-') // Заменяем все что не буквы/цифры на дефис
    .replace(/^-+|-+$/g, '') // Убираем дефисы в начале и конце
    .replace(/-+/g, '-'); // Убираем множественные дефисы
}

/**
 * Генерирует уникальный slug, добавляя суффикс при конфликте
 * @param baseSlug - Базовый slug
 * @param existingSlugs - Массив существующих slug'ов
 * @returns Уникальный slug
 */
export function makeUniqueSlug(
  baseSlug: string,
  existingSlugs: string[],
): string {
  let slug = baseSlug;
  let counter = 2;

  while (existingSlugs.includes(slug)) {
    slug = `${baseSlug}-${counter}`;
    counter++;
  }

  return slug;
}

/**
 * Валидирует slug
 * @param slug - slug для проверки
 * @returns true если slug валидный
 */
export function isValidSlug(slug: string): boolean {
  return /^[a-z0-9]+(?:-[a-z0-9]+)*$/.test(slug);
}
