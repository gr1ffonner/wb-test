# WB Test Application

A Go-based microservice application with PostgreSQL, Redis, and NATS messaging.

## 🚀 Quick Start

### Prerequisites
- Docker & Docker Compose
- Go 1.21+ (for local development)


# Server
SERVER_PORT=8080

# Redis
REDIS_ADDR=redis:6379
REDIS_PASSWORD=
REDIS_DB=0

# NATS
NATS_URL=nats://nats-streaming:4222

# Logging
LOG_LEVEL=info
```

### Running with Docker Compose

Start all services:
```bash
make up
```

Start only infrastructure (without app):
```bash
make up-dev
```

Stop all services:
```bash
make down
```

### Local Development

1. Install dependencies:
```bash
go mod tidy
```

2. Start only infrastructure (without app):
```bash
make up-dev
```

3. Run the application:
```bash
make run
```

## 📋 Services & Ports

| Service | Port | Description |
|---------|------|-------------|
| **App** | `8080` | Main Go application |
| **PostgreSQL** | `5432` | Database |
| **Redis** | `6379` | Cache |
| **NATS Streaming** | `4222` | Message broker |
| **NATS Monitoring** | `8222` | NATS monitoring UI |



## 📝 Makefile Commands

```bash
make up          # Start all services with test profile
make down        # Stop all services
make build       # Build Go application
make run         # Run application locally
make test        # Run tests
make clean       # Clean build artifacts
```


# ТЗ
```
WB Tech: level # 0 (Golang)		 	 	
Тестовое задание
Необходимо разработать демонстрационный сервис с простейшим интерфейсом, отображающий данные о заказе. 
Модель данных в формате JSON прилагается к заданию.	
				
Что нужно сделать:
Развернуть локально PostgreSQL
Создать свою БД
Настроить своего пользователя
Создать таблицы для хранения полученных данных
Разработать сервис
Реализовать подключение и подписку на канал в nats-streaming
Полученные данные записывать в БД
Реализовать кэширование полученных данных в сервисе с помощью редис
В случае падения сервиса необходимо восстанавливать кэш из БД
Запустить http-сервер и выдавать данные по id из кэша
Разработать простейший интерфейс отображения полученных данных по id заказа
Советы				
Данные статичны, исходя из этого подумайте насчет модели хранения в кэше и в PostgreSQL. Модель в файле model.json
Подумайте как избежать проблем, связанных с тем, что в канал могут закинуть что-угодно
Чтобы проверить работает ли подписка онлайн, сделайте себе отдельный скрипт, для публикации данных в канал
Подумайте как не терять данные в случае ошибок или проблем с сервисом
Nats-streaming разверните локально (не путать с Nats)
						
Бонус-задание						
Покройте сервис автотестами — будет плюсик вам в карму.
Устройте вашему сервису стресс-тест: выясните на что он способен.
						
Воспользуйтесь утилитами WRK и Vegeta, попробуйте оптимизировать код.
```



