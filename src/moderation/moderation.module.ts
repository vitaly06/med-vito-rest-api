import { Module } from '@nestjs/common';

import { AuthModule } from 'src/auth/auth.module';
import { PrismaModule } from 'src/prisma/prisma.module';

import { ModerationController } from './moderation.controller';

@Module({
  imports: [AuthModule, PrismaModule],
  controllers: [ModerationController],
})
export class ModerationModule {}
