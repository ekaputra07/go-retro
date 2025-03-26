FROM golang:1.23-alpine AS go-builder

# Set the working directory
WORKDIR /app

# Download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build app
RUN CGO_ENABLED=0 GOOS=linux go build -o go-retro

FROM node:22-alpine AS web-builder

# Set the working directory
WORKDIR /app

# Copy source code
COPY web web

# Build assets
WORKDIR /app/web/assets
RUN npm install && npm run build

FROM alpine:edge

# Set the environment variables
ENV GORETRO_HOST=0.0.0.0
ENV GORETRO_PORT=8080

# Set the working directory
WORKDIR /app

# copy the binary and statics from the build stage
COPY --from=go-builder /app/go-retro .
COPY --from=web-builder /app/web/templates ./web/templates
COPY --from=web-builder /app/web/public ./web/public

# Set the entry point
ENTRYPOINT [ "./go-retro" ]