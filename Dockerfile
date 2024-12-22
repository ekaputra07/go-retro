FROM golang:1.23-alpine AS build

# Set the working directory
WORKDIR /app

# Download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o go-retro

FROM alpine:edge

# Set the environment variables
ENV GORETRO_HOST=0.0.0.0
ENV GORETRO_PORT=8080

# Set the working directory
WORKDIR /app

# copy the binary and statics from the build stage
COPY --from=build /app/go-retro .
COPY --from=build /app/web web

# Set the entry point
ENTRYPOINT [ "./go-retro" ]