# syntax=docker/dockerfile:1

############################
# Base deps
############################
FROM node:20-alpine AS base
WORKDIR /app
ENV NODE_ENV=production

############################
# Dependencies
############################
FROM base AS deps

RUN apk add --no-cache libc6-compat

COPY package.json ./

RUN npm install

############################
# Build
############################
FROM base AS builder

COPY --from=deps /app/node_modules ./node_modules
COPY . .

RUN npm run build

############################
# Runtime
############################
FROM base AS runner

WORKDIR /app

ENV NODE_ENV=production
ENV PORT=3000

# Create non-root user
RUN addgroup -g 1001 -S nodejs && adduser -S nextjs -u 1001

COPY --from=builder /app/next.config.mjs ./
COPY --from=builder /app/public ./public
COPY --from=builder /app/.next ./.next
COPY --from=builder /app/package.json ./package.json

USER nextjs

EXPOSE 3000

HEALTHCHECK --interval=15s --timeout=5s --retries=5 CMD wget -qO- http://localhost:3000/ || exit 1

CMD ["npm", "start"]
