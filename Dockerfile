# Build stage
# ---------------
FROM alpine AS build

WORKDIR /app

# Request latest golang package
RUN apk add go \
--no-cache \
--repository=https://dl-cdn.alpinelinux.org/alpine/edge/community

# Copy deps list
COPY go.mod go.sum ./
RUN go mod download

# Copy source codes
COPY . /app/

# Optimize build size/speed
ENV GOCACHE=/root/.cache/go-build

RUN --mount=type=cache,target="/root/.cache/go-build" \
  go build -o /app/main

RUN chmod 500 /app/main

# Runtime stage
# ---------------
FROM scratch AS runtime

USER 1000:1000
WORKDIR /app

# Copy final binary & TLS Root CA certs
COPY --from=build --chown=1000:1000 /app/main .
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

ENTRYPOINT ["/app/main"]
  