import { CACHE_MANAGER } from '@nestjs/cache-manager';
import {
  BadRequestException,
  ForbiddenException,
  Inject,
  Injectable,
  Logger,
  NotFoundException,
  UnauthorizedException,
} from '@nestjs/common';
import { ConfigService } from '@nestjs/config';

import { MailerService } from '@nestjs-modules/mailer';
import * as bcrypt from 'bcrypt';
import * as cacheManager_1 from 'cache-manager';
import { randomBytes } from 'crypto';
import { Response } from 'express';

import { generateUniqueId } from '@/common/utils/id-generator';
import { PrismaService } from '@/prisma/prisma.service';

import { SignInDto } from './dto/sign-in.dto';
import { SignUpDto } from './dto/sign-up.dto';

@Injectable()
export class AuthService {
  baseUrl: string;
  private readonly logger = new Logger('Auth');
  constructor(
    private readonly prisma: PrismaService,
    private readonly configService: ConfigService,
    private readonly mailerSerivce: MailerService,
    @Inject(CACHE_MANAGER) private cacheManager: cacheManager_1.Cache,
  ) {
    this.baseUrl = this.configService.get<string>(
      'BASE_URL',
      'http://localhost:3000',
    );
  }

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

    const role = await this.prisma.role.findUnique({
      where: { name: 'default' },
    });
    if (!role) {
      throw new NotFoundException('Роль default не найдена');
    }
    try {
      // Генерируем уникальный семизначный ID
      const userId = await generateUniqueId(this.prisma, 'user');

      await this.prisma.user.create({
        data: {
          id: userId,
          ...dto,
          password: hashedPassword,
          roleId: role?.id,
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

    // Генерируем ID сессии
    const sessionId = this.generateSessionId();

    // Сохраняем данные сессии в Redis (30 дней)
    await this.cacheManager.set(
      `session:${sessionId}`,
      JSON.stringify({
        userId: checkUser.id,
        email: checkUser.email,
        profileType: checkUser.profileType,
      }),
      30 * 24 * 60 * 60 * 1000, // 30 дней
    );

    // Устанавливаем session id в cookie
    response.cookie('session_id', sessionId, {
      httpOnly: true,
      secure: process.env.NODE_ENV === 'production',
      sameSite: process.env.NODE_ENV === 'production' ? 'strict' : 'lax',
      maxAge: 30 * 24 * 60 * 60 * 1000, // 30 дней
    });

    this.logger.log(`Успешная авторизация: ${login}`);

    return {
      message: 'Вы успешно авторизовались!',
      user: {
        id: checkUser.id,
        email: checkUser.email,
        fullName: checkUser.fullName,
        phoneNumber: checkUser.phoneNumber,
        profileType: checkUser.profileType,
        photo: checkUser.photo ? `${this.baseUrl}${checkUser.photo}` : null,
      },
    };
  }

  async logout(sessionId: string, response: Response) {
    await this.cacheManager.del(`session:${sessionId}`);

    response.clearCookie('session_id');

    return { message: 'Вы успешно вышли из системы!' };
  }

  async getSessionData(sessionId: string) {
    const sessionDataStr = await this.cacheManager.get<string>(
      `session:${sessionId}`,
    );
    if (!sessionDataStr) {
      return null;
    }
    return JSON.parse(sessionDataStr);
  }

  async getCurrentUser(userId: number) {
    const user = await this.prisma.user.findUnique({
      where: { id: userId },
      include: {
        role: true,
      },
    });

    if (!user) {
      throw new NotFoundException('Пользователь не найден');
    }

    return {
      id: user.id,
      email: user.email,
      fullName: user.fullName,
      phoneNumber: user.phoneNumber,
      profileType: user.profileType,
      photo: user.photo ? `${this.baseUrl}${user.photo}` : null,
      rating: user.rating,
      isAnswersCall: user.isAnswersCall,
      role: user.role?.name,
    };
  }

  private generateSessionId(): string {
    return randomBytes(32).toString('hex');
  }

  async forgotPassword(email: string) {
    const checkUser = await this.prisma.user.findUnique({
      where: { email },
    });

    if (!checkUser) {
      throw new BadRequestException(
        'Пользователя с такой почтой не существует',
      );
    }

    const code = await this.generateVerifyCode();

    await this.cacheManager.set(
      `forgot-password:${code}`,
      JSON.stringify({
        id: checkUser.id.toString(),
        code,
      }),
      3600 * 1000,
    );

    await this.sendVerificationEmail(
      email,
      'Код восстановления пароля - Торгуй Сам',
      code,
      'forgot-password',
    );

    return { message: 'Письмо с кодом подтверждения отправлено на почту' };
  }

  async verifyCode(code: string) {
    const cachedDataStr = await this.cacheManager.get<string>(
      `forgot-password:${code}`,
    );
    const cachedData = cachedDataStr ? JSON.parse(cachedDataStr) : null;

    if (!cachedData) {
      console.log('Данные не найдены в кеше');
    }
    if (cachedData.code !== code) {
      throw new BadRequestException('Неверный код подтверждения');
    }
    const user = await this.prisma.user.findUnique({
      where: { id: +cachedData.id },
    });

    if (!user) {
      throw new NotFoundException('Такого пользователя не существует');
    }

    await this.prisma.user.update({
      where: { id: user.id },
      data: {
        isResetVerified: true,
      },
    });

    await this.cacheManager.del(`forgot-password:${code}`);
    return { userId: user.id };
  }

  async changePassword(userId: number, password: string) {
    const user = await this.prisma.user.findUnique({
      where: { id: userId },
    });

    if (!user) {
      throw new NotFoundException('Пользователь не найден');
    }

    if (!user.isResetVerified) {
      throw new ForbiddenException('Требуется подтверждение сброса пароля');
    }

    await this.prisma.user.update({
      where: { id: userId },
      data: {
        password: await bcrypt.hash(password, 10),
        isResetVerified: false,
      },
    });

    return { message: 'Пароль успешно изменён' };
  }

  private async sendVerificationEmail(
    email: string,
    text: string,
    code: string,
    template: string,
  ) {
    try {
      this.logger.log(`Отправка письма на ${email} с кодом ${code}`);

      const result = await this.mailerSerivce.sendMail({
        to: email,
        subject: text,
        template,
        context: {
          code,
        },
      });

      this.logger.log('Письмо успешно отправлено:', result);
    } catch (error) {
      this.logger.error('Ошибка отправки письма:', error);
      throw new BadRequestException(`Ошибка отправки письма: ${error.message}`);
    }
  }

  async generateVerifyCode(): Promise<string> {
    return Math.floor(100000 + Math.random() * 900000).toString(); // Генерирует 6 цифр (от 100000 до 999999)
  }
}
