# Submission Checklist (Oct 27, 2025)

## Microservices
- [x] user-service (Node + Express + MySQL): signup, login, profile, password-reset stub
- [x] catalog-service (Go + JSON storage): /books (GET, POST), /books/{id} (GET, PUT, DELETE)

## Containerization & Compose
- [x] Dockerfiles for both services
- [x] docker-compose.yml (MySQL 8.4 exposed on host 3307)

## REST & Tests
- [x] REST endpoints implemented
- [x] user-service tests (Jest): health + auth
- [x] catalog-service tests (go test): read/write; handlers covered by smoke in CD

## CI/CD
- [x] CI (user-service): Node build, Jest tests, Docker build
- [x] CI (catalog-service): Go build, go test, Docker build
- [x] CD: deploy-compose workflow builds, brings up stack, waits for /health, runs smoke checks, tears down

## Documentation & Evidence
- [x] README with run cmds and curl examples
- [x] ASSUMPTIONS.md and KNOWN_ISSUES.md
- [ ] evidence/ screenshots: compose up, health OK, curl flows, Actions green

## Repo Link
- https://github.com/ibilalkhan1/Class-Task-27-10-25-45775-M.Bilal

