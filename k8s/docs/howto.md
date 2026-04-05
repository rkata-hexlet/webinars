Ниже представлен полный сценарий лабораторной работы в формате **Markdown** (который легко конвертируется в PDF). Инструкция охватывает всё: от установки на Linux до глубокой отладки состояний.

---

# Лабораторная работа: Развертывание StreamFlow в Kubernetes (Minikube)

**Цель:** Научиться устанавливать локальный кластер, разворачивать микросервисную архитектуру, управлять ресурсами и отлаживать ошибки.

---

## 1. Подготовка окружения (Linux)

### Установка Minikube и kubectl
Для работы на Ubuntu/Debian выполните следующие команды:

```bash
# 1. Установка kubectl (пульт управления)
curl -LO "https://dl.k8s.io/release/$(curl -L -s https://dl.k8s.io/release/stable.txt)/bin/linux/amd64/kubectl"
sudo install -o root -g root -m 0755 kubectl /usr/local/bin/kubectl

# 2. Установка Minikube
curl -LO https://storage.googleapis.com/minikube/releases/latest/minikube-linux-amd64
sudo install minikube-linux-amd64 /usr/local/bin/minikube

# 3. Запуск (требуется установленный Docker)
minikube start --driver=docker
```

### Проверка связи
```bash
kubectl get nodes
# Вывод должен показать одну ноду minikube в статусе Ready
```

---

## 2. Сборка "Универсального" образа
Мы используем один образ для всех ролей (API, Catalog, Transcoder).



```bash
# Направляем Docker в контекст Minikube
eval $(minikube docker-env)

# Создаем main.go (код из предыдущего этапа) и собираем:
docker build -t streamflow-mock:v1 .
```

---

## 3. Развертывание базы (Манифест 6 элементов)

Создайте файл `streamflow.yaml`. В нем мы пропишем **Resources**, **Probes** и **ConfigMaps**.

```yaml
# --- КОНФИГУРАЦИЯ ---
apiVersion: v1
kind: ConfigMap
metadata:
  name: stream-config
data:
  WELCOME_MSG: "StreamFlow v1.0 Production"
---
# --- CATALOG (Backend) ---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: catalog-deploy
spec:
  replicas: 2
  selector:
    matchLabels:
      app: catalog
  template:
    metadata:
      labels:
        app: catalog
    spec:
      containers:
      - name: server
        image: streamflow-mock:v1
        env:
        - name: SERVICE_NAME
          value: "catalog"
        # ТЕМА: РЕСУРСЫ
        resources:
          requests:
            memory: "64Mi"
            cpu: "100m"
          limits:
            memory: "128Mi"
            cpu: "200m"
        # ТЕМА: ПРОВЕРКА ЗДОРОВЬЯ
        livenessProbe:
          httpGet:
            path: /healthz
            port: 8080
          initialDelaySeconds: 3
          periodSeconds: 5
---
apiVersion: v1
kind: Service
metadata:
  name: catalog-service
spec:
  selector:
    app: catalog
  ports:
  - port: 80
    targetPort: 8080
```
*(Аналогично добавьте блоки для api-deploy и transcoder-deploy)*

---

## 4. Эксплуатация и отладка (Практические задачи)

### Задача 1: Состояние Пода и Ошибки старта
Попробуйте изменить образ в манифесте на несуществующий: `image: streamflow-mock:v999`.
1. Примените: `kubectl apply -f streamflow.yaml`.
2. Посмотрите статус: `kubectl get pods`.
   * **Результат:** Вы увидите `ImagePullBackOff`.
3. Исправьте обратно и уроните сервис через код (эндпоинт `/fail`).
   * **Результат:** Вы увидите `CrashLoopBackOff`.
   * **Как лечить:** `kubectl logs <pod_name> --previous` — это покажет причину падения перед рестартом.

### Задача 2: Работа с логами
Логи — это единственный способ понять, что происходит в распределенной системе.
* **Следить в реальном времени:** `kubectl logs -f -l app=catalog`
* **Посмотреть логи конкретного контейнера:** `kubectl logs <pod_name>`

### Задача 3: Рескейлинг (Масштабирование)
Представьте, что нагрузка выросла.
```bash
# Увеличиваем количество каталогов до 5
kubectl scale deployment catalog-deploy --replicas=5

# Наблюдаем за процессом
kubectl get pods -w
```

### Задача 4: Проверка здоровья (Liveness Probe)
Если ваш сервис зависнет (Deadlock), Liveness Probe это заметит.
* Проверьте описание пода: `kubectl describe pod <catalog_pod_name>`.
* Найдите секцию `Events`. Там будет история проверок.



---

## 5. Полезная шпаргалка (Cheat Sheet)

| Команда | Что делает |
| :--- | :--- |
| `kubectl get all` | Показать всё: поды, сервисы, деплойменты |
| `kubectl describe pod <name>` | Найти причину, почему под не стартует (Events) |
| `kubectl exec -it <name> -- sh` | "Зайти" внутрь пода (как по SSH) |
| `minikube dashboard` | Открыть красивый UI в браузере |
| `kubectl rollout restart deploy <name>` | Принудительно перезапустить все поды |

---

**Задание для закрепления:**
1. Разверните систему.
2. Добейтесь, чтобы через `minikube tunnel` открылась страница `localhost`.
3. Убейте один под каталога и убедитесь, что через 5 секунд появится новый с тем же функционалом.