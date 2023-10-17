FROM golang:1.20-alpine AS builder

ENV GO111MODULE=on

WORKDIR /build

COPY . .
RUN go mod download


RUN CGO_ENABLED=0 GOOS=linux go build -o ./interview-test

FROM alpine:latest as baseImage

WORKDIR /opt

COPY --from=builder /build/interview-test /opt/interview-test

CMD ["/opt/interview-test"]
