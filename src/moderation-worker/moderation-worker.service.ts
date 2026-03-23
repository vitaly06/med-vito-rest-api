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

      // ── Step 2: Vision moderation (first image only) ─────────────────────
      let visionDecision: 'APPROVED' | 'MANUAL' | 'DENIED' = 'APPROVED';
      let visionReason = '';

      if (product.images.length > 0) {
        const visionResult = await this.gptVision.moderateImage(
          product.images[0],
        );
        visionDecision = visionResult.decision;
        visionReason = visionResult.reason;
        this.logger.log(`[#${id}] Vision → ${visionDecision}`);
      }

      if (visionDecision === 'DENIED') {
        await this.applyDecision(id, 'DENIDED', visionReason);
        return;
      }

      // ── Step 3: Final decision ────────────────────────────────────────────
      if (textResult.category === 'MANUAL' || visionDecision === 'MANUAL') {
        const reasons: string[] = [];
        if (textResult.category === 'MANUAL' && textResult.reason) {
          reasons.push(`Текст: ${textResult.reason}`);
        }
        const isVisionTechnicalFailure =
          visionDecision === 'MANUAL' &&
          visionReason.includes(VISION_TECHNICAL_ERROR_REASON);

        // Ignore transient technical failures of image analysis.
        if (
          visionDecision === 'MANUAL' &&
          visionReason &&
          !isVisionTechnicalFailure
        ) {
          reasons.push(`Фото: ${visionReason}`);
        }
        if (reasons.length === 0) {
          // If AI marked MANUAL but could not explain why, treat it as APPROVED.
          await this.applyDecision(id, 'APPROVED', AI_APPROVED_REASON);
          this.logger.log(`[#${id}] → APPROVED (${AI_APPROVED_REASON})`);
          return;
        }

        const manualReason = reasons.join(' / ');
        await this.applyDecision(id, 'AI_REVIEWED', manualReason);
        this.logger.log(`[#${id}] → MANUAL (${manualReason})`);
        return;
      }

      await this.applyDecision(id, 'APPROVED', AI_APPROVED_REASON);
      this.logger.log(`[#${id}] → APPROVED (${AI_APPROVED_REASON})`);
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
