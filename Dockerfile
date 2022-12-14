FROM golang:1.19.0-alpine3.16 AS build

WORKDIR /app

FROM build as dev
# Install git.
# Git is required for fetching the dependencies.
RUN apk update \
 && apk add --no-cache git \
 && apk add --no-cach bash  \
 && apk add build-base

COPY go.mod ./
COPY go.sum ./

RUN go mod download
RUN go get -u golang.org/x/lint/golint \
&&  go get -u github.com/gofiber/fiber/v2 \
&&  go get -u github.com/swaggo/swag/cmd/swag \
&&  go get -u github.com/arsmn/fiber-swagger/v2 \
&&  go get -u github.com/lib/pq

COPY ./db ./db/
COPY ./docs ./docs/
COPY ./main.go ./

RUN go build -o /app/aroundhome

FROM alpine:3.16 AS release
WORKDIR /app
COPY --from=dev /app/aroundhome ./

CMD ["./aroundhome"]

EXPOSE 3000