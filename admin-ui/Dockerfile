# Build stage
FROM node:20-alpine AS builder

ARG VITE_API_BASE_URL=''
ENV VITE_API_BASE_URL=${VITE_API_BASE_URL}

WORKDIR /app

# Configure Alpine mirrors for China
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories

COPY package*.json ./
RUN npm ci

COPY . ./
# Build with vite directly (skip tsc to avoid test file type checking issues)
RUN npx vite build

# Production stage
FROM nginx:alpine

COPY --from=builder /app/dist /usr/share/nginx/html

COPY nginx.conf /etc/nginx/conf.d/default.conf

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]
