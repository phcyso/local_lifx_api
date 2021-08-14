FROM golang:1.16.7-alpine as builder

WORKDIR /app 

COPY go.mod .

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" .

FROM golang:1.16.7-alpine

WORKDIR /app

COPY --from=builder /app/local_lifx_api /usr/bin/

COPY ./scenes.yaml /app/scenes.yaml

CMD ["local_lifx_api"]
