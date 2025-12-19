import {
  BadRequestException,
  Inject,
  Injectable,
  Logger,
  NotFoundException,
} from '@nestjs/common';
import { PrismaService } from 'src/prisma/prisma.service';
import { UpdateSettingsDto } from './dto/update-settings.dto';
import { ConfigService } from '@nestjs/config';
import { MailerService } from '@nestjs-modules/mailer';
import * as cacheManager_1 from 'cache-manager';
import { CACHE_MANAGER } from '@nestjs/cache-manager';
import { S3Service } from 'src/s3/s3.service';

@Injectable()
export class UserService {
  baseUrl: string;
  profileTypes = {
    IP: 'Индивидуальный предприниматель',
    OOO: 'Юридическое лицо',
    INDIVIDUAL: 'Физическое лицо',
  };
  private readonly logger = new Logger('Auth');
  constructor(
    private readonly prisma: PrismaService,
    private readonly configService: ConfigService,
    private readonly mailerSerivce: MailerService,
    @Inject(CACHE_MANAGER) private cacheManager: cacheManager_1.Cache,
    private readonly s3Service: S3Service,
  ) {
    this.baseUrl = this.configService.get<string>(
      'BASE_URL',
      'http://localhost:3000',
    );
  }

  async findAll() {
    const users = await this.prisma.user.findMany({
      select: {
        id: true,
        isBanned: true,
        fullName: true,
        email: true,
        phoneNumber: true,
        rating: true,
        profileType: true,
        photo: true,
        balance: true,
        bonusBalance: true,
        _count: {
          select: {
            products: true,
          },
        },
      },
    });

    return users.map((user) => ({
      ...user,
      products: user._count.products,
      _count: undefined,
    }));
  }

  async toggleBanned(id: number) {
    const checkUser = await this.prisma.user.findUnique({
      where: { id },
      select: { id: true, isBanned: true },
    });

    if (!checkUser) {
      throw new NotFoundException('Пользователь не найден');
    }

    await this.prisma.user.update({
      where: { id },
      data: {
        isBanned: !checkUser.isBanned,
      },
    });

    return checkUser.isBanned
      ? { message: 'Пользователь успешно разбанен' }
      : { message: 'Пользователь успешно забанен' };
  }

  async deleteUser(id: number) {
    const checkUser = await this.prisma.user.findUnique({
      where: { id },
    });

    if (!checkUser) {
      throw new NotFoundException('Пользователь не найден');
    }

    await this.prisma.user.delete({
      where: { id },
    });

    return { message: 'Пользователь успешно удалён' };
  }

  async getUserInfo(id: number) {
    const checkUser = await this.prisma.user.findUnique({
      where: { id },
      select: {
        id: true,
        fullName: true,
        profileType: true,
        phoneNumber: true,
        balance: true,
        bonusBalance: true,
        photo: true,
      },
    });

    const reviews = await this.prisma.review.findMany({
      where: { reviewedUserId: id },
      select: {
        rating: true,
      },
    });

    if (!checkUser) {
      throw new BadRequestException('Пользователь не найден');
    }

    const { bonusBalance, ...userWithoutBonus } = checkUser;

    return {
      ...userWithoutBonus,
      balance: checkUser.balance + bonusBalance,
      profileType: this.profileTypes[checkUser.profileType],
      rating:
        reviews.reduce((sum, review) => sum + review.rating, 0) /
          reviews.length || 0,
      reviewsCount: reviews.length,
    };
  }

  async getPhoneNumber(userId: number, myId: number) {
    const checkUser = await this.prisma.user.findUnique({
      where: { id: userId },
    });

    if (!checkUser) {
      throw new NotFoundException('Продавец с данным id не найден');
    }

    if (!checkUser.phoneNumber) {
      throw new BadRequestException('У продавца отсутствует номер телефона');
    }

    // Проверяем, не смотрел ли уже этот пользователь номер этого продавца
    const existingView = await this.prisma.phoneNumberView.findUnique({
      where: {
        viewedById_viewedUserId: {
          viewedById: myId,
          viewedUserId: userId,
        },
      },
    });

    // Если ещё не смотрел, записываем в статистику
    if (!existingView) {
      await this.prisma.phoneNumberView.create({
        data: {
          viewedById: myId,
          viewedUserId: userId,
        },
      });
    }

    return {
      phoneNumber: checkUser.phoneNumber,
    };
  }

