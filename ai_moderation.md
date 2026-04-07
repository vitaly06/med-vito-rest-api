# ИИ-модерация объявлений

## Обзор

Модерация объявлений выполняется отдельным процессом — **moderation-worker** (`src/moderation-worker/`), который запускается как самостоятельный Docker-контейнер. Он работает в фоновом режиме и не зависит от основного backend-а.

Используются две модели Yandex Cloud:
- **YandexGPT** (`yandexgpt/latest`) — анализ текста объявления
- **Gemma 3 27B** (`gemma-3-27b-it/latest`) — анализ фотографий

---

## Жизненный цикл объявления

```
Пользователь создаёт объявление
        │
        ▼
moderateState = MODERATE  (ожидает проверки)
        │
        ▼
ИИ-воркер подхватывает (каждые 30 сек)
        │
   ┌────┴────────────┐
   │                 │
DENIED        APPROVED / AI_REVIEWED
```

### Статусы в базе данных (`ProductModerate`)

| Статус | Значение |
|---|---|
| `MODERATE` | Ожидает проверки ИИ |
| `APPROVED` | Одобрено автоматически ИИ |
| `AI_REVIEWED` | ИИ направил на ручную проверку |
| `DENIDED` | Автоматически отклонено |

---

## Алгоритм обработки

Воркер каждые **30 секунд** делает выборку из БД: берёт до **5 товаров** со статусом `MODERATE`, у которых `isHide = false` и `moderationRejectionReason = null` (не проходили через ИИ), отсортированных по дате создания (старые первые).

### Шаг 1 — Текстовая модерация

Отправляет название, описание, категорию, подкатегорию и цену в **YandexGPT** (`llm.api.cloud.yandex.net/foundationModels/v1/completion`).

Проверяются 4 критерия:

| Критерий | Что проверяется |
|---|---|
| `categorization` | Соответствие товара категории, подозрительно низкая цена |
| `spam` | Повторы ключевых слов, бессмысленный текст |
| `fraud` | Признаки мошенничества: «предоплата», «переведи деньги» |
| `contacts` | Телефоны, email, ники мессенджеров в тексте, внешние ссылки |

Каждый критерий получает оценку: `OK` / `SUSPICIOUS` / `VIOLATION`.

**Результат текстовой проверки:**
- `DENIED` (хотя бы один `VIOLATION`) → объявление **сразу отклоняется**, фото не проверяются
- `MANUAL` (хотя бы один `SUSPICIOUS`) → помечается для ручной проверки, фото всё равно анализируются
- `APPROVED` → переходим к анализу фото

---

### Шаг 2 — Визуальная модерация

Проверяются **все фотографии** объявления по одной. Каждое фото:

1. **Скачивается** с S3 (таймаут 30 сек)
2. **Кодируется в base64** и отправляется в **Gemma 3 27B** (`ai.api.cloud.yandex.net/v1/chat/completions`) — Yandex требует именно base64 `data:image/...`, прямые URL не поддерживаются

Что проверяется на фото:
- Оружие, боеприпасы, взрывчатка, ножи
- NSFW / откровенный контент
- Насилие, кровь, шокирующие материалы
- Скриншоты сторонних сайтов с чужими контактами
- Наркотики, алкоголь, табак

**Итоговый вердикт по всем фото** — берётся наихудший результат:
- Если хотя бы одно фото `DENIED` → объявление **отклоняется**, остальные фото не проверяются
- Если хотя бы одно фото `MANUAL` → направляется на ручную проверку
- Если vision API недоступен технически → тоже направляется на ручную проверку (не одобряется автоматически)

---

### Шаг 3 — Финальное решение

| Текст | Фото | Итог |
|---|---|---|
| `APPROVED` | `APPROVED` | ✅ `APPROVED` — «Одобрено ИИ автоматически» |
| `MANUAL` | любое | 🔍 `AI_REVIEWED` — причина из текста (+ из фото если есть) |
| любое | `MANUAL` | 🔍 `AI_REVIEWED` — причина из фото |
| `DENIED` | — | ❌ `DENIDED` — причина из текста |
| любое | `DENIED` | ❌ `DENIDED` — причина из фото |
| техническая ошибка vision | — | 🔍 `AI_REVIEWED` — «Ошибка анализа фото, требуется ручная проверка» |

---

## Структура файлов

