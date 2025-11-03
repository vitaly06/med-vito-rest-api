import {
  Injectable,
  ExecutionContext,
  UnauthorizedException,
  CanActivate,
} from '@nestjs/common';
import { JwtService } from '@nestjs/jwt';
import { PrismaService } from 'src/prisma/prisma.service';
import { Socket } from 'socket.io';

@Injectable()
export class WsJwtAuthGuard implements CanActivate {
  constructor(
    private jwtService: JwtService,
    private prisma: PrismaService,
  ) {}

  async canActivate(context: ExecutionContext): Promise<boolean> {
    try {
      const client: Socket = context.switchToWs().getClient();

      // Получаем токен из куков (приоритет) или из других источников
      let token: string | null = null;

      // 1. Пытаемся получить из куков (как в основной стратегии)
      const cookies = client.handshake.headers.cookie;
      if (cookies) {
        const accessTokenMatch = cookies.match(/access_token=([^;]+)/);
        if (accessTokenMatch) {
          token = accessTokenMatch[1];
        }
      }

      // 2. Fallback к другим методам, если куки недоступны
      if (!token) {
        token =
          (client.handshake.auth?.token as string) ||
          client.handshake.headers?.authorization?.replace('Bearer ', '') ||
          (client.handshake.query?.token as string);
      }

      if (!token) {
        throw new UnauthorizedException('Токен не найден');
      }

      // Проверяем токен
      const payload = this.jwtService.verify(token, {
        secret: process.env.JWT_SECRET || 'fallback-secret',
      });

      // Получаем пользователя
      const user = await this.prisma.user.findUnique({
        where: { id: payload.sub },
      });

      if (!user) {
        throw new UnauthorizedException('Пользователь не найден');
      }

      // Сохраняем пользователя в client для дальнейшего использования
      (client as any).user = user;
      (client as any).userId = user.id;

      return true;
    } catch (error) {
      console.log('WebSocket Auth Error:', error.message);
      return false;
    }
  }
}
