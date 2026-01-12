FROM node:24-alpine AS builder

WORKDIR /app

COPY package.json yarn.lock ./
RUN yarn install --frozen-lockfile

COPY prisma ./prisma
RUN yarn prisma generate

COPY . .
RUN yarn build

FROM node:24-alpine

WORKDIR /app

COPY package.json yarn.lock ./
RUN yarn install --production --frozen-lockfile

COPY --from=builder /app/prisma ./prisma
COPY --from=builder /app/dist ./dist
COPY --from=builder /app/prisma.config.ts ./prisma.config.ts

EXPOSE 3000

CMD ["sh", "-c", "sleep 5 && yarn prisma db push --skip-generate --accept-data-loss && node dist/src/main"]