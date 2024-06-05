# Birthday Notification Service

## Описание

Birthday Notification Service - это сервис для удобного поздравления сотрудников с днем рождения. Сервис предоставляет следующие возможности:
- Регистрация пользователей
- Получение списка сотрудников
- Подписка/отписка от оповещений о дне рождения
- Оповещение о дне рождения того, на кого подписан пользователь
- Внешнее взаимодействие через Telegram бота

## Установка и запуск

### Шаг 1: Клонирование репозитория

```sh
git clone https://github.com/PurpleCorvette/birthday-service.git
cd birthday-service 
```

### Шаг 2: Настройка базы данных PostgreSQL
```sql
CREATE DATABASE yourdatabase;
CREATE USER yourusername WITH ENCRYPTED PASSWORD 'yourpassword';
GRANT ALL PRIVILEGES ON DATABASE yourdatabase TO yourusername;

CREATE TABLE users (
                       id SERIAL PRIMARY KEY,
                       username VARCHAR(255) UNIQUE NOT NULL,
                       password_hash VARCHAR(255) NOT NULL,
                       telegram_id BIGINT UNIQUE
);

CREATE TABLE employees (
                           id SERIAL PRIMARY KEY,
                           name VARCHAR(255) NOT NULL,
                           birthday DATE NOT NULL
);

CREATE TABLE subscriptions (
                               user_id INT REFERENCES users(id),
                               employee_id INT REFERENCES employees(id),
                               PRIMARY KEY (user_id, employee_id)
);
```

### Шаг 3: Настройка переменных окружения
Создайте файл .env в корне проекта и укажите в нем необходимые переменные окружения:
```.env
DATABASE_URL=postgres://username:password@localhost:5432/yourdatabase
TELEGRAM_BOT_TOKEN=your-telegram-bot-token
TELEGRAM_CHAT_ID=your-chat-id
```

### Шаг 4: Запуск с использованием Docker
```sh
docker build -t birthday-service .
docker run --env-file .env -p 8080:8080 birthday-service
```

### Использование
Регистрация пользователя через Telegram
#### Откройте Telegram и найдите вашего бота.
#### Отправьте команду /start, чтобы получить приветственное сообщение.
#### Отправьте команду /register <username> <password>, чтобы зарегистрироваться.

### Примеры:

#### Добавление сотрудника:
curl -X POST http://localhost:8080/employee -H "Content-Type: application/json" -d '{"name":"Вася Пупкин", "birthday":"1980-01-01"}'

#### Получение списка всех сотрудников:
curl -X GET http://localhost:8080/employees

#### Подписка на уведомления через Telegram:
curl -X GET http://localhost:8080/employees

#### Отправка уведомлений в групповой чат Telegram:
curl -X POST http://localhost:8080/trigger-notifications

#### Дополнительные команды:
##### Обновление сотрудника:
curl -X PUT http://localhost:8080/employee/1 -H "Content-Type: application/json" -d '{"name":"Вася Пупкин Updated", "birthday":"1980-01-02"}'
##### Удаление сотрудника: 
curl -X DELETE http://localhost:8080/employee/1