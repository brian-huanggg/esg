# --- Build stage: full Node toolchain, discarded after producing dist/ ---
FROM node:24-alpine AS build
WORKDIR /app

COPY package.json package-lock.json ./
RUN npm ci

COPY . .
RUN npm run build

# --- Runtime stage: static files + Nginx only, runs as non-root (listens on 8080) ---
FROM nginxinc/nginx-unprivileged:1.29-alpine

# Files in /etc/nginx/templates/*.template are run through envsubst at container
# start and written to /etc/nginx/conf.d/, so the backend address is injected per environment.
COPY nginx.conf /etc/nginx/templates/default.conf.template
COPY --from=build /app/dist /usr/share/nginx/html

# Default for local docker run; compose / Helm override it with the real service name.
ENV BACKEND_URL=http://host.docker.internal:8080
EXPOSE 8080
