import { PrismaClient } from '@prisma/client';

const prisma = new PrismaClient();

async function checkSlugs() {
  try {
    const categories = await prisma.category.findMany();
    const subCategories = await prisma.subCategory.findMany();
    const types = await prisma.subcategotyType.findMany();

    console.log('\n=== Categories ===');
    categories.forEach((c) => {
      console.log(`ID: ${c.id}, Name: ${c.name}, Slug: ${c.slug}`);
    });

    console.log('\n=== SubCategories ===');
    subCategories.forEach((s) => {
      console.log(
        `ID: ${s.id}, Name: ${s.name}, Slug: ${s.slug}, CategoryId: ${s.categoryId}`,
      );
    });

    console.log('\n=== Types ===');
    console.log(`Total types: ${types.length}`);
    const typesWithoutSlug = types.filter((t) => !t.slug || t.slug === '');
    console.log(`Types without slug: ${typesWithoutSlug.length}`);

    if (typesWithoutSlug.length > 0) {
      console.log('\nTypes without slug:');
      typesWithoutSlug.forEach((t) => {
        console.log(`ID: ${t.id}, Name: ${t.name}`);
      });
    }

    // Check for duplicates
    const categorySlugCounts = new Map<string, number>();
    categories.forEach((c) => {
      categorySlugCounts.set(c.slug, (categorySlugCounts.get(c.slug) || 0) + 1);
    });
    const duplicateCategorySlugs = Array.from(
      categorySlugCounts.entries(),
    ).filter(([_, count]) => count > 1);
    if (duplicateCategorySlugs.length > 0) {
      console.log('\n⚠️  Duplicate category slugs:', duplicateCategorySlugs);
    }

    const subCategorySlugCounts = new Map<string, number>();
    subCategories.forEach((s) => {
      const key = `${s.slug}-${s.categoryId}`;
      subCategorySlugCounts.set(key, (subCategorySlugCounts.get(key) || 0) + 1);
    });
    const duplicateSubCategorySlugs = Array.from(
      subCategorySlugCounts.entries(),
    ).filter(([_, count]) => count > 1);
    if (duplicateSubCategorySlugs.length > 0) {
      console.log(
        '\n⚠️  Duplicate subcategory slugs:',
        duplicateSubCategorySlugs,
      );
    }

    console.log('\n✅ All checks completed!');
  } catch (error) {
    console.error('Error:', error);
  } finally {
    await prisma.$disconnect();
  }
}

checkSlugs();
