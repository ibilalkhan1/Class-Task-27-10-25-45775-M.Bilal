# Online Bookstore — Microservices CI/CD

## Services
- **user-service** (Node.js + Express + MySQL): signup, login, profile, password-reset stub  
- **catalog-service** (Go + JSON storage): `/books` (GET, POST), `/books/{id}` (GET, PUT, DELETE)

## Run locally
```bash
docker compose up -d --build
curl http://localhost:3000/health
curl http://localhost:4000/health
