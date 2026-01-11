import { Module } from '@nestjs/common';

import { AuthModule } from '@/auth/auth.module';
import { WsSessionAuthGuard } from '@/auth/guards/ws-session-auth.guard';
import { PrismaModule } from '@/prisma/prisma.module';

import { ChatController } from './chat.controller';
import { ChatService } from './chat.service';
import { ChatGateway } from './gateway';

@Module({
  imports: [AuthModule, PrismaModule],
  controllers: [ChatController],
  providers: [ChatService, ChatGateway, WsSessionAuthGuard],
  exports: [ChatService, ChatGateway],
})
export class ChatModule {}
