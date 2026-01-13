import { Module } from '@nestjs/common';

import { AuthModule } from 'src/auth/auth.module';

import { PrismaModule } from '../prisma/prisma.module';
import { SupportController } from './support.controller';
import { SupportGateway } from './support.gateway';
import { SupportService } from './support.service';

@Module({
  imports: [PrismaModule, AuthModule],
  controllers: [SupportController],
  providers: [SupportService, SupportGateway],
  exports: [SupportService, SupportGateway],
})
export class SupportModule {}
