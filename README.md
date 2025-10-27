# Online Bookstore — Microservices CI/CD

## Services
- **user-service** (Node.js + Express + MySQL): signup, login, profile, password-reset stub
- **catalog-service** (Go + JSON storage): `/books` (GET, POST), `/books/{id}` (GET, PUT, DELETE)

## Run locally
```bash
docker compose up -d --build
# health
curl -s http://localhost:3000/health
curl -s http://localhost:4000/health

# user
curl -s -X POST http://localhost:3000/signup -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"pass123","name":"Bilal"}'
curl -s -X POST http://localhost:3000/login -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"pass123"}'

# catalog
curl -s -X POST http://localhost:4000/books -H "Content-Type: application/json" \
  -d '{"id":"b1","title":"Clean Code","author":"Robert C. Martin","price":29.99,"available":true}'
curl -s -X PUT  http://localhost:4000/books/b1 -H "Content-Type: application/json" \
  -d '{"price":31.99,"available":false}'
curl -s       http://localhost:4000/books/b1
curl -si -X DELETE http://localhost:4000/books/b1 | head -n1
curl -s       http://localhost:4000/books
