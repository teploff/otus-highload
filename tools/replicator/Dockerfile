FROM golang:1.15 as builder
LABEL mainater="Alexander Teplov teploff.aa@gmail.com"
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOSE=linux GO111MODULE=on go build -mod=vendor -a -installsuffix nocgo -o replicator /app/cmd/replicator/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY --from=builder /app/ ./

CMD ["./replicator", "--config=./configs/config.yaml"]