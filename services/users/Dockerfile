ARG GO_VERSION=1.23
FROM golang:${GO_VERSION}-alpine as builder

# Install build tools
RUN apk add --no-cache make gcc musl-dev

WORKDIR /app
COPY ../../go.mod ../../go.sum ./
RUN go mod download && go mod verify
COPY . .
RUN make build

FROM alpine:latest

ENV SERVICE_APP_ENV=production
COPY --from=builder /app/main /usr/local/bin/
EXPOSE 8080
CMD ["main"]
