# Build stage
FROM node:20-alpine AS builder

WORKDIR /app

# Enable corepack for yarn
RUN corepack enable

# Copy package files
COPY package.json ./

# Install dependencies with yarn (without lockfile requirement)
RUN yarn set version stable && \
    yarn install

# Copy source code
COPY . .

# Build the application
RUN yarn build

# Production stage with Caddy
FROM caddy:2.7-alpine

# Copy built files
COPY --from=builder /app/dist /usr/share/caddy

# Copy Caddyfile
COPY Caddyfile /etc/caddy/Caddyfile

# Expose port
EXPOSE 3000

# Caddy will start automatically
