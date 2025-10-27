FROM golang:1.22-alpine AS build
WORKDIR /src
COPY go.mod ./
COPY src ./src
RUN go build -o /bin/catalog ./src/main.go

FROM alpine:3.20
WORKDIR /app
COPY --from=build /bin/catalog /app/catalog
VOLUME ["/data"]
EXPOSE 4000
CMD ["/app/catalog"]
