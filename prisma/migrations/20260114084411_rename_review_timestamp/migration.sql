/*
  Warnings:

  - The values [PENDING] on the enum `TicketStatus` will be removed. If these variants are still used in the database, this will fail.
  - You are about to drop the column `brand` on the `Product` table. All the data in the column will be lost.
  - You are about to drop the column `model` on the `Product` table. All the data in the column will be lost.
  - You are about to drop the column `reviewedAt` on the `Review` table. All the data in the column will be lost.
  - You are about to drop the column `content` on the `SupportMessage` table. All the data in the column will be lost.
  - You are about to drop the column `createdAt` on the `SupportMessage` table. All the data in the column will be lost.
  - You are about to drop the column `isRead` on the `SupportMessage` table. All the data in the column will be lost.
  - You are about to drop the column `readAt` on the `SupportMessage` table. All the data in the column will be lost.
  - You are about to drop the column `senderId` on the `SupportMessage` table. All the data in the column will be lost.
  - You are about to drop the column `updatedAt` on the `SupportMessage` table. All the data in the column will be lost.
  - You are about to drop the column `assignedToId` on the `SupportTicket` table. All the data in the column will be lost.
  - You are about to drop the column `description` on the `SupportTicket` table. All the data in the column will be lost.
  - You are about to drop the column `lastMessageAt` on the `SupportTicket` table. All the data in the column will be lost.
  - You are about to drop the column `lastMessageId` on the `SupportTicket` table. All the data in the column will be lost.
  - You are about to drop the column `unreadCountModerator` on the `SupportTicket` table. All the data in the column will be lost.
  - You are about to drop the column `unreadCountUser` on the `SupportTicket` table. All the data in the column will be lost.
  - You are about to drop the column `refreshToken` on the `User` table. All the data in the column will be lost.
  - You are about to drop the column `refreshTokenExpiresAt` on the `User` table. All the data in the column will be lost.
  - A unique constraint covering the columns `[slug]` on the table `Category` will be added. If there are existing duplicate values, this will fail.
  - A unique constraint covering the columns `[slug,categoryId]` on the table `SubCategory` will be added. If there are existing duplicate values, this will fail.
  - A unique constraint covering the columns `[slug,subcategoryId]` on the table `SubcategotyType` will be added. If there are existing duplicate values, this will fail.
  - Made the column `address` on table `Product` required. This step will fail if there are existing NULL values in that column.
  - Added the required column `authorId` to the `SupportMessage` table without a default value. This is not possible if the table is not empty.
  - Added the required column `text` to the `SupportMessage` table without a default value. This is not possible if the table is not empty.
  - Added the required column `theme` to the `SupportTicket` table without a default value. This is not possible if the table is not empty.

*/
-- CreateEnum
CREATE TYPE "ProductModerate" AS ENUM ('MODERATE', 'APPROVED', 'DENIDED');

-- CreateEnum
CREATE TYPE "BannerPlace" AS ENUM ('PRODUCT_FEED', 'PROFILE', 'FAVORITES', 'CHATS');

-- CreateEnum
CREATE TYPE "TicketTheme" AS ENUM ('TECHNICAL_ISSUE', 'ACCOUNT_PROBLEM', 'PAYMENT_ISSUE', 'PRODUCT_QUESTION', 'COMPLAINT', 'SUGGESTION', 'OTHER');

-- AlterEnum
BEGIN;
CREATE TYPE "TicketStatus_new" AS ENUM ('OPEN', 'IN_PROGRESS', 'RESOLVED', 'CLOSED');
ALTER TABLE "public"."SupportTicket" ALTER COLUMN "status" DROP DEFAULT;
ALTER TABLE "SupportTicket" ALTER COLUMN "status" TYPE "TicketStatus_new" USING ("status"::text::"TicketStatus_new");
ALTER TYPE "TicketStatus" RENAME TO "TicketStatus_old";
ALTER TYPE "TicketStatus_new" RENAME TO "TicketStatus";
DROP TYPE "public"."TicketStatus_old";
ALTER TABLE "SupportTicket" ALTER COLUMN "status" SET DEFAULT 'OPEN';
COMMIT;

-- DropForeignKey
ALTER TABLE "SupportMessage" DROP CONSTRAINT "SupportMessage_senderId_fkey";

-- DropForeignKey
ALTER TABLE "SupportTicket" DROP CONSTRAINT "SupportTicket_assignedToId_fkey";

-- DropForeignKey
ALTER TABLE "SupportTicket" DROP CONSTRAINT "SupportTicket_lastMessageId_fkey";

-- DropForeignKey
ALTER TABLE "User" DROP CONSTRAINT "User_roleId_fkey";

-- AlterTable
ALTER TABLE "Category" ADD COLUMN     "createdAt" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
ADD COLUMN     "slug" TEXT NOT NULL DEFAULT '',
ADD COLUMN     "updatedAt" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP;

