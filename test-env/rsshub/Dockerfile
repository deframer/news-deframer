# Use a Node image that matches RSSHub's requirement
FROM node:24-bookworm-slim

RUN apt-get update \
    && apt-get install -y git \
    && rm -rf /var/lib/apt/lists/* \
    && npm install -g pnpm

WORKDIR /app

RUN git clone --depth 1 https://github.com/DIYgod/RSSHub.git .

RUN pnpm install

COPY custom_routes/ /app/lib/routes/

RUN pnpm build

EXPOSE 1200

CMD ["npm", "run", "start"]
