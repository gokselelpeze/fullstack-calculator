FROM node:24-alpine AS frontend
WORKDIR /app/frontend
COPY frontend/package*.json ./
RUN npm ci
COPY frontend/ ./
RUN npm run build

FROM golang:1.27.1-alpine AS backend
WORKDIR /app/backend
COPY backend/go.mod backend/go.sum ./
RUN go mod download
COPY backend/ ./
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /calculator ./cmd/server

FROM scratch
WORKDIR /app
COPY --from=backend /calculator /app/calculator
COPY --from=frontend /app/frontend/dist /app/public
ENV ADDR=:8080 STATIC_DIR=/app/public
USER 65532:65532
EXPOSE 8080
ENTRYPOINT ["/app/calculator"]
