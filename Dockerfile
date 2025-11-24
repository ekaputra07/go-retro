FROM golang:1.24-alpine AS go-builder

# Set the working directory
WORKDIR /app

# Download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build app
RUN CGO_ENABLED=0 GOOS=linux go build -v -o dist/goretro-web ./cmd/web

FROM node:22-alpine AS web-builder

# Set the working directory
WORKDIR /app

# Copy source code
COPY web web

# Build assets
WORKDIR /app/web/assets
RUN npm install && npm run build

FROM alpine:edge

ARG version="latest"
ENV GORETRO_VERSION=${version}

# Set the working directory
WORKDIR /app

# copy the binary and statics from the build stage
COPY --from=go-builder /app/dist/goretro-web .
COPY --from=web-builder /app/web/templates ./web/templates
COPY --from=web-builder /app/web/public ./web/public

# Set the entry point
ENTRYPOINT [ "./goretro-web" ]