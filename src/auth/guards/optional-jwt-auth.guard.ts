import { Injectable, ExecutionContext, CanActivate } from '@nestjs/common';
import { JwtService } from '@nestjs/jwt';
import { PrismaService } from 'src/prisma/prisma.service';

@Injectable()
export class OptionalJwtAuthGuard implements CanActivate {
  constructor(
    private jwtService: JwtService,
    private prisma: PrismaService,
  ) {}

  async canActivate(context: ExecutionContext): Promise<boolean> {
    const request = context.switchToHttp().getRequest();

    try {
      // Пытаемся получить токен из куков (приоритет) или заголовков
      let token = null;

      // 1. Получаем из куков
      if (request?.cookies?.access_token) {
        token = request.cookies.access_token;
      }

      // 2. Fallback к заголовку Authorization
      if (!token && request.headers.authorization) {
        token = request.headers.authorization.replace('Bearer ', '');
      }

      // Если токена нет - это нормально, просто продолжаем без аутентификации
      if (!token) {
        request.user = null;
        return true;
      }

      // Проверяем токен
      const payload = this.jwtService.verify(token, {
        secret: process.env.JWT_SECRET || 'fallback-secret',
      });

      // Получаем пользователя
      const user = await this.prisma.user.findUnique({
        where: { id: payload.sub },
      });

      // Если пользователь найден - добавляем в request
      request.user = user || null;

      return true;
    } catch (error) {
      // Если произошла ошибка с токеном - просто игнорируем и продолжаем без аутентификации
      console.log('Optional JWT Auth failed:', error.message);
      request.user = null;
      return true;
    }
  }
}
