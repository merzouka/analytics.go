FROM golang:alpine AS build
WORKDIR /app

COPY go.* ./
RUN go mod download

COPY . ./
RUN CGO_ENABLED=0 GOOS=linux go build -o /server

FROM alpine:3.20 AS release
WORKDIR /

COPY --from=build /server /server

CMD ["/server"]

