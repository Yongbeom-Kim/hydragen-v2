FROM node:20-alpine AS base
WORKDIR /app
COPY package.json pnpm-lock.yaml* ./
RUN corepack enable
RUN pnpm install --frozen-lockfile
COPY . .
EXPOSE 3000

FROM base AS dev
CMD ["pnpm", "run", "dev", "--", "--host", "0.0.0.0", "--port", "3000"]

FROM base AS build
RUN pnpm run build

FROM node:20-alpine AS runtime
WORKDIR /app
COPY --from=build /app/package.json /app/pnpm-lock.yaml* ./
RUN corepack enable 
RUN pnpm install --prod --frozen-lockfile

COPY --from=build /app/.output ./.output

ENV NODE_ENV=production
ENV NITRO_HOST=0.0.0.0
ENV NITRO_PORT=3000

EXPOSE 3000

CMD ["node", ".output/server/index.mjs"]