```
src/moderation-worker/
├── main.worker.ts                   # Точка входа отдельного процесса
├── moderation-worker.module.ts      # NestJS модуль
├── moderation-worker.service.ts     # Главная логика: поллинг, алгоритм решения
└── providers/
    ├── yandex-gpt-text.service.ts   # Текстовая модерация (YandexGPT)
    └── yandex-gpt-vision.service.ts # Визуальная модерация (Gemma 3 27B)
```

---

## API-эндпоинты Yandex Cloud

| Endpoint | Протокол | Используется для |
|---|---|---|
| `llm.api.cloud.yandex.net/foundationModels/v1/completion` | Нативный Yandex | Текст (YandexGPT) |
| `ai.api.cloud.yandex.net/v1/chat/completions` | OpenAI-совместимый | Фото (Gemma 3 27B) |

Аутентификация: заголовок `Authorization: Api-Key <YANDEX_API_KEY>`.  
Для vision endpoint дополнительно: заголовок `x-folder-id: <YANDEX_FOLDER_ID>`.

---

## Переменные окружения

| Переменная | Описание |
|---|---|
| `YANDEX_API_KEY` | API-ключ Yandex Cloud |
| `YANDEX_FOLDER_ID` | ID папки (folder) в Yandex Cloud |

---

## Важные особенности

- **Воркер не блокирует создание товара** — объявление сразу сохраняется, модерация асинхронная
- **Batch size = 5** — за один цикл обрабатывается не более 5 товаров, чтобы не перегружать API
- **Ошибка ≠ одобрение** — при любой технической ошибке vision товар уходит на ручную проверку, а не одобряется автоматически
- **Повторная обработка исключена** — условие `moderationRejectionReason = null` гарантирует, что уже проверенные товары не берутся снова
- **Фото в base64** — Yandex vision API не принимает прямые URL изображений, поэтому каждое фото скачивается и кодируется перед отправкой

---

## Содержание

