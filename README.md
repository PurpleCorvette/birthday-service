# Birthday Service

## Описание

`Birthday Service` - это сервис для удобного поздравления сотрудников с днем рождения. Основные возможности включают:

- Регистрацию и авторизацию пользователей.
- Получение списка сотрудников из различных источников (API, LDAP, прямая регистрация).
- Подписку и отписку от оповещений о днях рождения.
- Оповещения о днях рождения для подписанных пользователей.
- Внешнее взаимодействие через JSON API или телеграм-бота.
- Настройку времени оповещений до дня рождения на почту (при использовании фронтенда).

## Основные функции

1. **Авторизация**
    - Регистрация и авторизация пользователей с использованием JWT токенов.

2. **Управление сотрудниками**
    - Добавление, получение, редактирование и удаление сотрудников.
    - Интеграция с внешними источниками данных.

3. **Подписка на оповещения**
    - Возможность подписаться или отписаться от оповещений о днях рождения сотрудников.

4. **Оповещения**
    - Оповещения о днях рождения через различные каналы (телеграм-бот, электронная почта).

5. **Внешнее взаимодействие**
    - JSON API для взаимодействия с другими сервисами.
    - (Опционально) Телеграм-бот для оповещений.

## Установка и запуск

### Клонирование репозитория