-- AlterTable
ALTER TABLE "Chat" ADD COLUMN     "isModerationChat" BOOLEAN NOT NULL DEFAULT false,
ALTER COLUMN "productId" DROP NOT NULL;

-- AlterTable
ALTER TABLE "Message" ADD COLUMN     "relatedProductId" INTEGER;

-- AlterTable
ALTER TABLE "Product" DROP COLUMN "brand",
DROP COLUMN "model",
ADD COLUMN     "isHide" BOOLEAN NOT NULL DEFAULT false,
ADD COLUMN     "moderateState" "ProductModerate" NOT NULL DEFAULT 'MODERATE',
ADD COLUMN     "moderationRejectionReason" TEXT,
ADD COLUMN     "typeId" INTEGER,
ADD COLUMN     "videoUrl" TEXT,
ALTER COLUMN "id" DROP DEFAULT,
ALTER COLUMN "address" SET NOT NULL;
DROP SEQUENCE "Product_id_seq";

-- AlterTable
ALTER TABLE "Review" DROP COLUMN "reviewedAt",
ADD COLUMN     "createdAt" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP;

-- AlterTable
ALTER TABLE "SubCategory" ADD COLUMN     "createdAt" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
ADD COLUMN     "slug" TEXT NOT NULL DEFAULT '',
ADD COLUMN     "updatedAt" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP;

-- AlterTable
ALTER TABLE "SubcategotyType" ADD COLUMN     "createdAt" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
ADD COLUMN     "slug" TEXT NOT NULL DEFAULT '',
ADD COLUMN     "updatedAt" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP;

-- AlterTable
ALTER TABLE "SupportMessage" DROP COLUMN "content",
DROP COLUMN "createdAt",
DROP COLUMN "isRead",
DROP COLUMN "readAt",
DROP COLUMN "senderId",
DROP COLUMN "updatedAt",
ADD COLUMN     "authorId" INTEGER NOT NULL,
ADD COLUMN     "sentAt" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
ADD COLUMN     "text" TEXT NOT NULL;

-- AlterTable
ALTER TABLE "SupportTicket" DROP COLUMN "assignedToId",
DROP COLUMN "description",
DROP COLUMN "lastMessageAt",
DROP COLUMN "lastMessageId",
DROP COLUMN "unreadCountModerator",
DROP COLUMN "unreadCountUser",
ADD COLUMN     "moderatorId" INTEGER,
ADD COLUMN     "theme" "TicketTheme" NOT NULL;

-- AlterTable
ALTER TABLE "User" DROP COLUMN "refreshToken",
DROP COLUMN "refreshTokenExpiresAt",
ADD COLUMN     "balance" DOUBLE PRECISION NOT NULL DEFAULT 0,
ADD COLUMN     "bonusBalance" DOUBLE PRECISION NOT NULL DEFAULT 0,
ADD COLUMN     "freeAdsLimit" INTEGER NOT NULL DEFAULT 12,
ADD COLUMN     "isAnswersCall" BOOLEAN DEFAULT false,
ADD COLUMN     "isBanned" BOOLEAN NOT NULL DEFAULT false,
ADD COLUMN     "isEmailVerified" BOOLEAN NOT NULL DEFAULT false,
ADD COLUMN     "lastAdLimitReset" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
ADD COLUMN     "photo" TEXT,
ADD COLUMN     "usedFreeAds" INTEGER NOT NULL DEFAULT 0,
ALTER COLUMN "id" DROP DEFAULT,
ALTER COLUMN "roleId" DROP NOT NULL;
DROP SEQUENCE "User_id_seq";

