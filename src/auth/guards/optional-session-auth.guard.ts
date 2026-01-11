import { CanActivate, ExecutionContext, Injectable } from '@nestjs/common';

import { Request } from 'express';

import { PrismaService } from '@/prisma/prisma.service';

import { AuthService } from '../auth.service';

@Injectable()
export class OptionalSessionAuthGuard implements CanActivate {
  constructor(
    private authService: AuthService,
    private prisma: PrismaService,
  ) {}

  async canActivate(context: ExecutionContext): Promise<boolean> {
    const request = context.switchToHttp().getRequest<Request>();
    const sessionId = request.cookies?.session_id;

    if (!sessionId) {
      // Если нет сессии, пропускаем, но не прикрепляем пользователя
      return true;
    }

    const sessionData = await this.authService.getSessionData(sessionId);

    if (!sessionData) {
      // Если сессия недействительна, пропускаем
      return true;
    }

    // Получаем полные данные пользователя из БД
    const user = await this.prisma.user.findUnique({
      where: { id: sessionData.userId },
      include: {
        role: true,
      },
    });

    if (user) {
      // Прикрепляем пользователя к запросу
      request['user'] = user;
    }

    return true;
  }
}
