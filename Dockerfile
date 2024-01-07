FROM node:latest AS frontend-builder
WORKDIR /frontend
COPY /frontend ./
RUN npm i && \
  npm run build

FROM golang:alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
COPY --from=frontend-builder /frontend/dist ./frontend/dist
RUN go mod download
# Copy everything from the root directory
COPY . .
RUN go build -v -o ./server ./cmd/server/

# This starts a new container and copies the binary here
FROM alpine
WORKDIR /app
COPY .env .env
COPY --from=builder /app/server ./server
CMD ./server
