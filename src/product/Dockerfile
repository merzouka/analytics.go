FROM golang:alpine AS base
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . ./
RUN go build -o /server .

FROM alpine:3.20 AS release
WORKDIR /

COPY --from=base /server /server

RUN apk update
RUN apk --no-cache add curl

CMD [ "/server" ]
