# Song library

Приложение реализует онлайн библиотеку песен.

Приложение состоит из:
- Backend API - REST API для управлением данными через HTTP запросы;
- External API - внешний API с обогащенной информацией (не реализован в рамках данного задания)
- База данных - PostgreSQL.

Для сервиса создан Dockerfile и файл compose, которые собирают образы и запускают сервис и БД.

Необходимо реализовать следующее:
1. Выставить REST методы:
- получение данных библиотеки с фильтрацией по всем полям и пагинацией;
- получение текста песни с пагинацией по куплетам;
- удаление песни;
- изменение данных песни;
- добавление новой песни в формате JSON
```JSON
  {
  "group": "Muse",
  "song": "Supermassive Black Hole"
  }
```
2. При добавлении сделать запрос в API, описанного сваггером. API, описанный сваггером, 
реализовывать отдельно не нужно.
3. Обогащенную информацию положить в БД PostgreSQL (структура БД должна быть создана путем 
миграций при старте сервиса).
4. Покрыть код debug- и info-логами.
5. Вынести конфигурационные данные в .env файл.
6. Сгенерировать сваггер на реализованное API.

### Запуск

Для запуска приложения необходимо:
1. Склонировать репозиторий
```bash
git clone https://github.com/DeaPacis/song_library.git
```
2. Перейти в директорию
```bash
cd song_library
```
3. Создать .env файл на основе .env.example файла
```bash
cp .env.example .env
```
4. Запустить compose файл
```bash
docker compose up -d
```
5.Документация проекта доступна по ссылке http://localhost:8080/swagger/index.html