# 🏢 Reservify - Sistema de Reservaciones

Sistema completo de gestión de reservaciones desarrollado con **Go (Gin)**, **Angular** y **MySQL**.

## 🚀 Características

- ✅ Autenticación y autorización con JWT
- 📅 Sistema de reservaciones con calendario interactivo
- 🏢 Gestión de recursos y servicios
- 👥 Panel de administración
- 🔔 Sistema de notificaciones
- 📊 Dashboard con estadísticas
- 🎨 Interfaz responsiva con Angular Material

## 🛠️ Tecnologías

### Backend
- **Go 1.21+**
- **Gin** - Framework web
- **GORM** - ORM para MySQL
- **JWT** - Autenticación
- **MySQL 8.0** - Base de datos

### Frontend
- **Angular 17+**
- **Angular Material** - UI Components
- **FullCalendar** - Calendario interactivo
- **Chart.js** - Gráficas y estadísticas
- **TypeScript**

## 📁 Estructura del Proyecto

```
reservify/
├── backend/          # API REST en Go
│   ├── cmd/          # Punto de entrada
│   ├── config/       # Configuración
│   ├── models/       # Modelos y migraciones
│   ├── controllers/  # Controladores MVC
│   ├── services/     # Lógica de negocio
│   ├── repositories/ # Acceso a datos
│   ├── middleware/   # Middlewares
│   ├── routes/       # Definición de rutas
│   ├── dto/          # Data Transfer Objects
│   └── utils/        # Utilidades
├── frontend/         # Aplicación Angular
└── docker-compose.yml
```

## 🚀 Instalación y Uso

### Prerrequisitos
- Go 1.21 o superior
- Node.js 18+ y npm
- MySQL 8.0

### Instalación de MySQL

1. **Descargar MySQL:**
   - https://dev.mysql.com/downloads/installer/
   - Instalar "Developer Default"

2. **Crear base de datos:**
   ```sql
   CREATE DATABASE reservify_db;
   CREATE USER 'reservify'@'localhost' IDENTIFIED BY 'reservify123';
   GRANT ALL PRIVILEGES ON reservify_db.* TO 'reservify'@'localhost';
   FLUSH PRIVILEGES;
   ```

### Instalación del Proyecto

#### Backend
```bash
cd backend

# Crear archivo .env
cp .env.example .env

# Instalar dependencias de Go
go mod download

# Ejecutar el servidor (las migraciones se ejecutan automáticamente)
go run cmd/api/main.go
```

El backend estará corriendo en `http://localhost:8080`
Las tablas se crearán automáticamente en la primera ejecución.

#### Frontend
```bash
cd frontend

# Instalar dependencias
npm install

# Ejecutar en modo desarrollo
ng serve
```

El frontend estará disponible en `http://localhost:4200`

## 📋 Roadmap de Desarrollo

- [x] Fase 0: Setup inicial
- [x] Fase 1: Base de datos y modelos
- [x] Fase 2: Autenticación y autorización
- [ ] Fase 3: Gestión de usuarios
- [ ] Fase 4: Gestión de recursos
- [ ] Fase 5: Sistema de reservaciones
- [ ] Fase 6: Notificaciones
- [ ] Fase 7: Dashboard y estadísticas
- [ ] Fase 8: Testing
- [ ] Fase 9: UI/UX mejoras
- [ ] Fase 10: Deployment

## 🤝 Contribuciones

Este es un proyecto de aprendizaje personal, pero si quieres contribuir:

1. Fork el proyecto
2. Crea una rama para tu feature (`git checkout -b feature/AmazingFeature`)
3. Commit tus cambios (`git commit -m 'Add some AmazingFeature'`)
4. Push a la rama (`git push origin feature/AmazingFeature`)
5. Abre un Pull Request

## 📝 Licencia

Este proyecto está bajo la Licencia MIT - ver el archivo [LICENSE](LICENSE) para más detalles.

## 👤 Autor

**Stormdead**
- GitHub: [@Stormdead](https://github.com/Stormdead)
- LinkedIn: [Tu perfil de LinkedIn]

## 📸 Screenshots

_(Aquí agregarás screenshots cuando tengas la aplicación funcionando)_

---

⭐️ Si este proyecto te ayudó, dale una estrella en GitHub!