# Task App

Учебное fullstack-приложение для управления задачами, развёрнутое в Kubernetes.

## Стек

| Слой | Технология |
|------|-----------|
| Backend | Go (net/http), порт 5000 |
| Frontend | Node.js + Express + Axios, порт 3000 |
| Оркестрация | Kubernetes (namespace `tasks`) |

## Структура проекта

```
task-app/
├── backend/        # Go REST API
├── frontend/       # Node.js + Express
└── k8s/            # Kubernetes манифесты
    ├── namespace.yaml
    ├── backend-deployment.yaml
    ├── backend-service.yaml
    ├── frontend-deployment.yaml
    └── frontend-service.yaml
```

## API

| Метод | Endpoint | Описание |
|-------|----------|----------|
| GET | `/api/tasks` | Получить все задачи |
| POST | `/api/tasks` | Создать задачу `{"title": "..."}` |
| PUT | `/api/tasks/:id` | Переключить статус done |
| DELETE | `/api/tasks/:id` | Удалить задачу |
| GET | `/health` | Healthcheck |

## Запуск в Kubernetes

```bash
# Применить все манифесты
kubectl apply -f k8s/

# Проверить статус
kubectl get all -n tasks
```

## Локальный запуск

```bash
# Backend
cd backend && go run main.go

# Frontend
cd frontend && npm install && npm start
```