  async updateSettings(
    dto: UpdateSettingsDto,
    userId: number,
    file: Express.Multer.File | null,
  ) {
    const updatedData = { ...dto };

    // Если загружено новое фото
    if (file) {
      // Получаем текущего пользователя для удаления старого фото
      const currentUser = await this.prisma.user.findUnique({
        where: { id: userId },
        select: { photo: true },
      });

      // Удаляем старое фото из S3, если оно есть
      if (currentUser?.photo) {
        await this.s3Service.deleteFile(currentUser.photo);
      }

      // Загружаем новое фото в S3
      const photoUrl = await this.s3Service.uploadFile(file, 'users');
      updatedData['photo'] = photoUrl;
    }

    await this.prisma.user.update({
      where: { id: userId },
      data: {
        ...updatedData,
      },
    });

    return { ...updatedData };
  }

  async getProfileSettings(userId: number) {
    const checkUser = await this.prisma.user.findUnique({
      where: { id: userId },
      select: {
        id: true,
        fullName: true,
        phoneNumber: true,
        isAnswersCall: true,
        profileType: true,
        photo: true,
      },
    });

    if (!userId) {
      throw new NotFoundException('Пользователь не найден');
    }

    return {
      ...checkUser,
      photo: checkUser?.photo || null, // URL уже полный из S3
    };
  }

  async verifyCode(code: string) {
    const cachedDataStr = await this.cacheManager.get<string>(
      `verify-email:${code}`,
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
        isEmailVerified: true,
      },
    });

    await this.cacheManager.del(`verify-email:${code}`);
    return { message: 'Почта успешно подтверждена' };
  }

  async verifyEmail(userId: number) {
    const currentUser = await this.prisma.user.findUnique({
      where: { id: userId },
    });

    if (!currentUser) {
      throw new NotFoundException('Пользователь не найден');
    }

    const code = await this.generateVerifyCode();

    await this.cacheManager.set(
      `verify-email:${code}`,
      JSON.stringify({
        id: currentUser.id.toString(),
        code,
      }),
      3600 * 1000,
    );

    await this.sendVerificationEmail(
      currentUser.email,
      'Код подтверждения почты - Торгуй Сам',
      code,
      'verify-email',
    );

    return { message: 'Письмо с кодом подтверждения отправлено на почту' };
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

  async setBalance(userId: number, balance: string) {
    const checkUser = await this.prisma.user.findUnique({
      where: { id: userId },
    });

    if (!checkUser) {
      throw new NotFoundException('Пользователь не найден');
    }

    const parseBalance = parseFloat(balance);

    if (isNaN(parseBalance)) {
      throw new BadRequestException('Неудалось преобразовать баланс в float');
    }
    if (parseBalance < 0) {
      throw new BadRequestException('Нельзя установить отрицательный баланс');
    }

    await this.prisma.user.update({
      where: { id: userId },
      data: {
        bonusBalance: parseBalance,
      },
    });

    return { message: 'Баланс успешно обновлён' };
  }

  async updateUser(userId: number, dto: any) {
    const checkUser = await this.prisma.user.findUnique({
      where: { id: userId },
    });

    if (!checkUser) {
      throw new NotFoundException('Пользователь не найден');
    }

    // Проверяем уникальность email если он обновляется
    if (dto.email && dto.email !== checkUser.email) {
      const emailExists = await this.prisma.user.findUnique({
        where: { email: dto.email },
      });

      if (emailExists) {
        throw new BadRequestException(
          'Пользователь с таким email уже существует',
        );
      }
    }

    // Проверяем уникальность phoneNumber если он обновляется
    if (dto.phoneNumber && dto.phoneNumber !== checkUser.phoneNumber) {
      const phoneExists = await this.prisma.user.findUnique({
        where: { phoneNumber: dto.phoneNumber },
      });

      if (phoneExists) {
        throw new BadRequestException(
          'Пользователь с таким номером телефона уже существует',
        );
      }
    }

    await this.prisma.user.update({
      where: { id: userId },
      data: {
        ...dto,
      },
    });

    return { message: 'Данные пользователя успешно обновлены' };
  }
}
