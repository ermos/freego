FROM golang:1.21-alpine AS builder

WORKDIR /build

COPY go.mod go.sum ./

RUN go mod download

COPY main.go ./
COPY cmd ./cmd
COPY internal ./internal

RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-w -s" -o freego -v ./

FROM alpine:latest

RUN mkdir -p /freego

COPY --from=builder /build/freego /bin/freego

EXPOSE 80
EXPOSE 443

ENTRYPOINT ["freego"]
CMD ["serve", "--docker"]
