FROM golang:1.25-alpine AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /task-manager-api ./cmd/server

FROM alpine:3.22

RUN addgroup -S app && adduser -S app -G app
WORKDIR /app

COPY --from=build /task-manager-api /app/task-manager-api
USER app

EXPOSE 8089
ENTRYPOINT ["/app/task-manager-api"]
