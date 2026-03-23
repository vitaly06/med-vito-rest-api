import {
  Injectable,
  Logger,
  OnModuleDestroy,
  OnModuleInit,
} from '@nestjs/common';

import { PrismaService } from 'src/prisma/prisma.service';

import {
  TextModerationResult,
  YandexGptTextService,
} from './providers/yandex-gpt-text.service';
import { YandexGptVisionService } from './providers/yandex-gpt-vision.service';

/** How often the worker polls the database for new products to moderate. */
const POLL_INTERVAL_MS = 30_000;

/** Max products to pull from DB in a single poll cycle. */
const BATCH_SIZE = 5;

/** Marker for products approved automatically by AI. */
const AI_APPROVED_REASON = 'Одобрено ИИ автоматически';
const VISION_TECHNICAL_ERROR_REASON =
  'Ошибка анализа фото, требуется ручная проверка';

interface ProductWithRelations {
  id: number;
  name: string;
  description: string | null;
  price: number;
  images: string[];
  category: { name: string };
  subCategory: { name: string };
}

@Injectable()
export class ModerationWorkerService implements OnModuleInit, OnModuleDestroy {
  private readonly logger = new Logger(ModerationWorkerService.name);
  private pollTimer: NodeJS.Timeout | null = null;
  private isProcessing = false;

  constructor(
    private readonly prisma: PrismaService,
    private readonly gptText: YandexGptTextService,
    private readonly gptVision: YandexGptVisionService,
  ) {}

  onModuleInit() {
    this.logger.log(
      `Moderation worker started. Poll interval: ${POLL_INTERVAL_MS / 1000}s`,
    );
    // Run first poll immediately, then schedule recurring polls.
    void this.runPoll().then(() => this.scheduleNextPoll());
  }

  onModuleDestroy() {
    if (this.pollTimer) {
      clearTimeout(this.pollTimer);
      this.pollTimer = null;
    }
  }

  private scheduleNextPoll() {
    this.pollTimer = setTimeout(() => {
      void this.runPoll().then(() => this.scheduleNextPoll());
    }, POLL_INTERVAL_MS);
  }

  private async runPoll(): Promise<void> {
    if (this.isProcessing) {
      this.logger.warn('Previous poll still running, skipping this tick');
      return;
    }
    this.isProcessing = true;

    try {
      const products = await this.prisma.product.findMany({
        where: {
          moderateState: 'MODERATE',
          isHide: false,
          // Only pick up products that haven't been through AI yet
          moderationRejectionReason: null,
        },
        orderBy: { createdAt: 'asc' },
        take: BATCH_SIZE,
        include: {
          category: { select: { name: true } },
          subCategory: { select: { name: true } },
        },
      });

      if (products.length === 0) {
        this.logger.debug('No products pending moderation');
        return;
      }

      this.logger.log(`Found ${products.length} product(s) pending moderation`);

      for (const product of products) {
        await this.processProduct(product);
      }
    } catch (err) {
      this.logger.error(`Poll cycle error: ${(err as Error).message}`);
    } finally {
      this.isProcessing = false;
    }
  }

  private async processProduct(product: ProductWithRelations): Promise<void> {
    const { id } = product;
    this.logger.log(`[#${id}] Processing "${product.name}"`);

    try {
      // ── Step 1: Text moderation ──────────────────────────────────────────
      const textResult: TextModerationResult = await this.gptText.moderateText(
        product.name,
        product.description,
        product.category.name,
        product.subCategory.name,
        product.price,
      );

      this.logger.log(
        `[#${id}] Text → ${textResult.category}` +
          (textResult.reason ? ` (${textResult.reason})` : ''),
      );

      if (textResult.category === 'DENIED') {
        await this.applyDecision(id, 'DENIDED', textResult.reason);
        return;
      }

      // ── Step 2: Vision moderation (all images) ───────────────────────────
      let visionDecision: 'APPROVED' | 'MANUAL' | 'DENIED' = 'APPROVED';
      let visionReason = '';
      let visionTechnicalFailure = false;

      for (const imageUrl of product.images) {
        const visionResult = await this.gptVision.moderateImage(imageUrl);
        this.logger.log(
          `[#${id}] Vision [${imageUrl}] → ${visionResult.decision}`,
        );

        // Technical failure from the vision service
        if (
          visionResult.decision === 'MANUAL' &&
          visionResult.reason === VISION_TECHNICAL_ERROR_REASON
        ) {
          visionTechnicalFailure = true;
          continue; // check remaining images
        }

        // Escalate to the worst outcome seen across all images
        if (visionResult.decision === 'DENIED') {
          visionDecision = 'DENIED';
          visionReason = visionResult.reason;
          break;
        }
        if (visionResult.decision === 'MANUAL') {
          visionDecision = 'MANUAL';
          visionReason = visionResult.reason;
        }
      }

      if (visionDecision === 'DENIED') {
        await this.applyDecision(id, 'DENIDED', visionReason);
        return;
      }

      // ── Step 3: Final decision ────────────────────────────────────────────
      const reasons: string[] = [];

      if (textResult.category === 'MANUAL' && textResult.reason) {
        reasons.push(`Текст: ${textResult.reason}`);
      }
      if (visionDecision === 'MANUAL' && visionReason) {
        reasons.push(`Фото: ${visionReason}`);
      }
      // Vision failed technically — cannot confirm photos are safe, require human
      if (visionTechnicalFailure && visionDecision !== 'MANUAL') {
        reasons.push(VISION_TECHNICAL_ERROR_REASON);
      }

      if (reasons.length > 0) {
        const manualReason = reasons.join(' / ');
        await this.applyDecision(id, 'AI_REVIEWED', manualReason);
        this.logger.log(`[#${id}] → MANUAL (${manualReason})`);
        return;
      }

      if (textResult.category === 'APPROVED' && visionDecision === 'APPROVED') {

        await this.applyDecision(id, 'APPROVED', AI_APPROVED_REASON);
        this.logger.log(`[#${id}] → APPROVED (${AI_APPROVED_REASON})`);
        return;
      }

      // Fallback — should not normally reach here
      await this.applyDecision(id, 'AI_REVIEWED', 'Требуется ручная проверка');
      this.logger.log(`[#${id}] → MANUAL (fallback)`);
    } catch (err) {
      // Log and continue — don't let one failure stop the whole batch
      this.logger.error(`[#${id}] Unexpected error: ${(err as Error).message}`);
    }
  }

  private async applyDecision(
    productId: number,
    state: 'APPROVED' | 'DENIDED' | 'AI_REVIEWED',
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
