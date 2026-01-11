import { CACHE_MANAGER } from '@nestjs/cache-manager';
import {
  BadRequestException,
  ForbiddenException,
  Inject,
  Injectable,
  NotFoundException,
} from '@nestjs/common';
import { ConfigService } from '@nestjs/config';

import * as cacheManager from 'cache-manager';

import { PrismaService } from '@/prisma/prisma.service';

@Injectable()
export class ChatService {
  constructor(
    private readonly prisma: PrismaService,
    private readonly configService: ConfigService,
    @Inject(CACHE_MANAGER) private cacheManager: cacheManager.Cache,
  ) {}

  async getSessionData(sessionId: string) {
    const sessionDataStr = await this.cacheManager.get<string>(
      `session:${sessionId}`,
    );
    if (!sessionDataStr) {
      return null;
    }
    return JSON.parse(sessionDataStr);
  }

  async getUserById(userId: number) {
    return await this.prisma.user.findUnique({
      where: { id: userId },
      include: { role: true },
    });
  }

  // Начать чат по товару
  async startChat(productId: number, buyerId: number) {
    // Проверяем, что товар существует
    const product = await this.prisma.product.findUnique({
      where: { id: productId },
      include: {
        user: {
          select: {
            id: true,
            fullName: true,
            phoneNumber: true,
          },
        },
      },
    });

    if (!product) {
      throw new NotFoundException('Товар не найден');
    }

    // Проверяем, что покупатель не пытается написать самому себе
    if (product.userId === buyerId) {
      throw new BadRequestException('Нельзя писать самому себе');
    }

    // Проверяем, существует ли уже чат между этими пользователями по этому товару
    const existingChat = await this.prisma.chat.findFirst({
      where: {
        productId,
        buyerId,
        sellerId: product.userId,
      },
      include: {
        product: {
          select: {
            id: true,
            name: true,
            images: true,
            price: true,
          },
        },
        seller: {
          select: {
            id: true,
            fullName: true,
            phoneNumber: true,
          },
        },
        buyer: {
          select: {
            id: true,
            fullName: true,
            phoneNumber: true,
          },
        },
      },
    });

    if (existingChat) {
      return {
        ...existingChat,
        product: existingChat.product
          ? {
              ...existingChat.product,
              images: existingChat.product.images.map((image) => {
                image;
              }),
            }
          : null,
      };
    }

    // Создаем новый чат
    const newChat = await this.prisma.chat.create({
      data: {
        productId,
        buyerId,
        sellerId: product.userId,
      },
      include: {
        product: {
          select: {
            id: true,
            name: true,
            images: true,
            price: true,
          },
        },
        seller: {
          select: {
            id: true,
            fullName: true,
            phoneNumber: true,
          },
        },
        buyer: {
          select: {
            id: true,
            fullName: true,
            phoneNumber: true,
          },
        },
      },
    });

    return {
      ...newChat,
      product: newChat.product
        ? {
            ...newChat.product,
            images: newChat.product.images.map((image) => `${image}`),
          }
        : null,
    };
  }

  // Получить все чаты пользователя
  async getUserChats(userId: number) {
    const chats = await this.prisma.chat.findMany({
      where: {
        OR: [{ buyerId: userId }, { sellerId: userId }],
      },
      include: {
        product: {
          select: {
            id: true,
            name: true,
            images: true,
            price: true,
          },
        },
        buyer: {
          select: {
            id: true,
            fullName: true,
            phoneNumber: true,
          },
        },
        seller: {
          select: {
            id: true,
            fullName: true,
            phoneNumber: true,
          },
        },
        lastMessage: {
          select: {
            id: true,
            content: true,
            createdAt: true,
            senderId: true,
            isRead: true,
          },
        },
      },
      orderBy: {
        lastMessageAt: 'desc',
      },
    });

    // Форматируем данные для фронтенда
    return chats.map((chat) => {
      const isUserBuyer = chat.buyerId === userId;
      const companion = isUserBuyer ? chat.seller : chat.buyer;
      const unreadCount = isUserBuyer
        ? chat.unreadCountBuyer
        : chat.unreadCountSeller;

      // Форматируем дату в dd.mm.yy
      const formatDate = (date: Date) => {
        const day = date.getDate().toString().padStart(2, '0');
        const month = (date.getMonth() + 1).toString().padStart(2, '0');
        const year = date.getFullYear().toString().slice(-2);
        return `${day}.${month}.${year}`;
      };

      return {
        id: chat.id,
        isModerationChat: chat.isModerationChat, // Признак чата модерации
        // Информация о товаре (может быть null для чатов модерации)
        product: chat.product
          ? {
              id: chat.product.id,
              name: chat.product.name,
              price: chat.product.price,
              image: chat.product.images[0]
                ? `${chat.product.images[0]}`
                : null, // Первое фото товара
            }
          : null,
        // Информация о собеседнике (продавце для покупателя, покупателе для продавца)
        companion: {
          id: companion.id,
          fullName: companion.fullName,
          phoneNumber: companion.phoneNumber,
          // TODO: Добавить аватарку когда будет поле в User
          avatar: null,
        },
        // Последнее сообщение
        lastMessage: chat.lastMessage
          ? {
              content: chat.lastMessage.content,
              createdAt: chat.lastMessage.createdAt,
              formattedDate: formatDate(chat.lastMessage.createdAt),
              isFromMe: chat.lastMessage.senderId === userId,
              isRead: chat.lastMessage.isRead,
            }
          : null,
        // Счетчик непрочитанных
        unreadCount,
        // Дата последнего обновления чата
        lastActivity: formatDate(chat.lastMessageAt || chat.createdAt),
        createdAt: chat.createdAt,
      };
    });
  }

