import { Module } from '@nestjs/common';
import { ChatService } from './chat.service';
import { ChatController } from './chat.controller';
import { ChatGateway } from './gateway';
import { AuthModule } from 'src/auth/auth.module';
import { PrismaModule } from 'src/prisma/prisma.module';
import { WsJwtAuthGuard } from 'src/auth/guards/ws-jwt-auth.guard';

@Module({
  imports: [AuthModule, PrismaModule],
  controllers: [ChatController],
  providers: [ChatService, ChatGateway, WsJwtAuthGuard],
  exports: [ChatService],
})
export class ChatModule {}
