# 🌐 URL Analyzer

A web application built with Golang and Gin that analyzes a given webpage URL and returns structured metadata such as:

- HTML version
- Page title
- Heading tag count (`h1` to `h6`)
- Number of internal, external, and inaccessible links
- Presence of a login form

---

## 🧭 Project Structure

```
url-analyzer/
├── cmd/                      # Entry point for the application (main.go)
├── config/                   # YAML configuration files
├── docs/                     # Swagger auto-generated API documentation
├── internal/
│   ├── config/               # Loads YAML config
│   ├── handler/              # Gin route handlers
│   ├── model/                # Request and response models
│   ├── service/              # Core HTML analysis logic
│   └── utils/                # Utility functions (URL validation, HTTP checking)
├── static/
│   ├── css/style.css         # UI styles
│   └── js/app.js             # JavaScript for form submission
├── templates/index.html      # User-facing HTML UI
├── tests/                    # External integration & unit tests
├── Dockerfile                # Docker build config
├── go.mod / go.sum           # Go module dependencies
└── README.md                 # You're here!
```

---

## 🚀 Getting Started

### ✅ 1. Clone the repository

```bash
git clone https://github.com/your-username/url-analyzer.git
cd url-analyzer
```

### ✅ 2. Run locally with Go

```bash
go mod tidy
go run ./cmd/main.go
```

Visit: [http://localhost:8080](http://localhost:8080)

---

## 🐳 Docker Deployment

### ✅ Build Docker Image

```bash
docker build -t url-analyzer .
```

### ✅ Run Container

```bash
docker run -p 8080:8080 url-analyzer
```

---

## 📘 API Documentation

Swagger docs:

```bash
swag init -g cmd/main.go --output docs
```

Docs: [http://localhost:8080/swagger/index.html](http://localhost:8080/swagger/index.html)

---

## 🧪 Run Tests

```bash
go test ./internal/... ./tests/... -v -cover
```

Generate coverage report:

```bash
go test ./internal/... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

---

## ⚙️ Configuration

```yaml
server:
  port: "8080"
```

---
