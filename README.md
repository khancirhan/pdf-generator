# PDF Generator

A Go-based REST API service for rendering HTML templates and generating PDFs using the Liquid templating engine and Gotenberg.

## Features

- Render HTML templates with dynamic data using Liquid syntax
- Generate PDFs from templates with customizable options (paper size, margins, orientation)
- Template management API
- Docker-based deployment with Gotenberg integration

## Requirements

- Go 1.25+
- Docker & Docker Compose

## Quick Start

```bash
docker-compose up --build
```

The API will be available at `http://localhost:8080`.

## API Endpoints

### Health Check

```
GET /health
```

### Templates

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/templates/` | List all templates |
| GET | `/api/templates/:name` | Get template by name |
| POST | `/api/templates/html` | Render template as HTML |
| POST | `/api/templates/pdf` | Render template as PDF |

### Render HTML

```bash
curl -X POST http://localhost:8080/api/templates/html \
  -H "Content-Type: application/json" \
  -d '{
    "template": "my-template.html",
    "data": {
      "title": "Hello World",
      "items": ["item1", "item2"]
    }
  }'
```

### Render PDF

```bash
curl -X POST http://localhost:8080/api/templates/pdf \
  -H "Content-Type: application/json" \
  -d '{
    "template": "my-template.html",
    "data": {
      "title": "Hello World"
    },
    "options": {
      "paperWidth": 8.5,
      "paperHeight": 11,
      "marginTop": 0.5,
      "marginBottom": 0.5,
      "marginLeft": 0.5,
      "marginRight": 0.5,
      "landscape": false,
      "printBackground": true
    }
  }' --output document.pdf
```

## PDF Options

| Option | Type | Description |
|--------|------|-------------|
| `paperWidth` | float | Paper width in inches (default: 8.5) |
| `paperHeight` | float | Paper height in inches (default: 11) |
| `marginTop` | float | Top margin in inches |
| `marginBottom` | float | Bottom margin in inches |
| `marginLeft` | float | Left margin in inches |
| `marginRight` | float | Right margin in inches |
| `landscape` | bool | Landscape orientation |
| `printBackground` | bool | Print background graphics |
| `preferCssPageSize` | bool | Use CSS-defined page size |
| `waitDelay` | string | Wait duration before conversion (e.g., "1s") |
| `waitForExpression` | string | JS expression to wait for |

## Template Syntax

Templates use Liquid syntax:

- **Variables**: `{{ variableName }}`
- **Conditionals**: `{% if condition %} ... {% endif %}`
- **Loops**: `{% for item in items %} ... {% endfor %}`
- **Filters**: `{{ text | upcase }}`, `{{ text | slice: 0, 1 }}`

## Configuration

Environment variables:

| Variable | Default | Description |
|----------|---------|-------------|
| `PORT` | `8080` | Server port |
| `TEMPLATES_DIR` | `./templates` | Templates directory path |
| `GOTENBERG_URL` | `http://localhost:3000` | Gotenberg service URL |
| `GIN_MODE` | `debug` | Gin mode (`debug`, `release`, `test`) |

## Development

Run locally without Docker:

```bash
# Start Gotenberg
docker run -p 3000:3000 gotenberg/gotenberg:8

# Run the server
go run ./cmd/server
```

## Project Structure

```
├── cmd/server/          # Application entrypoint
├── internal/
│   ├── api/
│   │   ├── handlers/    # HTTP handlers
│   │   ├── middlewares/ # Middleware (error handling)
│   │   └── routes/      # Route registration
│   ├── config/          # Configuration
│   ├── domain/          # Domain models and errors
│   ├── pdfgen/          # PDF generation (Gotenberg client)
│   └── services/        # Business logic
├── templates/           # HTML templates
├── Dockerfile
└── docker-compose.yml
```
