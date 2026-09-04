# Think Skeleton Application

Official starter skeleton application for **[Think Framework](https://github.com/go-think/think)**.

---

## 🚀 Quick Start

### 1. Requirements
- Go 1.23+

### 2. Run Locally
```bash
# Clone or navigate to directory
cd think-app

# Run application
go run main.go
```

The application will start on `http://localhost:8080`.

---

## 📂 Directory Structure

- `app/http/controllers/` - HTTP Controllers
- `app/http/middleware/` - Application Middlewares
- `app/http/kernel.go` - HTTP Kernel (Global stack, Groups, Aliases)
- `app/models/` - Data & Domain Models
- `app/providers/` - Service Providers (App, Event, Route)
- `bootstrap/app.go` - Framework instance bootstrapper
- `config/` - Configuration files
- `resources/views/` - HTML templates
- `routes/web.go` - Web routes (Session & Cookie enabled)
- `routes/api.go` - API routes (JSON responses with `/api` prefix)
- `storage/` - Cache, Session files, Views, and Logs

---

## 🧪 Available Endpoints

- `GET /` - Welcome page
- `GET /status` - Live server metadata & session visit counter
- `GET /api/ping` - JSON ping response
- `GET /api/user` - Protected endpoint example (`auth` middleware alias)
