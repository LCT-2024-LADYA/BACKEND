# DEPLOY

## CLONE
- Через https:
    ```bash
    git clone https://github.com/LCT-2024-LADYA/BACKEND.git
    ```
- Через ssh:
    ```bash
    git clone git@github.com:LCT-2024-LADYA/BACKEND.git
    ```

## ENVIRONMENT
Необходимо настроить файл виртуального окружения. Для этого создайте
в корне клонированного проекта файл `.env` и скопируйте в него содержимое
файла `.env.example`. Заполните файл в соответствии с конфигурацией вашего проекта (возможный пример заполнения будет лежать на Яндекс Диске в `data/.env`).

## START

### Загрузка фото

Перенесите фотографии тренировок с Яндекс Диска (`data/photos`) в папку `/static/img/exercises`

### Запуск проект
Запустите проект (из корня клонированного репозитория):
```bash
docker compose up --build
```

### Загрузка данных

Это можно сделать через загрузку `dump.sql` или же через swagger.
Откройте swagger (http://localhost:8080/swagger/index.html) и откройте описание ручки `Admin Authorization` (1 в списке).
В параметр `X-API-KEY` передайте `API_KEY` из файла `.env`. Полученный `access_token` необходимо вставить в параметры ручек `Create Exercises` (POST http://localhost:8080/api/training/exercise),
`Create Training Base` (POST http://localhost:8080/api/training/base) для загрузки в базу упражнений и базовых тренировок. Данные для вставки body лежат в папке `data/exercises.json` и `data/trainings.json` на переданном вам Яндекс Диске.

### Регистрация тренера

Регистрация тренера осуществляется программно, для этого имеется ручка `Trainer Register` (5 в списке). Туда необходимо передать `access_token` администратора, почту и пароль будущего тренера.
