import { Module } from '@nestjs/common';
import { AuthModule } from './auth/auth.module';
import { PrismaModule } from './prisma/prisma.module';
import { UserModule } from './user/user.module';
import { ProductModule } from './product/product.module';
import { ConfigModule } from '@nestjs/config';
import { CategoryModule } from './category/category.module';
import { MailModule } from './mail/mail.module';
import { CacheModule } from '@nestjs/cache-manager';
import { ChatModule } from './chat/chat.module';
import { StatisticsModule } from './statistics/statistics.module';
import { ReviewModule } from './review/review.module';
import { SupportModule } from './support/support.module';
import { SubcategoryModule } from './subcategory/subcategory.module';
import { AddressModule } from './address/address.module';
import { S3Module } from './s3/s3.module';
import { SubcategoryTypeModule } from './subcategory-type/subcategory-type.module';
import { TypeFieldModule } from './type-field/type-field.module';
import { PromotionModule } from './promotion/promotion.module';
import { BannerModule } from './banner/banner.module';
import { PaymentModule } from './payment/payment.module';
import * as redisStore from 'cache-manager-ioredis';

@Module({
  imports: [
    AuthModule,
    PrismaModule,
    UserModule,
    ProductModule,
    ConfigModule.forRoot({
      isGlobal: true,
    }),
    CategoryModule,
    MailModule,
    CacheModule.register({
      ttl: 60 * 60 * 100,
      isGlobal: true,
      store: redisStore,
      // host: 'localhost',
      host: 'redis',
      port: 6379,
    }),
    ChatModule,
    StatisticsModule,
    ReviewModule,
    SupportModule,
    SubcategoryModule,
    AddressModule,
    S3Module,
    SubcategoryTypeModule,
    TypeFieldModule,
    PromotionModule,
    BannerModule,
    PaymentModule,
  ],
})
export class AppModule {}
