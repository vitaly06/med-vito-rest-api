import { Logger } from '@nestjs/common';
import { NestFactory } from '@nestjs/core';

import { ModerationWorkerModule } from './moderation-worker.module';

async function bootstrap() {
  const app = await NestFactory.createApplicationContext(
    ModerationWorkerModule,
    { logger: ['log', 'warn', 'error', 'debug'] },
  );

  await app.init();

  const logger = new Logger('ModerationWorker');
  logger.log('Worker process is running. Press Ctrl+C to stop.');

  const shutdown = async (signal: string) => {
    logger.log(`${signal} received — shutting down gracefully...`);
    await app.close();
    process.exit(0);
  };

  process.on('SIGTERM', () => void shutdown('SIGTERM'));
  process.on('SIGINT', () => void shutdown('SIGINT'));
}

bootstrap().catch((err: unknown) => {
  console.error('Failed to start moderation worker:', err);
  process.exit(1);
});
