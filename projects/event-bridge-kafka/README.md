# 🌉 Event Bridge Kafka

> **Gateway que publica eventos blockchain o IA en Kafka, con métricas Prometheus.**

## 📋 **Descripción**

Event Bridge Kafka es el punto de entrada del ecosistema StreamForge. Se encarga de recibir eventos de fuentes externas (blockchain, IA, IoT) y publicarlos en Kafka para su procesamiento posterior.

## 🛠️ **Stack Tecnológico**

- **Backend**: NestJS + TypeScript
- **Messaging**: KafkaJS
- **Observabilidad**: Prometheus + Grafana
- **Containerización**: Docker
- **Base de datos**: PostgreSQL

## 🚀 **Características**

- ✅ Recepción de eventos de múltiples fuentes
- ✅ Validación y transformación de datos
- ✅ Publicación en Kafka con particionado inteligente
- ✅ Métricas Prometheus integradas
- ✅ Health checks y logging estructurado
- ✅ Rate limiting y autenticación
- ✅ Retry automático y dead letter queue

## 📊 **Arquitectura**

```mermaid
graph LR
    A[Blockchain APIs] --> D[Event Bridge]
    B[AI Models] --> D
    C[IoT Sensors] --> D
    D --> E[Kafka Topics]
    D --> F[Prometheus Metrics]
    D --> G[PostgreSQL Logs]
```

## 🔧 **Configuración**

### **Variables de Entorno**

```bash
# Kafka
KAFKA_BROKERS=localhost:9092
KAFKA_TOPIC_EVENTS=blockchain-events
KAFKA_TOPIC_AI=ai-events
KAFKA_TOPIC_IOT=iot-events

# Base de datos
DATABASE_URL=postgresql://user:pass@localhost:5432/streamforge

# Observabilidad
PROMETHEUS_PORT=9091
LOG_LEVEL=info

# Autenticación
JWT_SECRET=your-secret-key
API_RATE_LIMIT=1000
```

### **Endpoints de API**

```yaml
POST /api/v1/events/blockchain
POST /api/v1/events/ai
POST /api/v1/events/iot
GET  /api/v1/health
GET  /api/v1/metrics
```

## 🚀 **Inicio Rápido**

```bash
# Instalar dependencias
npm install

# Configurar variables de entorno
cp .env.example .env

# Levantar en desarrollo
make up

# Ver logs
make logs

# Ejecutar tests
make test
```

## 📈 **Métricas Prometheus**

- `event_bridge_events_received_total`
- `event_bridge_events_published_total`
- `event_bridge_events_failed_total`
- `event_bridge_processing_duration_seconds`
- `event_bridge_kafka_lag_seconds`

## 🧪 **Testing**

```bash
# Tests unitarios
npm run test

# Tests de integración
npm run test:integration

# Tests E2E
npm run test:e2e

# Coverage
npm run test:coverage
```

## 📚 **API Documentation**

### **Publicar Evento Blockchain**

```bash
curl -X POST http://localhost:3001/api/v1/events/blockchain \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d '{
    "blockNumber": 12345,
    "transactionHash": "0x...",
    "from": "0x...",
    "to": "0x...",
    "value": "1000000000000000000",
    "timestamp": 1640995200
  }'
```

### **Publicar Evento IA**

```bash
curl -X POST http://localhost:3001/api/v1/events/ai \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d '{
    "modelId": "gpt-4",
    "input": "Hello world",
    "output": "Hello! How can I help you?",
    "confidence": 0.95,
    "timestamp": 1640995200
  }'
```

## 🔍 **Monitoreo**

### **Health Check**

```bash
curl http://localhost:3001/api/v1/health
```

### **Métricas**

```bash
curl http://localhost:3001/api/v1/metrics
```

## 🐳 **Docker**

```bash
# Construir imagen
docker build -t streamforge/event-bridge-kafka .

# Ejecutar contenedor
docker run -p 3001:3000 \
  -e KAFKA_BROKERS=localhost:9092 \
  -e DATABASE_URL=postgresql://user:pass@localhost:5432/streamforge \
  streamforge/event-bridge-kafka
```

## 📁 **Estructura del Proyecto**

```
event-bridge-kafka/
├── src/
│   ├── controllers/        # Controladores de API
│   ├── services/          # Lógica de negocio
│   ├── entities/          # Entidades de base de datos
│   ├── dto/              # Data Transfer Objects
│   ├── guards/           # Guards de autenticación
│   ├── interceptors/      # Interceptores
│   └── utils/            # Utilidades
├── test/                 # Tests
├── docker/              # Configuración Docker
├── docs/               # Documentación
└── scripts/            # Scripts de automatización
```

## 🤝 **Contribuir**

1. Fork el proyecto
2. Crea tu feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit tus cambios (`git commit -m 'Add some AmazingFeature'`)
4. Push a la branch (`git push origin feature/AmazingFeature`)
5. Abre un Pull Request

## 📄 **Licencia**

Este proyecto está bajo la Licencia MIT - ver el archivo [LICENSE](LICENSE) para detalles.

---

**Parte del ecosistema StreamForge** 🚀
