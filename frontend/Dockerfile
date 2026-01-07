# Stage 1: Generate API client (needs backend swagger)
FROM golang:1.24-alpine AS swagger

WORKDIR /app
COPY backend/go.mod backend/go.sum* ./backend/
RUN cd backend && go mod download

COPY backend/ ./backend/
RUN cd backend && go run github.com/swaggo/swag/cmd/swag@latest init -g cmd/server/main.go -o docs --parseDependency --st

# Stage 2: Build frontend
FROM node:20-alpine AS builder

WORKDIR /app

# Copy package files
COPY frontend/package.json frontend/package-lock.json* ./

# Install dependencies
RUN npm ci

# Copy swagger.json from swagger stage
COPY --from=swagger /app/backend/docs/swagger.json ./swagger.json

# Copy source code
COPY frontend/ .

# Update orval config to use local swagger.json
RUN sed -i "s|../backend/docs/swagger.json|./swagger.json|g" orval.config.ts

# Generate API client
RUN npm run generate:api

# Generate route tree (needed before tsc)
RUN npx @tanstack/router-cli generate

# Build the application
RUN npm run build

# Stage 3: Production
FROM nginx:alpine

# Copy custom nginx config
COPY --from=builder /app/nginx.conf /etc/nginx/conf.d/default.conf

# Copy built assets
COPY --from=builder /app/dist /usr/share/nginx/html

# Expose port
EXPOSE 80

# Start nginx
CMD ["nginx", "-g", "daemon off;"]
