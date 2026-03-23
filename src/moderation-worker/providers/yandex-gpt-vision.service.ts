import { Injectable, Logger } from '@nestjs/common';
import { ConfigService } from '@nestjs/config';

import axios from 'axios';

export interface VisionModerationResult {
  decision: 'APPROVED' | 'MANUAL' | 'DENIED';
  reason: string;
}

interface OpenAICompatResponse {
  choices: Array<{
    message: {
      content: string;
    };
  }>;
}

const VISION_PROMPT = `Это фото для объявления на маркетплейсе.

Проверь наличие ЛЮБОГО из следующих нарушений:
1. Оружие, боеприпасы, взрывчатка, ножи как основной товар
2. NSFW / откровенный контент / части тела
3. Насилие, кровь, шокирующие материалы
4. Скриншот стороннего сайта/приложения с контактами (телефон, email, ник)
5. Наркотики, алкоголь, табак

Если хотя бы одно нарушение есть — DENIED.
Если фото нечёткое, подозрительное или невозможно определить товар — MANUAL.
Если фото обычного товара (в т.ч. немедицинского бытового) без нарушений — APPROVED.

Ответь СТРОГО в JSON без markdown-обёртки:
{
  "decision": "APPROVED" | "MANUAL" | "DENIED",
  "reason": "Объяснение на русском что именно нарушено (пустая строка если APPROVED)"
}

DENIED — явное нарушение из списка выше
MANUAL — сомнительно, нужна ручная проверка
APPROVED — фото подходит`;

const DOWNLOAD_TIMEOUT_MS = 30_000;
const REQUEST_TIMEOUT_MS = 60_000;
const VISION_ENDPOINT = 'https://ai.api.cloud.yandex.net/v1/chat/completions';

@Injectable()
export class YandexGptVisionService {
  private readonly logger = new Logger(YandexGptVisionService.name);

  constructor(private readonly configService: ConfigService) {}

  private async downloadImageAsBase64(
    imageUrl: string,
  ): Promise<{ base64: string; mimeType: string }> {
    this.logger.log(`[Vision] Downloading image: ${imageUrl}`);
    const resp = await axios.get<ArrayBuffer>(imageUrl, {
      responseType: 'arraybuffer',
      timeout: DOWNLOAD_TIMEOUT_MS,
      headers: { 'User-Agent': 'MedVito-ModerationWorker/1.0' },
    });
    const base64 = Buffer.from(resp.data).toString('base64');
    const mimeType =
      (resp.headers['content-type'] as string | undefined) ?? 'image/jpeg';
    this.logger.log(
      `[Vision] Downloaded OK: ${base64.length} chars base64, mime=${mimeType}`,
    );
    return { base64, mimeType };
  }

  async moderateImage(imageUrl: string): Promise<VisionModerationResult> {
    const folderId = this.configService.getOrThrow<string>('YANDEX_FOLDER_ID');
    const apiKey = this.configService.getOrThrow<string>('YANDEX_API_KEY');
    const model = `gpt://${folderId}/gemma-3-27b-it/latest`;

    this.logger.log(`[Vision] START imageUrl=${imageUrl}`);
    this.logger.log(`[Vision] endpoint=${VISION_ENDPOINT} model=${model}`);

    let rawResponseData: unknown = null;

    try {
      const { base64, mimeType } = await this.downloadImageAsBase64(imageUrl);

      const requestBody = {
        model,
        max_tokens: 150,
        temperature: 0.05,
        messages: [
          {
            role: 'user',
            content: [
              {
                type: 'image_url',
                // Yandex requires base64 data URI — direct URLs are not supported
                image_url: { url: `data:${mimeType};base64,${base64}` },
              },
              { type: 'text', text: VISION_PROMPT },
            ],
          },
        ],
      };

      this.logger.debug(
        `[Vision] Request body: ${JSON.stringify(requestBody)}`,
      );

      const response = await axios.post(VISION_ENDPOINT, requestBody, {
        headers: {
          Authorization: `Api-Key ${apiKey}`,
          'Content-Type': 'application/json',
          'x-folder-id': folderId,
        },
        timeout: REQUEST_TIMEOUT_MS,
      });

      rawResponseData = response.data;
      this.logger.log(
        `[Vision] HTTP ${response.status}. Raw response: ${JSON.stringify(response.data)}`,
      );

      // Cast to typed response shape (OpenAI-compatible endpoint)
      const data = response.data as OpenAICompatResponse;
      const messageContent = data?.choices?.[0]?.message?.content;

      if (!messageContent) {
        this.logger.error(
          `[Vision] choices[0].message.content is missing. Full response: ${JSON.stringify(response.data)}`,
        );
        throw new Error(
          'choices[0].message.content is missing in API response',
        );
      }

      const raw = messageContent.trim();
      this.logger.log(`[Vision] Raw content from model: ${raw}`);

      const cleaned = raw.replace(/```json|```/g, '').trim();
      this.logger.log(`[Vision] Cleaned JSON string: ${cleaned}`);

      let parsed: VisionModerationResult;
      try {
        parsed = JSON.parse(cleaned) as VisionModerationResult;
      } catch (parseErr) {
        this.logger.error(
          `[Vision] JSON.parse failed: ${(parseErr as Error).message}. Input was: ${cleaned}`,
        );
        throw parseErr;
      }

      this.logger.log(
        `[Vision] Parsed result: decision=${parsed.decision} reason=${parsed.reason}`,
      );

      if (
        !['APPROVED', 'MANUAL', 'DENIED'].includes(parsed.decision) ||
        typeof parsed.reason !== 'string'
      ) {
        this.logger.error(
          `[Vision] Validation failed — unexpected shape: ${JSON.stringify(parsed)}`,
        );
        throw new Error('Unexpected vision response shape');
      }

      this.logger.log(`[Vision] SUCCESS → ${parsed.decision}`);
      return parsed;
    } catch (error) {
      const err = error as {
        code?: string;
        message?: string;
        isAxiosError?: boolean;
        config?: { url?: string; method?: string };
        response?: {
          status?: number;
          statusText?: string;
          data?: unknown;
          headers?: unknown;
        };
      };

      // Detailed Axios error dump
      if (err.isAxiosError) {
        this.logger.error(`[Vision] AXIOS ERROR:`);
        this.logger.error(`  message   : ${err.message ?? 'n/a'}`);
        this.logger.error(`  code      : ${err.code ?? 'n/a'}`);
        this.logger.error(`  url       : ${err.config?.url ?? 'n/a'}`);
        this.logger.error(
          `  http status: ${err.response?.status ?? 'no response'} ${err.response?.statusText ?? ''}`,
        );
        this.logger.error(
          `  resp body : ${JSON.stringify(err.response?.data ?? null)}`,
        );
        this.logger.error(
          `  resp hdrs : ${JSON.stringify(err.response?.headers ?? null)}`,
        );
      } else {
        this.logger.error(
          `[Vision] NON-AXIOS ERROR: ${err.message ?? String(error)}`,
        );
        if (rawResponseData !== null) {
          this.logger.error(
            `[Vision] Last raw API response was: ${JSON.stringify(rawResponseData)}`,
          );
        }
      }

      return {
        decision: 'MANUAL',
        reason: 'Ошибка анализа фото, требуется ручная проверка',
      };
    }
  }
}