-- CreateTable
CREATE TABLE "Log" (
    "id" SERIAL NOT NULL,
    "userId" INTEGER NOT NULL,
    "action" TEXT NOT NULL,

    CONSTRAINT "Log_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "Promotion" (
    "id" SERIAL NOT NULL,
    "name" TEXT NOT NULL,
    "pricePerDay" INTEGER NOT NULL,
    "createdAt" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updatedAt" TIMESTAMP(3) NOT NULL,

    CONSTRAINT "Promotion_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "ProductPromotion" (
    "id" SERIAL NOT NULL,
    "productId" INTEGER NOT NULL,
    "promotionId" INTEGER NOT NULL,
    "userId" INTEGER NOT NULL,
    "days" INTEGER NOT NULL,
    "totalPrice" INTEGER NOT NULL,
    "startDate" TIMESTAMP(3) NOT NULL,
    "endDate" TIMESTAMP(3) NOT NULL,
    "isActive" BOOLEAN NOT NULL DEFAULT true,
    "isPaid" BOOLEAN NOT NULL DEFAULT false,
    "createdAt" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updatedAt" TIMESTAMP(3) NOT NULL,

    CONSTRAINT "ProductPromotion_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "TypeField" (
    "id" SERIAL NOT NULL,
    "name" TEXT NOT NULL,
    "isRequired" BOOLEAN NOT NULL DEFAULT false,
    "typeId" INTEGER NOT NULL,

    CONSTRAINT "TypeField_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "ProductFieldValue" (
    "id" SERIAL NOT NULL,
    "value" TEXT NOT NULL,
    "fieldId" INTEGER NOT NULL,
    "productId" INTEGER NOT NULL,

    CONSTRAINT "ProductFieldValue_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "Banner" (
    "id" SERIAL NOT NULL,
    "photoUrl" TEXT NOT NULL,
    "place" "BannerPlace" NOT NULL,
    "navigateToUrl" TEXT NOT NULL,
    "createdAt" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updatedAt" TIMESTAMP(3) NOT NULL,

    CONSTRAINT "Banner_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "Payment" (
    "id" SERIAL NOT NULL,
    "orderId" TEXT NOT NULL,
    "paymentId" TEXT NOT NULL,
    "userId" INTEGER NOT NULL,
    "amount" DOUBLE PRECISION NOT NULL,
    "status" TEXT NOT NULL,
    "paymentUrl" TEXT,
    "createdAt" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updatedAt" TIMESTAMP(3) NOT NULL,

    CONSTRAINT "Payment_pkey" PRIMARY KEY ("id")
);

-- CreateIndex
CREATE UNIQUE INDEX "Promotion_name_key" ON "Promotion"("name");

-- CreateIndex
CREATE UNIQUE INDEX "ProductFieldValue_fieldId_productId_key" ON "ProductFieldValue"("fieldId", "productId");

-- CreateIndex
CREATE UNIQUE INDEX "Payment_orderId_key" ON "Payment"("orderId");

-- CreateIndex
CREATE UNIQUE INDEX "Payment_paymentId_key" ON "Payment"("paymentId");

-- CreateIndex
CREATE UNIQUE INDEX "Category_slug_key" ON "Category"("slug");

-- CreateIndex
CREATE INDEX "Category_slug_idx" ON "Category"("slug");

-- CreateIndex
CREATE INDEX "SubCategory_slug_idx" ON "SubCategory"("slug");

-- CreateIndex
CREATE UNIQUE INDEX "SubCategory_slug_categoryId_key" ON "SubCategory"("slug", "categoryId");

-- CreateIndex
CREATE INDEX "SubcategotyType_slug_idx" ON "SubcategotyType"("slug");

-- CreateIndex
CREATE UNIQUE INDEX "SubcategotyType_slug_subcategoryId_key" ON "SubcategotyType"("slug", "subcategoryId");

-- AddForeignKey
ALTER TABLE "Log" ADD CONSTRAINT "Log_userId_fkey" FOREIGN KEY ("userId") REFERENCES "User"("id") ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "User" ADD CONSTRAINT "User_roleId_fkey" FOREIGN KEY ("roleId") REFERENCES "Role"("id") ON DELETE SET NULL ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "Product" ADD CONSTRAINT "Product_typeId_fkey" FOREIGN KEY ("typeId") REFERENCES "SubcategotyType"("id") ON DELETE SET NULL ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "ProductPromotion" ADD CONSTRAINT "ProductPromotion_productId_fkey" FOREIGN KEY ("productId") REFERENCES "Product"("id") ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "ProductPromotion" ADD CONSTRAINT "ProductPromotion_promotionId_fkey" FOREIGN KEY ("promotionId") REFERENCES "Promotion"("id") ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "ProductPromotion" ADD CONSTRAINT "ProductPromotion_userId_fkey" FOREIGN KEY ("userId") REFERENCES "User"("id") ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "TypeField" ADD CONSTRAINT "TypeField_typeId_fkey" FOREIGN KEY ("typeId") REFERENCES "SubcategotyType"("id") ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "ProductFieldValue" ADD CONSTRAINT "ProductFieldValue_fieldId_fkey" FOREIGN KEY ("fieldId") REFERENCES "TypeField"("id") ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "ProductFieldValue" ADD CONSTRAINT "ProductFieldValue_productId_fkey" FOREIGN KEY ("productId") REFERENCES "Product"("id") ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "Message" ADD CONSTRAINT "Message_relatedProductId_fkey" FOREIGN KEY ("relatedProductId") REFERENCES "Product"("id") ON DELETE SET NULL ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "SupportTicket" ADD CONSTRAINT "SupportTicket_moderatorId_fkey" FOREIGN KEY ("moderatorId") REFERENCES "User"("id") ON DELETE SET NULL ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "SupportMessage" ADD CONSTRAINT "SupportMessage_authorId_fkey" FOREIGN KEY ("authorId") REFERENCES "User"("id") ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "Payment" ADD CONSTRAINT "Payment_userId_fkey" FOREIGN KEY ("userId") REFERENCES "User"("id") ON DELETE CASCADE ON UPDATE CASCADE;
