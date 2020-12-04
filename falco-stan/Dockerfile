FROM golang:alpine3.12 AS builder

RUN mkdir /build
ADD main.go /build/
WORKDIR /build

RUN apk --no-cache add git
# Install all deps
RUN go get -u github.com/nats-io/stan.go && go get -u github.com/nats-io/nats.go

# Build 'stan-pub' binary
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-s" -o stan-pub main.go

FROM alpine:latest
COPY --from=builder /build/stan-pub /bin/stan-pub

CMD ["/bin/stan-pub"]

