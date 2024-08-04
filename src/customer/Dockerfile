FROM golang:1.22-alpine AS base
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . ./
RUN CGO_ENABLED=0 GOOS=linux go build -o /server .

FROM alpine:3.14 AS release
WORKDIR /

COPY --from=base /server /server

RUN apk update
RUN apk --no-cache add curl

ENV GIN_MODE=release
EXPOSE 8080
CMD [ "/server" ]

