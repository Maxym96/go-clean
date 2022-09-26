FROM golang:1.18-alpine

RUN set -ex && \
    apk add --no-cache gcc musl-dev

WORKDIR /app/src

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

COPY migrations ../migrations

RUN go build -o ../go-clean cmd/app/main.go

WORKDIR /app

RUN rm -rf src

EXPOSE 8080

CMD [ "/app/go-clean" ]