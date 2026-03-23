import { Module } from '@nestjs/common';
import { ConfigModule } from '@nestjs/config';

import { PrismaModule } from 'src/prisma/prisma.module';

import { ModerationWorkerService } from './moderation-worker.service';
import { YandexGptTextService } from './providers/yandex-gpt-text.service';
import { YandexGptVisionService } from './providers/yandex-gpt-vision.service';

@Module({
  imports: [
    ConfigModule.forRoot({ isGlobal: true }),
    PrismaModule,
  ],
  providers: [
    ModerationWorkerService,
    YandexGptTextService,
    YandexGptVisionService,
  ],
})
export class ModerationWorkerModule {}
