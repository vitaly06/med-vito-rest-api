import {
  CanActivate,
  ExecutionContext,
  ForbiddenException,
  Injectable,
  UnauthorizedException,
} from '@nestjs/common';
import { Request } from 'express';
import { AuthService } from '../auth.service';
import { PrismaService } from 'src/prisma/prisma.service';

@Injectable()
export class AdminSessionAuthGuard implements CanActivate {
  constructor(
    private authService: AuthService,
    private prisma: PrismaService,
  ) {}

  async canActivate(context: ExecutionContext): Promise<boolean> {
    const request = context.switchToHttp().getRequest<Request>();
    const sessionId = request.cookies?.session_id;

    if (!sessionId) {
      throw new UnauthorizedException('Сессия не найдена');
    }

    const sessionData = await this.authService.getSessionData(sessionId);

    if (!sessionData) {
      throw new UnauthorizedException('Сессия недействительна или истекла');
    }

    // Получаем полные данные пользователя из БД
    const user = await this.prisma.user.findUnique({
      where: { id: sessionData.userId },
      include: {
        role: true,
      },
    });

    if (!user) {
      throw new UnauthorizedException('Пользователь не найден');
    }

    // Проверяем, что у пользователя роль admin
    if (!user.role || user.role.name !== 'admin') {
      throw new ForbiddenException('Доступ запрещён');
    }

    // Прикрепляем пользователя к запросу
    request['user'] = user;

    return true;
  }
}
