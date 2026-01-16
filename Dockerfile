FROM node:24-alpine AS builder

WORKDIR /app

COPY package.json yarn.lock ./
RUN yarn install --frozen-lockfile

COPY . .

RUN yarn prisma generate

RUN yarn build

FROM node:24-alpine

WORKDIR /app

# Устанавливаем postgresql-client для psql
RUN apk add --no-cache postgresql-client

COPY package.json yarn.lock ./
RUN yarn install --production --frozen-lockfile

COPY --from=builder /app/dist ./dist
COPY --from=builder /app/prisma ./prisma
COPY --from=builder /app/prisma.config.ts ./
COPY --from=builder /app/dumps ./dumps
COPY start.sh ./

RUN chmod +x start.sh

EXPOSE 3000

CMD ["sh", "./start.sh"]