  // Получить информацию о чате с деталями товара и продавца
  async getChatInfo(chatId: number, userId: number) {
    const chat = await this.prisma.chat.findUnique({
      where: { id: chatId },
      include: {
        product: {
          select: {
            id: true,
            name: true,
            images: true,
            price: true,
            description: true,
          },
        },
        buyer: {
          select: {
            id: true,
            fullName: true,
            phoneNumber: true,
          },
        },
        seller: {
          select: {
            id: true,
            fullName: true,
            phoneNumber: true,
          },
        },
      },
    });

    if (!chat) {
      throw new NotFoundException('Чат не найден');
    }

    // Проверяем доступ к чату
    if (chat.buyerId !== userId && chat.sellerId !== userId) {
      throw new ForbiddenException('Нет доступа к этому чату');
    }

    const isUserBuyer = chat.buyerId === userId;
    const companion = isUserBuyer ? chat.seller : chat.buyer;

    return {
      id: chat.id,
      isModerationChat: chat.isModerationChat,
      // Информация о товаре (отображается сверху в чате, может быть null для чатов модерации)
      product: chat.product
        ? {
            id: chat.product.id,
            name: chat.product.name,
            price: chat.product.price,
            image: chat.product.images[0] ? `${chat.product.images[0]}` : null,
            description: chat.product.description,
          }
        : null,
      // Информация о собеседнике (продавце)
      companion: {
        id: companion.id,
        fullName: companion.fullName,
        phoneNumber: companion.phoneNumber,
        role: isUserBuyer ? 'seller' : 'buyer',
      },
      // Мета информация
      isUserBuyer,
      unreadCount: isUserBuyer ? chat.unreadCountBuyer : chat.unreadCountSeller,
      createdAt: chat.createdAt,
    };
  }

  // Получить сообщения чата с пагинацией
  async getChatMessages(
    chatId: number,
    userId: number,
    page: number = 1,
    limit: number = 50,
  ) {
    // Проверяем доступ к чату
    const chat = await this.prisma.chat.findUnique({
      where: { id: chatId },
      select: { buyerId: true, sellerId: true },
    });

    if (!chat) {
      throw new NotFoundException('Чат не найден');
    }

    if (chat.buyerId !== userId && chat.sellerId !== userId) {
      throw new ForbiddenException('Нет доступа к этому чату');
    }

    // Получаем сообщения с пагинацией
    const messages = await this.prisma.message.findMany({
      where: { chatId },
      include: {
        sender: {
          select: {
            id: true,
            fullName: true,
          },
        },
        relatedProduct: {
          select: {
            id: true,
            name: true,
            images: true,
            price: true,
            moderateState: true,
          },
        },
      },
      orderBy: {
        createdAt: 'desc',
      },
      skip: (page - 1) * limit,
      take: limit,
    });

    // Форматируем сообщения
    const formattedMessages = messages.reverse().map((message) => ({
      id: message.id,
      content: message.content,
      senderId: message.senderId,
      sender: message.sender,
      relatedProduct: message.relatedProduct, // Добавляем информацию о связанном товаре
      isFromMe: message.senderId === userId,
      isRead: message.isRead,
      readAt: message.readAt,
      createdAt: message.createdAt,
      // Время отправки в удобном формате
      timeString: message.createdAt.toLocaleTimeString('ru-RU', {
        hour: '2-digit',
        minute: '2-digit',
      }),
    }));

    // Подсчитываем общее количество сообщений
    const total = await this.prisma.message.count({
      where: { chatId },
    });

    return {
      messages: formattedMessages,
      pagination: {
        page,
        limit,
        total,
        pages: Math.ceil(total / limit),
      },
    };
  }

