# 🔧 StreamForge Scripts

> **Scripts de automatización para el ecosistema StreamForge.**

## 📋 **Scripts Disponibles**

### **🚀 Desarrollo**
- `setup-dev.sh` - Configurar entorno de desarrollo
- `install-dependencies.sh` - Instalar dependencias
- `start-dev.sh` - Iniciar entorno de desarrollo
- `stop-dev.sh` - Parar entorno de desarrollo

### **🐳 Docker**
- `build-all.sh` - Construir todas las imágenes
- `clean-docker.sh` - Limpiar contenedores y volúmenes
- `docker-logs.sh` - Ver logs de todos los servicios
- `docker-stats.sh` - Estadísticas de contenedores

### **☸️ Kubernetes**
- `deploy-k8s.sh` - Desplegar en Kubernetes
- `undeploy-k8s.sh` - Eliminar despliegue
- `scale-services.sh` - Escalar servicios
- `k8s-logs.sh` - Ver logs de pods

### **📊 Monitoreo**
- `setup-monitoring.sh` - Configurar monitoreo
- `import-dashboards.sh` - Importar dashboards Grafana
- `export-metrics.sh` - Exportar métricas
- `health-check.sh` - Verificar salud del sistema

### **🔒 Seguridad**
- `setup-security.sh` - Configurar seguridad
- `generate-certs.sh` - Generar certificados
- `audit-logs.sh` - Analizar logs de auditoría
- `security-scan.sh` - Escanear vulnerabilidades

### **📈 Performance**
- `load-test.sh` - Pruebas de carga
- `benchmark.sh` - Benchmarks
- `profile-services.sh` - Profiling de servicios
- `optimize-resources.sh` - Optimizar recursos

### **🗄️ Base de Datos**
- `init-db.sh` - Inicializar base de datos
- `migrate-db.sh` - Ejecutar migraciones
- `backup-db.sh` - Backup de base de datos
- `restore-db.sh` - Restaurar base de datos

### **🔄 CI/CD**
- `ci-setup.sh` - Configurar CI/CD
- `run-tests.sh` - Ejecutar tests
- `build-release.sh` - Construir release
- `deploy-staging.sh` - Desplegar en staging
- `deploy-production.sh` - Desplegar en producción

## 🚀 **Uso de Scripts**

### **Configurar Entorno de Desarrollo**
```bash
# Hacer ejecutable
chmod +x scripts/setup-dev.sh

# Ejecutar
./scripts/setup-dev.sh
```

### **Construir Todas las Imágenes**
```bash
# Construir imágenes
./scripts/build-all.sh

# Con limpieza previa
./scripts/build-all.sh --clean
```

### **Desplegar en Kubernetes**
```bash
# Desplegar
./scripts/deploy-k8s.sh

# Con configuración específica
./scripts/deploy-k8s.sh --env=production --replicas=3
```

### **Configurar Monitoreo**
```bash
# Configurar monitoreo completo
./scripts/setup-monitoring.sh

# Solo Prometheus
./scripts/setup-monitoring.sh --prometheus-only
```

## 🔧 **Configuración de Scripts**

### **Variables de Entorno**
```bash
# Configurar variables
export STREAMFORGE_ENV=development
export KUBECONFIG=/path/to/kubeconfig
export DOCKER_REGISTRY=your-registry.com
export GRAFANA_ADMIN_PASSWORD=secure-password
```

### **Parámetros Comunes**
```bash
# Parámetros disponibles
--env=development|staging|production
--replicas=3
--clean
--verbose
--dry-run
```

## 📁 **Estructura de Scripts**

```
scripts/
├── development/          # Scripts de desarrollo
├── docker/              # Scripts de Docker
├── kubernetes/          # Scripts de K8s
├── monitoring/          # Scripts de monitoreo
├── security/            # Scripts de seguridad
├── performance/         # Scripts de performance
├── database/            # Scripts de base de datos
├── cicd/                # Scripts de CI/CD
└── utils/               # Utilidades comunes
```

## 🤝 **Contribuir**

1. Fork el proyecto
2. Crea tu feature branch (`git checkout -b feature/AmazingScript`)
3. Commit tus cambios (`git commit -m 'Add some AmazingScript'`)
4. Push a la branch (`git push origin feature/AmazingScript`)
5. Abre un Pull Request

## 📄 **Licencia**

Este proyecto está bajo la Licencia MIT - ver el archivo [LICENSE](LICENSE) para detalles.

---

**Scripts del ecosistema StreamForge** 🚀
