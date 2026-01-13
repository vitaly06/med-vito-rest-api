// @ts-check
import eslint from '@eslint/js';
import eslintPluginPrettierRecommended from 'eslint-plugin-prettier/recommended';
import globals from 'globals';
import tseslint from 'typescript-eslint';
import { defineConfig } from 'eslint/config';
import perfectionist from 'eslint-plugin-perfectionist';

export default defineConfig(
  {
    ignores: ['eslint.config.mjs'],
  },
  eslint.configs.recommended,
  ...tseslint.configs.recommendedTypeChecked,
  eslintPluginPrettierRecommended,
  {
    languageOptions: {
      globals: {
        ...globals.node,
        ...globals.jest,
      },
      sourceType: 'commonjs',
      parserOptions: {
        projectService: true,
        tsconfigRootDir: import.meta.dirname,
      },
    },
  },
  {
    rules: {
      '@typescript-eslint/no-explicit-any': 'off',
      '@typescript-eslint/no-floating-promises': 'warn',
      '@typescript-eslint/no-unsafe-argument': 'warn',
      'prettier/prettier': ['error', { endOfLine: 'auto' }],
    },
  },
  {
    plugins: { perfectionist },
    rules: {
      'perfectionist/sort-named-imports': [
        'error',
        { order: 'asc', type: 'natural' },
      ],
      'perfectionist/sort-imports': [
        'error',
        {
          type: 'natural',
          order: 'asc',
          newlinesBetween: 1,
          groups: [
            'nestjs',
            'type',
            ['builtin', 'external'],
            'type-internal',
            ['internal'],
            ['type-parent', 'type-sibling', 'type-index'],
            ['parent', 'sibling', 'index'],
            'side-effect-style',
            'unknown',
          ],
          customGroups: [
            {
              groupName: 'nestjs',
              elementNamePattern: ['^@nestjs/.*'],
            },
          ],
          internalPattern: ['^@/.*', '^~/.*'],
        },
      ],
    },
  },
);
