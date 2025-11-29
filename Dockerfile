# --- Build UI ---
FROM node:22-alpine AS ui-builder

# Set the working directory
WORKDIR /app

# Copy source code
COPY web/ui ui

# Build assets
WORKDIR /app/ui
RUN npm install && npm run build

# --- Build golang app ---
FROM golang:1.24-alpine AS go-builder

# Set the working directory
WORKDIR /app

# Download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .
# Copy UI dist
COPY --from=ui-builder /app/ui/dist ./web/ui/dist

# Build app
RUN CGO_ENABLED=0 GOOS=linux go build -v -o dist/goretro-web ./cmd/web

# --- Package it! ---
FROM alpine:edge

ARG version="latest"
ENV GORETRO_VERSION=${version}

# Set the working directory
WORKDIR /app

# copy binary from the build stage
COPY --from=go-builder /app/dist/goretro-web .

# Set the entry point
ENTRYPOINT [ "./goretro-web" ]