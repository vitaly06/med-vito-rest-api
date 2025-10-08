import {
  BadRequestException,
  Injectable,
  Logger,
  UnauthorizedException,
} from '@nestjs/common';
import { JwtService } from '@nestjs/jwt';
import { PrismaService } from 'src/prisma/prisma.service';
import { SignUpDto } from './dto/sign-up.dto';
import * as bcrypt from 'bcrypt';
import { SignInDto } from './dto/sign-in.dto';
import { Response } from 'express';

@Injectable()
export class AuthService {
  private readonly logger = new Logger('Auth');
  constructor(
    private readonly prisma: PrismaService,
    private readonly jwtService: JwtService,
  ) {}

  async signUp(dto: SignUpDto) {
    const checkUser = await this.prisma.user.findFirst({
      where: {
        OR: [{ phoneNumber: dto.phoneNumber }, { email: dto.email }],
      },
    });

    if (checkUser) {
      throw new BadRequestException('Данный пользователь уже существует');
    }

    const hashedPassword = await bcrypt.hash(dto.password, 10);
    try {
      await this.prisma.user.create({
        data: {
          ...dto,
          password: hashedPassword,
        },
      });

      this.logger.log(`Успешная регистрация: ${dto.email}, ${dto.phoneNumber}`);

      return { message: 'Вы успешно зарегистрировались!' };
    } catch (e) {
      this.logger.error(`Ошибка: ${e}`);
    }
  }

  async signIn(dto: SignInDto, response: Response) {
    const { login, password } = { ...dto };
    // Проверка на почту / телефон
    const checkUser = await this.prisma.user.findFirst({
      where: {
        OR: [{ phoneNumber: login }, { email: login }],
      },
    });

    if (!checkUser) {
      throw new UnauthorizedException('Неверные данные для входа');
    }

    const passwordIsMatch = await bcrypt.compare(password, checkUser.password);

    if (!passwordIsMatch) {
      throw new UnauthorizedException('Неверные данные для входа');
    }

    // Генерируем токены
    const tokens = await this.generateTokens(checkUser.id, checkUser.email);

    // Сохраняем refresh токен в БД
    const refreshTokenExpiresAt = new Date();
    refreshTokenExpiresAt.setDate(refreshTokenExpiresAt.getDate() + 7);

    await this.prisma.user.update({
      where: { id: checkUser.id },
      data: {
        refreshToken: tokens.refreshToken,
        refreshTokenExpiresAt,
      },
    });

    this.setTokenCookies(response, tokens);

    this.logger.log(`Успешная авторизация: ${login}`);

    return {
      message: 'Вы успешно авторизовались!',
      user: {
        id: checkUser.id,
        email: checkUser.email,
        fullName: checkUser.fullName,
        phoneNumber: checkUser.phoneNumber,
        profileType: checkUser.profileType,
      },
    };
  }

  async logout(userId: number, response: Response) {
    // Очищаем refresh токен в БД
    await this.prisma.user.update({
      where: { id: userId },
      data: {
        refreshToken: null,
        refreshTokenExpiresAt: null,
      },
    });

    response.clearCookie('access_token');
    response.clearCookie('refresh_token');

    return { message: 'Вы успешно вышли из системы!' };
  }

  async refreshTokens(refreshToken: string, response: Response) {
    try {
      // Проверяем refresh токен
      const payload = this.jwtService.verify(refreshToken, {
        secret: process.env.JWT_REFRESH_SECRET || 'fallback-refresh-secret',
      });

      // Находим пользователя
      const user = await this.prisma.user.findFirst({
        where: {
          id: payload.sub,
          refreshToken: refreshToken,
          refreshTokenExpiresAt: {
            gt: new Date(),
          },
        },
      });

      if (!user) {
        throw new UnauthorizedException('Недействительный refresh токен');
      }

      // Генерируем новые токены
      const tokens = await this.generateTokens(user.id, user.email);

      // Обновляем refresh токен в БД
      const refreshTokenExpiresAt = new Date();
      refreshTokenExpiresAt.setDate(refreshTokenExpiresAt.getDate() + 7);

      await this.prisma.user.update({
        where: { id: user.id },
        data: {
          refreshToken: tokens.refreshToken,
          refreshTokenExpiresAt,
        },
      });

      // Устанавливаем новые токены в cookies
      this.setTokenCookies(response, tokens);

      return {
        message: 'Токены успешно обновлены!',
        user: {
          id: user.id,
          email: user.email,
          fullName: user.fullName,
          phoneNumber: user.phoneNumber,
          profileType: user.profileType,
        },
      };
    } catch (error) {
      throw new UnauthorizedException('Недействительный refresh токен');
    }
  }

  private async generateTokens(userId: number, email: string) {
    const payload = { sub: userId, email };

    const accessToken = this.jwtService.sign(payload, {
      secret: process.env.JWT_SECRET || 'fallback-secret',
      expiresIn: '15m',
    });

    const refreshToken = this.jwtService.sign(
      { sub: userId },
      {
        secret: process.env.JWT_REFRESH_SECRET || 'fallback-refresh-secret',
        expiresIn: '7d',
      },
    );

    return { accessToken, refreshToken };
  }

  private setTokenCookies(
    response: Response,
    tokens: { accessToken: string; refreshToken: string },
  ) {
    response.cookie('access_token', tokens.accessToken, {
      httpOnly: true,
      secure: process.env.NODE_ENV === 'production',
      sameSite: 'strict',
      maxAge: 15 * 60 * 1000, // 15 минут
    });

    response.cookie('refresh_token', tokens.refreshToken, {
      httpOnly: true,
      secure: process.env.NODE_ENV === 'production',
      sameSite: 'strict',
      maxAge: 7 * 24 * 60 * 60 * 1000, // 7 дней
    });
  }
}