1. [Архитектура и флоу](#1-архитектура-и-флоу)
2. [Структура файлов](#2-структура-файлов)
3. [Этап 1 — Установка зависимостей](#3-этап-1--установка-зависимостей)
4. [Этап 2 — Переменные окружения](#4-этап-2--переменные-окружения)
5. [Этап 3 — YandexGPT сервис (текст)](#5-этап-3--yandexgpt-сервис-текст)
6. [Этап 4 — YandexGPT сервис (фото)](#6-этап-4--yandexgpt-сервис-фото)
7. [Этап 5 — Orchestrator сервис](#7-этап-5--orchestrator-сервис)
8. [Этап 6 — BullMQ процессор](#8-этап-6--bullmq-процессор)
9. [Этап 7 — Moderation Module](#9-этап-7--moderation-module)
10. [Этап 8 — AppModule](#10-этап-8--appmodule)
11. [Этап 9 — Интеграция в ProductService](#11-этап-9--интеграция-в-productservice)
12. [Этап 10 — Admin контроллер](#12-этап-10--admin-контроллер)
13. [Промпты](#13-промпты)
14. [Тестирование](#14-тестирование)
15. [Стоимость](#15-стоимость)

---

## 1. Архитектура и флоу

### Логика маршрутизации

Каждый критерий получает статус: 🟢 ОК / 🟡 Сомнительно / 🔴 Нарушение

| Результат | Действие |
|---|---|
| Все критерии 🟢 | Автоматически публикуется |
| Хотя бы один 🟡 (нет 🔴) | Уходит в очередь модератора |
| Хотя бы один 🔴 | Автоматически отклоняется + копия модератору |

### ИИ-критерии

| Критерий | Модель | Описание |
|---|---|---|
| Категоризация | GPT 5.1 (текст) | Соответствует ли объявление выбранной категории |
| Спам | GPT 5.1 (текст) | SEO-спам, повторы, бессмысленный текст |
| Мошенничество | GPT 5.1 (текст) | Низкая цена, схемы обмана, предоплата |
| Контакты в тексте | GPT 5.1 (текст) | Скрытые телефоны/email/мессенджеры |
| Фото-анализ | GPT Pro multimodal | NSFW, нерелевантные фото |

### Асинхронный флоу

```
POST /product/create
        │
        ▼
ProductService.createProduct()
        │
        ├─ Сохраняет в БД: moderateState = 'MODERATE'
        ├─ Ставит job в BullMQ → возвращает 201 пользователю
        │
        ▼ (в фоне, через Redis)
ModerationProcessor.handleProduct()
        │
        ├─ GPT 5.1: анализ текста (название + описание)
        │     └─ Если 🔴 → сразу DENIDED, фото не проверяем
        │
        ├─ GPT Pro: анализ первого фото (если текст 🟢/🟡)
        │
        └─ Итоговое решение:
              APPROVED → product.moderateState = 'APPROVED'
              MANUAL   → product.moderateState = 'MODERATE' (ждёт модератора)
              DENIED   → product.moderateState = 'DENIDED' + reason
```

---

## 2. Структура файлов

Создать следующие файлы (все новые):

```
src/moderation/
├── moderation.module.ts
├── moderation.service.ts        ← ставит задачи в очередь
├── moderation.processor.ts      ← обрабатывает задачи из очереди
├── moderation.controller.ts     ← admin endpoints
├── moderation.constants.ts      ← константы очереди
└── providers/
    ├── yandex-gpt-text.service.ts   ← анализ текста через GPT 5.1
    └── yandex-gpt-vision.service.ts ← анализ фото через GPT Pro
```

---

## 3. Этап 1 — Установка зависимостей

```bash
npm install @nestjs/bull bull
npm install --save-dev @types/bull
npm install axios
```

> Redis уже подключён в проекте (`cache-manager-ioredis`, host: `redis`, port: `6379`). BullMQ будет использовать тот же Redis — ничего дополнительно поднимать не нужно.

---

## 4. Этап 2 — Переменные окружения

Добавить в `.env`:

```env
YANDEX_FOLDER_ID=b1g...         # ID папки в Yandex Cloud (из консоли)
YANDEX_API_KEY=AQVN...          # API ключ сервисного аккаунта (уже есть)
```

Добавить в `docker-compose.yml` в секцию `environment` сервиса API:

```yaml
- YANDEX_FOLDER_ID=${YANDEX_FOLDER_ID}
- YANDEX_API_KEY=${YANDEX_API_KEY}
```

> **Где взять FOLDER_ID:** зайти в [console.yandex.cloud](https://console.yandex.cloud) → выбрать папку → ID виден в URL: `console.yandex.cloud/folders/b1gXXXXXXXXXXXXXX`

---

## 5. Этап 3 — YandexGPT сервис (текст)

**Файл:** `src/moderation/providers/yandex-gpt-text.service.ts`

```typescript
import { Injectable, Logger } from '@nestjs/common';
import { ConfigService } from '@nestjs/config';
import axios from 'axios';

export interface TextModerationResult {
  category: 'APPROVED' | 'MANUAL' | 'DENIED';
  reason: string;
  details: {
    categorization: 'OK' | 'SUSPICIOUS' | 'VIOLATION';
    spam: 'OK' | 'SUSPICIOUS' | 'VIOLATION';
    fraud: 'OK' | 'SUSPICIOUS' | 'VIOLATION';
    contacts: 'OK' | 'SUSPICIOUS' | 'VIOLATION';
  };
}

@Injectable()
export class YandexGptTextService {
  private readonly logger = new Logger(YandexGptTextService.name);
  private readonly endpoint =
    'https://llm.api.cloud.yandex.net/foundationModels/v1/completion';

  constructor(private readonly configService: ConfigService) {}

  async moderateText(
    name: string,
    description: string | null,
    categoryName: string,
    subcategoryName: string,
    price: number,
  ): Promise<TextModerationResult> {
    const folderId = this.configService.get<string>('YANDEX_FOLDER_ID');
    const apiKey = this.configService.get<string>('YANDEX_API_KEY');

    const systemPrompt = `Ты — автоматический модератор объявлений на маркетплейсе медицинского оборудования MedVito.

Разрешённый контент: медицинское оборудование, инструменты, мебель, расходные материалы, ортопедия, стоматология, реабилитация, запчасти к медтехнике.

Запрещено:
- Контакты в тексте: телефон, email, ник в мессенджере (даже написанные словами: "восемь девятьсот", "собака", "тг @...")
- Внешние ссылки (http://, t.me/, wa.me/ и т.д.)
- Нецензурная лексика и оскорбления
- Товары не из медицинской сферы
- Мошеннические признаки: "предоплата", "переведи деньги", "только безнал", "аванс"
- SEO-спам: многократные повторы слов, бессмысленный набор текста
- Подозрительно низкая цена (менее 10% от рыночной для категории)

Для каждого из 4 критериев укажи статус: OK / SUSPICIOUS / VIOLATION.
- OK — нарушений нет
- SUSPICIOUS — есть признаки, но неоднозначно (требуется модератор)
- VIOLATION — явное нарушение (автоотказ)

Итоговое решение:
- APPROVED — все критерии OK
- MANUAL — хотя бы один SUSPICIOUS, нет VIOLATION
- DENIED — хотя бы один VIOLATION

Отвечай СТРОГО в JSON без markdown-обёртки:
{
  "category": "APPROVED" | "MANUAL" | "DENIED",
  "reason": "Объяснение на русском (пустая строка если APPROVED)",
  "details": {
    "categorization": "OK" | "SUSPICIOUS" | "VIOLATION",
    "spam": "OK" | "SUSPICIOUS" | "VIOLATION",
    "fraud": "OK" | "SUSPICIOUS" | "VIOLATION",
    "contacts": "OK" | "SUSPICIOUS" | "VIOLATION"
  }
}`;

    const userPrompt = `Категория: ${categoryName}
Подкатегория: ${subcategoryName}
Цена: ${price} руб.
Название: ${name}
Описание: ${description || 'не указано'}`;

    try {
      const response = await axios.post(
        this.endpoint,
        {
          modelUri: `gpt://${folderId}/yandexgpt/latest`, // GPT 5.1
          completionOptions: {
            stream: false,
            temperature: 0.05, // минимальная вариативность для классификации
            maxTokens: 200,
          },
          messages: [
            { role: 'system', text: systemPrompt },
            { role: 'user', text: userPrompt },
          ],
        },
        {
          headers: {
            Authorization: `Api-Key ${apiKey}`,
            'Content-Type': 'application/json',
          },
          timeout: 20000,
        },
      );

      const raw = response.data.result.alternatives[0].message.text.trim();
      const cleaned = raw.replace(/```json|```/g, '').trim();
      return JSON.parse(cleaned) as TextModerationResult;
    } catch (error) {
      this.logger.error(`Text moderation error: ${error?.message}`);
      // При ошибке API — отправляем на ручную проверку, не блокируем
      return {
        category: 'MANUAL',
        reason: 'Ошибка ИИ-сервиса, требуется ручная проверка',
        details: {
          categorization: 'OK',
          spam: 'OK',
          fraud: 'OK',
          contacts: 'OK',
        },
      };
    }
  }
}
```

---

## 6. Этап 4 — YandexGPT сервис (фото)

**Файл:** `src/moderation/providers/yandex-gpt-vision.service.ts`

> Фото скачивается с S3 и передаётся в API как base64. YandexGPT Pro принимает изображения в одном запросе с текстом.

```typescript
import { Injectable, Logger } from '@nestjs/common';
import { ConfigService } from '@nestjs/config';
import axios from 'axios';

export interface VisionModerationResult {
  decision: 'APPROVED' | 'MANUAL' | 'DENIED';
  reason: string;
}

@Injectable()
export class YandexGptVisionService {
  private readonly logger = new Logger(YandexGptVisionService.name);
  private readonly endpoint =
    'https://llm.api.cloud.yandex.net/foundationModels/v1/completion';

  constructor(private readonly configService: ConfigService) {}

  async moderateImage(imageUrl: string): Promise<VisionModerationResult> {
    const folderId = this.configService.get<string>('YANDEX_FOLDER_ID');
    const apiKey = this.configService.get<string>('YANDEX_API_KEY');

    try {
      // Скачиваем изображение с S3
      const imageResponse = await axios.get(imageUrl, {
        responseType: 'arraybuffer',
        timeout: 15000,
      });
      const base64 = Buffer.from(imageResponse.data).toString('base64');
      const mimeType = imageResponse.headers['content-type'] || 'image/jpeg';

      const response = await axios.post(
        this.endpoint,
        {
          modelUri: `gpt://${folderId}/yandexgpt-pro/latest`, // Pro — поддерживает изображения
          completionOptions: {
            stream: false,
            temperature: 0.05,
            maxTokens: 150,
          },
          messages: [
            {
              role: 'user',
              content: [
                {
                  type: 'image_url',
                  image_url: {
                    url: `data:${mimeType};base64,${base64}`,
                  },
                },
                {
                  type: 'text',
                  text: `Это фото для объявления на маркетплейсе медицинского оборудования.

Проверь:
1. Есть ли NSFW/откровенный контент, насилие, шокирующие материалы?
2. Это реальное фото товара, а не скриншот сайта/приложения?

Ответь СТРОГО в JSON без markdown-обёртки:
{
  "decision": "APPROVED" | "MANUAL" | "DENIED",
  "reason": "Объяснение на русском (пустая строка если APPROVED)"
}

APPROVED — фото подходит
MANUAL — сомнительно, нужна ручная проверка
DENIED — явное нарушение (NSFW, насилие)`,
                },
              ],
            },
          ],
        },
        {
          headers: {
            Authorization: `Api-Key ${apiKey}`,
            'Content-Type': 'application/json',
          },
          timeout: 30000,
        },
      );

      const raw = response.data.result.alternatives[0].message.text.trim();
      const cleaned = raw.replace(/```json|```/g, '').trim();
      return JSON.parse(cleaned) as VisionModerationResult;
    } catch (error) {
      this.logger.error(`Vision moderation error: ${error?.message}`);
      // При ошибке — на ручную проверку
      return {
        decision: 'MANUAL',
        reason: 'Ошибка анализа фото, требуется ручная проверка',
      };
    }
  }
}
```

---

## 7. Этап 5 — Orchestrator сервис

**Файл:** `src/moderation/moderation.service.ts`

```typescript
import { InjectQueue } from '@nestjs/bull';
import { Injectable, Logger } from '@nestjs/common';
import { Queue } from 'bull';

import { MODERATION_QUEUE } from './moderation.constants';

export interface ModerationJobData {
  productId: number;
}

@Injectable()
export class ModerationService {
  private readonly logger = new Logger(ModerationService.name);

  constructor(
    @InjectQueue(MODERATION_QUEUE) private readonly queue: Queue,
  ) {}

  async enqueueProductModeration(productId: number): Promise<void> {
    await this.queue.add(
      { productId } satisfies ModerationJobData,
      {
        attempts: 3,
        backoff: { type: 'exponential', delay: 5000 }, // retry: 5s → 10s → 20s
        removeOnComplete: true,
        removeOnFail: false, // оставлять failed jobs для отладки
      },
    );
    this.logger.log(`Queued moderation for product #${productId}`);
  }
}
```

**Файл:** `src/moderation/moderation.constants.ts`

```typescript
export const MODERATION_QUEUE = 'moderation';
```

---

## 8. Этап 6 — BullMQ процессор

**Файл:** `src/moderation/moderation.processor.ts`

```typescript
import { Process, Processor } from '@nestjs/bull';
import { Logger } from '@nestjs/common';
import { Job } from 'bull';

import { PrismaService } from 'src/prisma/prisma.service';

import { MODERATION_QUEUE } from './moderation.constants';
import { ModerationJobData } from './moderation.service';
import { YandexGptTextService } from './providers/yandex-gpt-text.service';
import { YandexGptVisionService } from './providers/yandex-gpt-vision.service';

@Processor(MODERATION_QUEUE)
export class ModerationProcessor {
  private readonly logger = new Logger(ModerationProcessor.name);

  constructor(
    private readonly prisma: PrismaService,
    private readonly gptText: YandexGptTextService,
    private readonly gptVision: YandexGptVisionService,
  ) {}

  @Process()
  async handle(job: Job<ModerationJobData>): Promise<void> {
    const { productId } = job.data;
    this.logger.log(`Processing moderation for product #${productId}`);

    const product = await this.prisma.product.findUnique({
      where: { id: productId },
      include: {
        category: { select: { name: true } },
        subCategory: { select: { name: true } },
      },
    });

    if (!product) {
      this.logger.warn(`Product #${productId} not found, skipping`);
      return;
    }

    // --- Шаг 1: Анализ текста ---
    const textResult = await this.gptText.moderateText(
      product.name,
      product.description,
      product.category.name,
      product.subCategory.name,
      product.price,
    );

    this.logger.log(
      `Product #${productId} text result: ${textResult.category}`,
    );

    // Если текст — явное нарушение, сразу отклоняем, фото не проверяем
    if (textResult.category === 'DENIED') {
      await this.updateProduct(productId, 'DENIDED', textResult.reason);
      return;
    }

    // --- Шаг 2: Анализ первого фото (если есть) ---
    let visionDecision: 'APPROVED' | 'MANUAL' | 'DENIED' = 'APPROVED';
    let visionReason = '';

    if (product.images && product.images.length > 0) {
      const visionResult = await this.gptVision.moderateImage(product.images[0]);
      visionDecision = visionResult.decision;
      visionReason = visionResult.reason;
      this.logger.log(
        `Product #${productId} vision result: ${visionDecision}`,
      );
    }

    // --- Шаг 3: Финальное решение ---
    // Если фото DENIED — отклоняем
    if (visionDecision === 'DENIED') {
      await this.updateProduct(productId, 'DENIDED', visionReason);
      return;
    }

    // Если текст или фото MANUAL — уходит модератору (остаётся MODERATE)
    if (textResult.category === 'MANUAL' || visionDecision === 'MANUAL') {
      this.logger.log(`Product #${productId} → MANUAL review`);
      // moderateState уже MODERATE, просто логируем
      return;
    }

    // Всё ОК — публикуем
    await this.updateProduct(productId, 'APPROVED', null);
    this.logger.log(`Product #${productId} → APPROVED`);
  }

  private async updateProduct(
    productId: number,
    state: 'APPROVED' | 'DENIDED',
    reason: string | null,
  ): Promise<void> {
    await this.prisma.product.update({
      where: { id: productId },
      data: {
        moderateState: state,
        moderationRejectionReason: reason,
      },
    });
  }
}
```

---

## 9. Этап 7 — Moderation Module

**Файл:** `src/moderation/moderation.module.ts`

```typescript
import { BullModule } from '@nestjs/bull';
import { Module } from '@nestjs/common';

import { PrismaModule } from 'src/prisma/prisma.module';

import { MODERATION_QUEUE } from './moderation.constants';
import { ModerationController } from './moderation.controller';
import { ModerationProcessor } from './moderation.processor';
import { ModerationService } from './moderation.service';
import { YandexGptTextService } from './providers/yandex-gpt-text.service';
import { YandexGptVisionService } from './providers/yandex-gpt-vision.service';

@Module({
  imports: [
    PrismaModule,
    BullModule.registerQueue({
      name: MODERATION_QUEUE,
      limiter: {
        max: 5,       // не более 5 запросов в секунду к Yandex API
        duration: 1000,
      },
    }),
  ],
  controllers: [ModerationController],
  providers: [
    ModerationService,
    ModerationProcessor,
    YandexGptTextService,
    YandexGptVisionService,
  ],
  exports: [ModerationService],
})
export class ModerationModule {}
```

---

## 10. Этап 8 — AppModule

**Файл:** `src/app.module.ts` — добавить два импорта:

```typescript
// 1. В список импортов добавить BullModule.forRoot (один раз, глобально):
import { BullModule } from '@nestjs/bull';
import { ModerationModule } from './moderation/moderation.module';

// 2. В массив imports:
BullModule.forRoot({
  redis: {
    host: 'redis',   // имя сервиса Redis из docker-compose
    port: 6379,
  },
}),
ModerationModule,
```

**Итоговый список imports в AppModule после добавления:**

```typescript
imports: [
  ConfigModule.forRoot({ isGlobal: true }),
  PrismaModule,
  CacheModule.register({ ... }),  // уже есть
  BullModule.forRoot({            // ← добавить
    redis: { host: 'redis', port: 6379 },
  }),
  ModerationModule,               // ← добавить
  // ... остальные модули без изменений
]
```

---

## 11. Этап 9 — Интеграция в ProductService

**Файл:** `src/product/product.service.ts`

### 11.1 Добавить импорт и зависимость в конструктор

```typescript
// Добавить импорт
import { ModerationService } from 'src/moderation/moderation.service';

// Добавить в конструктор
constructor(
  private readonly prisma: PrismaService,
  private readonly configService: ConfigService,
  private readonly userService: UserService,
  private readonly s3Service: S3Service,
  private readonly chatService: ChatService,
  private readonly chatGateway: ChatGateway,
  private readonly moderationService: ModerationService,  // ← добавить
) { ... }
```

### 11.2 Добавить вызов после создания продукта

Найти в `createProduct()` строку `return { message: 'Продукт успешно создан', ... }` и добавить *перед* ней:

```typescript
// Ставим в очередь на модерацию — не await, не блокируем ответ пользователю
this.moderationService
  .enqueueProductModeration(product.id)
  .catch((err) =>
    console.error('[Moderation] Failed to enqueue:', err),
  );

return {
  message: 'Продукт успешно создан',
  product: { ...product, images: product.images },
};
```

### 11.3 Добавить ModerationModule в ProductModule

**Файл:** `src/product/product.module.ts`

```typescript
import { ModerationModule } from 'src/moderation/moderation.module';

@Module({
  imports: [
    MulterModule.register({ storage: memoryStorage() }),
    AuthModule,
    PrismaModule,
    UserModule,
    S3Module,
    ChatModule,
    ModerationModule,  // ← добавить
  ],
  controllers: [ProductController],
  providers: [ProductService, OptionalSessionAuthGuard],
})
export class ProductModule {}
```

---

## 12. Этап 10 — Admin контроллер

**Файл:** `src/moderation/moderation.controller.ts`

Ручное переопределение решения ИИ — для случаев когда модератор не согласен с автоматическим решением.

```typescript
import {
  Body,
  Controller,
  Get,
  Param,
  Patch,
  Query,
  UseGuards,
} from '@nestjs/common';
import {
  ApiBearerAuth,
  ApiOperation,
  ApiQuery,
  ApiTags,
} from '@nestjs/swagger';
import { IsEnum, IsOptional, IsString } from 'class-validator';

import { AdminSessionAuthGuard } from 'src/auth/guards/admin-session-auth.guard';
import { PrismaService } from 'src/prisma/prisma.service';

class ManualDecisionDto {
  @IsEnum(['APPROVED', 'DENIDED'])
  decision: 'APPROVED' | 'DENIDED';

  @IsOptional()
  @IsString()
  reason?: string;
}

@ApiTags('Moderation (Admin)')
@ApiBearerAuth()
@UseGuards(AdminSessionAuthGuard)
@Controller('admin/moderation')
export class ModerationController {
  constructor(private readonly prisma: PrismaService) {}

  // Список объявлений ожидающих модерации
  @Get('queue')
  @ApiOperation({ summary: 'Очередь объявлений на ручную проверку' })
  @ApiQuery({ name: 'page', required: false })
  async getQueue(@Query('page') page = '1') {
    const take = 20;
    const skip = (parseInt(page) - 1) * take;

    const [items, total] = await Promise.all([
      this.prisma.product.findMany({
        where: { moderateState: 'MODERATE' },
        skip,
        take,
        orderBy: { createdAt: 'asc' }, // старые — первыми
        include: {
          user: { select: { id: true, fullName: true, email: true } },
          category: { select: { name: true } },
          subCategory: { select: { name: true } },
        },
      }),
      this.prisma.product.count({ where: { moderateState: 'MODERATE' } }),
    ]);

    return { items, total, page: parseInt(page), pages: Math.ceil(total / take) };
  }

  // Ручное решение по конкретному объявлению
  @Patch('product/:id')
  @ApiOperation({ summary: 'Ручная модерация объявления' })
  async moderateProduct(
    @Param('id') id: string,
    @Body() dto: ManualDecisionDto,
  ) {
    return this.prisma.product.update({
      where: { id: +id },
      data: {
        moderateState: dto.decision,
        moderationRejectionReason:
          dto.decision === 'DENIDED' ? (dto.reason ?? 'Отклонено модератором') : null,
      },
    });
  }
}
```

---

## 13. Промпты

### Советы по настройке промптов

**`temperature: 0.05`** — почти детерминированный ответ. Для классификации важна стабильность.

**Формат ответа:** всегда просить JSON без markdown. GPT иногда оборачивает в `` ```json `` — поэтому в коде есть `.replace(/```json|```/g, '')`.

**Контекст домена в промпте** — обязательно. Без него модель не знает что «медицинский тонометр» — нормальный товар, а «таблетки от всего» — нет.

### Пример хорошего ответа GPT 5.1

```json
{
  "category": "DENIED",
  "reason": "В описании обнаружены контактные данные: номер телефона записан словами",
  "details": {
    "categorization": "OK",
    "spam": "OK",
    "fraud": "OK",
    "contacts": "VIOLATION"
  }
}
```

### ModelUri для YandexGPT 5.1

```
gpt://{FOLDER_ID}/yandexgpt/latest
```

> Модель `yandexgpt` в Yandex Cloud соответствует актуальной версии GPT (сейчас 5.1). Суффикс `/latest` автоматически использует последнюю стабильную версию.

### ModelUri для YandexGPT Pro (с поддержкой изображений)

```
gpt://{FOLDER_ID}/yandexgpt-pro/latest
```

---

## 14. Тестирование

### 14.1 Проверить подключение к API вручную (curl)

```bash
curl -X POST https://llm.api.cloud.yandex.net/foundationModels/v1/completion \
  -H "Authorization: Api-Key AQVN..." \
  -H "Content-Type: application/json" \
  -d '{
    "modelUri": "gpt://b1g.../yandexgpt/latest",
    "completionOptions": { "stream": false, "temperature": 0.05, "maxTokens": 100 },
    "messages": [
      { "role": "user", "text": "Скажи \"OK\" и больше ничего." }
    ]
  }'
```

Ожидаемый ответ: `{ "result": { "alternatives": [{ "message": { "text": "OK" } }] } }`

### 14.2 Тестовые сценарии

| Сценарий | Ожидаемый результат |
|---|---|
| Обычное объявление «Тонометр Omron» | `APPROVED` |
| «Продам iPhone» в категории «Медоборудование» | `DENIED`, categorization: VIOLATION |
| «Звоните восемь девятьсот...» в описании | `DENIED`, contacts: VIOLATION |
| «Отдам даром» (цена 1 руб., подозрительно) | `MANUAL`, fraud: SUSPICIOUS |
| Размытое/нечёткое фото | `APPROVED` (не нарушение) |
| NSFW-фото | `DENIED` |

### 14.3 Проверить очередь BullMQ

После создания объявления в логах должно появиться:

```
[ModerationProcessor] Processing moderation for product #1234567
[ModerationProcessor] Product #1234567 text result: APPROVED
[ModerationProcessor] Product #1234567 vision result: APPROVED
[ModerationProcessor] Product #1234567 → APPROVED
```

### 14.4 Bull Board (опционально — UI для мониторинга очереди)

```bash
npm install @bull-board/nestjs @bull-board/api @bull-board/express
```

После настройки доступен на `/admin/queues` — показывает активные, завершённые и упавшие задачи.

---

## 15. Стоимость

### Тарифы Yandex Cloud

| Сервис | Цена за 1 000 токенов |
|---|---|
| YandexGPT (Lite/5.1) — вход | 1.20 ₽ |
| YandexGPT (Lite/5.1) — выход | 1.20 ₽ |
| YandexGPT Pro (vision) — вход | ~6.00 ₽ |
| YandexGPT Pro (vision) — выход | ~6.00 ₽ |

> Токен ≈ 4 символа. Уточнённые цены: [cloud.yandex.ru/prices](https://cloud.yandex.ru/prices)

### Стоимость одного объявления

| Операция | Токенов | Стоимость |
|---|---|---|
| Системный промпт + текст объявления (GPT 5.1) | ~800 вход + ~80 выход | ~1.06 ₽ |
| Анализ фото (GPT Pro, ~800 токенов фото) | ~900 вход + ~50 выход | ~5.70 ₽ |
| **Итого (текст + фото)** | | **~6.76 ₽** |
| **Только текст (нет фото или текст отклонён)** | | **~1.06 ₽** |

> **Оптимизация:** фото проверяется только при наличии фотографий И только если текст прошёл. Около 20% объявлений отклоняются на тексте → не тратим на них деньги за фото.

### Расчёт по объёму (реальная средняя ~4 ₽/объявление с фото)

| Объявлений/мес | ИИ-модерация | Штатный модератор |
|---|---|---|
| 500 | **~2 000 ₽** | ~45 000 ₽ |
| 2 000 | **~8 000 ₽** | ~45 000 ₽ |
| 5 000 | **~20 000 ₽** | ~60 000 ₽ |
| 10 000 | **~40 000 ₽** | ~90 000 ₽ |

---

## Чеклист перед запуском

- [ ] `YANDEX_FOLDER_ID` добавлен в `.env`
- [ ] `YANDEX_API_KEY` добавлен в `.env`
- [ ] Оба ключа добавлены в `docker-compose.yml`
- [ ] `npm install @nestjs/bull bull axios` выполнен
- [ ] Все 7 файлов модуля созданы
- [ ] `BullModule.forRoot` и `ModerationModule` добавлены в `AppModule`
- [ ] `ModerationModule` добавлен в `ProductModule`
- [ ] `ModerationService` инжектирован в `ProductService`
- [ ] Вызов `enqueueProductModeration` добавлен после `prisma.product.create`
- [ ] Проверен curl-запрос к API — ответ приходит
- [ ] Создано тестовое объявление — в логах видна работа процессора
