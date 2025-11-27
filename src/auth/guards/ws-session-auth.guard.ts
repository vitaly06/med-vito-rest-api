import {
  Injectable,
  ExecutionContext,
  UnauthorizedException,
  CanActivate,
} from '@nestjs/common';
import { PrismaService } from 'src/prisma/prisma.service';
import { Socket } from 'socket.io';
import { AuthService } from '../auth.service';

@Injectable()
export class WsSessionAuthGuard implements CanActivate {
  constructor(
    private authService: AuthService,
    private prisma: PrismaService,
  ) {}

  async canActivate(context: ExecutionContext): Promise<boolean> {
    try {
      const client: Socket = context.switchToWs().getClient();

      // Получаем session ID из куков
      let sessionId: string | null = null;

      const cookies = client.handshake.headers.cookie;
      if (cookies) {
        const sessionIdMatch = cookies.match(/session_id=([^;]+)/);
        if (sessionIdMatch) {
          sessionId = sessionIdMatch[1];
        }
      }

      if (!sessionId) {
        throw new UnauthorizedException('Сессия не найдена');
      }

      // Получаем данные сессии из Redis
      const sessionData = await this.authService.getSessionData(sessionId);

      if (!sessionData) {
        throw new UnauthorizedException('Сессия недействительна или истекла');
      }

      // Получаем пользователя из БД
      const user = await this.prisma.user.findUnique({
        where: { id: sessionData.userId },
        include: { role: true },
      });

      if (!user) {
        throw new UnauthorizedException('Пользователь не найден');
      }

      // Сохраняем пользователя в client для дальнейшего использования
      client.data.user = user;

      return true;
    } catch (error) {
      throw new UnauthorizedException('Ошибка аутентификации');
    }
  }
}
