import {
  Injectable,
  ExecutionContext,
  UnauthorizedException,
  CanActivate,
  NotFoundException,
  ForbiddenException,
} from '@nestjs/common';
import { JwtService } from '@nestjs/jwt';
import { PrismaService } from 'src/prisma/prisma.service';

@Injectable()
export class AdminJwtAuthGuard implements CanActivate {
  constructor(
    private jwtService: JwtService,
    private prisma: PrismaService,
  ) {}

  async canActivate(context: ExecutionContext): Promise<boolean> {
    const request = context.switchToHttp().getRequest();
    const accessToken = request.cookies?.access_token;

    if (!accessToken) {
      throw new UnauthorizedException('Access токен не найден');
    }

    try {
      // Проверяем access токен
      const payload = this.jwtService.verify(accessToken, {
        secret: process.env.JWT_SECRET || 'fallback-secret',
      });

      const user = await this.prisma.user.findUnique({
        where: { id: payload.sub },
      });

      if (!user) {
        throw new UnauthorizedException('Пользователь не найден');
      }

      const role = await this.prisma.role.findUnique({
        where: { name: 'admin' },
      });

      if (!role) {
        throw new NotFoundException('Роль admin не найдена');
      }

      if (role.id != user.roleId) {
        throw new ForbiddenException('Недостаточно прав');
      }

      if (user) request.user = user;
      return true;
    } catch (error) {
      // Если это HTTP исключение (UnauthorizedException, ForbiddenException и т.д.), пробрасываем его дальше
      if (
        error instanceof UnauthorizedException ||
        error instanceof ForbiddenException ||
        error instanceof NotFoundException
      ) {
        throw error;
      }

      // Для JWT ошибок
      if (error.name === 'TokenExpiredError') {
        throw new UnauthorizedException('Access токен истёк');
      }

      throw new UnauthorizedException('Недействительный токен');
    }
  }
}
