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

const SYSTEM_PROMPT = `Ты — автоматический модератор объявлений на маркетплейсе.


Запрещено:
- Контакты в тексте: телефон, email, ник в мессенджере (даже написанные словами: "восемь девятьсот", "собака", "тг @...")
- Внешние ссылки (http://, t.me/, wa.me/ и т.д.)
- Нецензурная лексика и оскорбления
- Мошеннические признаки: "предоплата", "переведи деньги", "только безнал", "аванс"
- SEO-спам: многократные повторы одних и тех же слов или ключевых фраз, бессмысленный набор текста, перечисление несвязанных слов для выдачи в поиске
- НЕ является спамом: вежливые фразы продавца ("советуем заглянуть", "в нашем профиле есть другие товары", "отличное качество"), стандартные описания товара, упоминание ассортимента магазина
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
    const folderId = this.configService.getOrThrow<string>('YANDEX_FOLDER_ID');
    const apiKey = this.configService.getOrThrow<string>('YANDEX_API_KEY');

    const userPrompt = `Категория: ${categoryName}
Подкатегория: ${subcategoryName}
Цена: ${price} руб.
Название: ${name}
Описание: ${description || 'не указано'}`;

    try {
      const response = await axios.post(
        this.endpoint,
        {
          modelUri: `gpt://${folderId}/yandexgpt/latest`,
          completionOptions: {
            stream: false,
            temperature: 0.05,
            maxTokens: 250,
          },
          messages: [
            { role: 'system', text: SYSTEM_PROMPT },
            { role: 'user', text: userPrompt },
          ],
        },
        {
          headers: {
            Authorization: `Api-Key ${apiKey}`,
            'Content-Type': 'application/json',
          },
          timeout: 25000,
        },
      );

      const raw = String(
        response.data.result.alternatives[0].message.text,
      ).trim();
      this.logger.log(`[Text] Raw response: ${raw}`);

      const cleaned = raw.replace(/```json|```/g, '').trim();

      // If the model returned a plain-text refusal instead of JSON, treat as DENIED
      if (!cleaned.startsWith('{')) {
        this.logger.warn(
          `[Text] Model returned non-JSON refusal: ${cleaned.substring(0, 200)}`,
        );
        return {
          category: 'DENIED',
          reason: `Модель отказалась обрабатывать объявление: ${cleaned.substring(0, 300)}`,
          details: {
            categorization: 'VIOLATION',
            spam: 'OK',
            fraud: 'OK',
            contacts: 'OK',
          },
        };
      }

      const parsed = JSON.parse(cleaned) as TextModerationResult;

      // Validate the expected shape to guard against prompt injection in the AI response
      if (
        !['APPROVED', 'MANUAL', 'DENIED'].includes(parsed.category) ||
        typeof parsed.reason !== 'string' ||
        !parsed.details
      ) {
        throw new Error(`Unexpected response shape from YandexGPT: ${cleaned}`);
      }

      return parsed;
    } catch (error) {
      this.logger.error(`Text moderation error: ${(error as Error).message}`);
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
