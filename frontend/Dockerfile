# # Используем официальный образ Node.js
# FROM node:16 AS build

# # Устанавливаем рабочую директорию внутри контейнера
# WORKDIR /app

# # Копируем package.json и package-lock.json для установки зависимостей
# COPY package*.json ./

# # Устанавливаем зависимости
# RUN npm install

# # Копируем все файлы проекта в контейнер
# COPY . .

# # Собираем приложение (если в продакшн)
# RUN npm run build

# # Если в режиме разработки, запускаем сервер
# CMD ["npm", "start"]

# # Экспонируем порт 3000
# EXPOSE 3000

# 1. Используем Node.js для сборки React-приложения
FROM node:18 AS build

WORKDIR /app

COPY package*.json ./
RUN npm install

COPY . .
RUN npm run build

# 2. Используем Nginx для раздачи статики
FROM nginx:alpine

# Копируем собранные файлы React в папку Nginx
COPY --from=build /app/build /usr/share/nginx/html

# Открываем 80-й порт для доступа к Nginx
EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]
