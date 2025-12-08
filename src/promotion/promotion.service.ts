import {
  BadRequestException,
  ForbiddenException,
  Injectable,
  NotFoundException,
} from '@nestjs/common';
import { PrismaService } from 'src/prisma/prisma.service';
import { AddPromotionDto } from './add-promotion.dto';

@Injectable()
export class PromotionService {
  constructor(private readonly prisma: PrismaService) {}

  async allPromotions() {
    const promotions = await this.prisma.promotion.findMany({
      select: {
        id: true,
        name: true,
        pricePerDay: true,
      },
    });

    return promotions;
  }

  async addPromotion(dto: AddPromotionDto, userId: number) {
    const { productId, promotionId, days } = { ...dto };

    const checkProduct = await this.prisma.product.findUnique({
      where: { id: productId },
      include: {
        user: true,
      },
    });

    if (!checkProduct) {
      throw new NotFoundException('Товар не найден');
    }

    if (checkProduct.user.id !== userId) {
      throw new ForbiddenException('Вы не можете продвигать чужой товар');
    }

    const checkPromo = await this.prisma.promotion.findUnique({
      where: {
        id: promotionId,
      },
    });

    if (!checkPromo) {
      throw new NotFoundException('Тип продвижения не найден');
    }

    const totalPrice = days * checkPromo.pricePerDay;

    // Проверяем баланс
    if (checkProduct.user.balance < totalPrice) {
      throw new BadRequestException(
        `На балансе недостаточно средств. Требуется: ${totalPrice}₽, доступно: ${checkProduct.user.balance}₽`,
      );
    }

    // Рассчитываем даты
    const startDate = new Date();
    const endDate = new Date();
    endDate.setDate(endDate.getDate() + days);

    // Создаем продвижение
    const promotion = await this.prisma.productPromotion.create({
      data: {
        productId,
        userId,
        promotionId,
        days,
        totalPrice,
        startDate,
        endDate,
        isActive: true,
        isPaid: false,
      },
    });

    // Списываем средства с баланса
    await this.prisma.user.update({
      where: { id: userId },
      data: {
        balance: {
          decrement: totalPrice,
        },
      },
    });

    // Помечаем как оплаченное
    await this.prisma.productPromotion.update({
      where: { id: promotion.id },
      data: {
        isPaid: true,
      },
    });

    return {
      message: 'Продвижение успешно активировано',
      promotion: {
        id: promotion.id,
        days,
        totalPrice,
        startDate,
        endDate,
      },
    };
  }
}
