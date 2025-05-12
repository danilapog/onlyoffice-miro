FROM node:20-alpine AS builder
WORKDIR /app

ARG VITE_MIRO_ONLYOFFICE_BACKEND
ENV VITE_MIRO_ONLYOFFICE_BACKEND=${VITE_MIRO_ONLYOFFICE_BACKEND}

COPY package*.json ./
COPY tsconfig.json ./
COPY vite.config.mts ./
COPY index.html ./
COPY application.html ./
COPY eslint.config.js ./

RUN npm install
COPY src ./src
RUN npm run build

FROM nginx:alpine
COPY --from=builder /app/dist /usr/share/nginx/html
EXPOSE 80

CMD ["nginx", "-g", "daemon off;"] 