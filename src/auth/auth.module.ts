import { Module } from '@nestjs/common';
import { AuthService } from './auth.service';
import { AuthController } from './auth.controller';
import { SessionAuthGuard } from './guards/session-auth.guard';
import { AdminSessionAuthGuard } from './guards/admin-session-auth.guard';
import { OptionalSessionAuthGuard } from './guards/optional-session-auth.guard';
import { WsSessionAuthGuard } from './guards/ws-session-auth.guard';

@Module({
  imports: [],
  controllers: [AuthController],
  providers: [
    AuthService,
    SessionAuthGuard,
    AdminSessionAuthGuard,
    OptionalSessionAuthGuard,
    WsSessionAuthGuard,
  ],
  exports: [
    AuthService,
    SessionAuthGuard,
    AdminSessionAuthGuard,
    OptionalSessionAuthGuard,
    WsSessionAuthGuard,
  ],
})
export class AuthModule {}
