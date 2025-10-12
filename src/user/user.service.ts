import { BadRequestException, Injectable } from '@nestjs/common';
import { PrismaService } from 'src/prisma/prisma.service';

@Injectable()
export class UserService {
  profileTypes = {
    IP: 'Индивидуальный предприниматель',
    OOO: 'Юридическое лицо',
    INDIVIDUAL: 'Физическое лицо',
  };
  constructor(private readonly prisma: PrismaService) {}

  async getUserInfo(id: number) {
    const checkUser = await this.prisma.user.findUnique({
      where: { id },
      select: {
        id: true,
        fullName: true,
        profileType: true,
        phoneNumber: true,
      },
    });

    if (!checkUser) {
      throw new BadRequestException('Пользователь не найден');
    }

    return {
      ...checkUser,
      profileType: this.profileTypes[checkUser.profileType],
    };
  }
}
