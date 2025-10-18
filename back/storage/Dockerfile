FROM node:22-alpine AS builder

WORKDIR /app

COPY package*.json ./
RUN npm install

COPY . .

RUN apk add aws-cli

# aws.shを実行
RUN chmod +x ./aws.sh && \
    echo "=== Running aws.sh ===" && \
    ./aws.sh && \
    echo "=== AWS config after aws.sh ===" && \
    ls -la ~/.aws/ || echo "No .aws directory found"

RUN npm run build

FROM node:22-alpine

WORKDIR /app

COPY package*.json ./
RUN npm install --production

# ビルド環境から、コンパイル済みのJavaScriptファイルだけをコピー
COPY --from=builder /app/dist ./dist

# アプリケーションがリッスンするポート
EXPOSE 3001

# コンテナ起動時に、ビルド済みのJSファイルを実行
CMD [ "node", "/app/dist/bin/www.js" ]