  // Отправить сообщение
  async sendMessage(chatId: number, senderId: number, content: string) {
    // Проверяем доступ к чату
    const chat = await this.prisma.chat.findUnique({
      where: { id: chatId },
    });

    if (!chat) {
      throw new NotFoundException('Чат не найден');
    }

    if (chat.buyerId !== senderId && chat.sellerId !== senderId) {
      throw new ForbiddenException('Нет доступа к этому чату');
    }

    // Создаем сообщение
    const message = await this.prisma.message.create({
      data: {
        content,
        senderId,
        chatId,
      },
      include: {
        sender: {
          select: {
            id: true,
            fullName: true,
          },
        },
      },
    });

    // Обновляем счетчики непрочитанных сообщений и последнее сообщение
    const isFromBuyer = senderId === chat.buyerId;

    await this.prisma.chat.update({
      where: { id: chatId },
      data: {
        lastMessageId: message.id,
        lastMessageAt: new Date(),
        // Увеличиваем счетчик непрочитанных для получателя
        ...(isFromBuyer
          ? { unreadCountSeller: { increment: 1 } }
          : { unreadCountBuyer: { increment: 1 } }),
      },
    });

    return {
      id: message.id,
      content: message.content,
      senderId: message.senderId,
      sender: message.sender,
      isRead: message.isRead,
      createdAt: message.createdAt,
      timeString: message.createdAt.toLocaleTimeString('ru-RU', {
        hour: '2-digit',
        minute: '2-digit',
      }),
    };
  }

  // Отправить сообщение с привязкой к товару (для чатов модерации)
  async sendMessageWithProduct(
    chatId: number,
    senderId: number,
    content: string,
    productId: number,
  ) {
    // Проверяем доступ к чату
    const chat = await this.prisma.chat.findUnique({
      where: { id: chatId },
    });

    if (!chat) {
      throw new NotFoundException('Чат не найден');
    }

    if (chat.buyerId !== senderId && chat.sellerId !== senderId) {
      throw new ForbiddenException('Нет доступа к этому чату');
    }

    // Создаем сообщение с привязкой к товару
    const message = await this.prisma.message.create({
      data: {
        content,
        senderId,
        chatId,
        relatedProductId: productId, // Привязываем к товару
      },
      include: {
        sender: {
          select: {
            id: true,
            fullName: true,
          },
        },
        relatedProduct: {
          select: {
            id: true,
            name: true,
            images: true,
            price: true,
          },
        },
      },
    });

    // Обновляем счетчики непрочитанных сообщений и последнее сообщение
    const isFromBuyer = senderId === chat.buyerId;

    await this.prisma.chat.update({
      where: { id: chatId },
      data: {
        lastMessageId: message.id,
        lastMessageAt: new Date(),
        // Увеличиваем счетчик непрочитанных для получателя
        ...(isFromBuyer
          ? { unreadCountSeller: { increment: 1 } }
          : { unreadCountBuyer: { increment: 1 } }),
      },
    });

    return {
      id: message.id,
      content: message.content,
      senderId: message.senderId,
      sender: message.sender,
      relatedProduct: message.relatedProduct,
      isRead: message.isRead,
      createdAt: message.createdAt,
      timeString: message.createdAt.toLocaleTimeString('ru-RU', {
        hour: '2-digit',
        minute: '2-digit',
      }),
    };
  }

  // Отметить сообщения как прочитанные
  async markMessagesAsRead(chatId: number, userId: number) {
    // Проверяем доступ к чату
    const chat = await this.prisma.chat.findUnique({
      where: { id: chatId },
    });

    if (!chat) {
      throw new NotFoundException('Чат не найден');
    }

    if (chat.buyerId !== userId && chat.sellerId !== userId) {
      throw new ForbiddenException('Нет доступа к этому чату');
    }

    // Отмечаем все непрочитанные сообщения как прочитанные
    await this.prisma.message.updateMany({
      where: {
        chatId,
        senderId: { not: userId }, // Не свои сообщения
        isRead: false,
      },
      data: {
        isRead: true,
        readAt: new Date(),
      },
    });

    // Обнуляем счетчик непрочитанных сообщений
    const isUserBuyer = chat.buyerId === userId;

    await this.prisma.chat.update({
      where: { id: chatId },
      data: isUserBuyer ? { unreadCountBuyer: 0 } : { unreadCountSeller: 0 },
    });

    return { success: true, message: 'Сообщения отмечены как прочитанные' };
  }
}
