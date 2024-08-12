FROM golang:alpine AS base
WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY . ./
RUN CGO_ENABLED=0 GOOS=linux go build -o /seeder .

FROM alpine:3.20 AS release
WORKDIR /

COPY --from=base /seeder /seeder

CMD [ "/seeder" ]
