FROM golang:1.24-alpine AS go-builder

# Set the working directory
WORKDIR /app

# Download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build app
RUN CGO_ENABLED=0 GOOS=linux go build ./cmd/web -o go-retro

FROM node:22-alpine AS ui-builder

# Set the working directory
WORKDIR /app

# Copy source code
COPY ui ui

# Build assets
WORKDIR /app/ui/assets
RUN npm install && npm run build

FROM alpine:edge

# Set the environment variables
ENV GORETRO_HOST=0.0.0.0
ENV GORETRO_PORT=8080

# Set the working directory
WORKDIR /app

# copy the binary and statics from the build stage
COPY --from=go-builder /app/go-retro .
COPY --from=ui-builder /app/ui/templates ./ui/templates
COPY --from=ui-builder /app/ui/public ./ui/public

# Set the entry point
ENTRYPOINT [ "./go-retro" ]