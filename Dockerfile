FROM dwarvesf/sql-migrate-git as migrate

FROM golang:1.16 as builder
WORKDIR /go/src/github.com/nnhuyhoang/simple_rest_project/backend
COPY . .
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o server cmd/server/*.go

FROM alpine:latest  
RUN apk --no-cache add ca-certificates
RUN apk --no-cache add tzdata
WORKDIR /app
COPY --from=builder /go/src/github.com/nnhuyhoang/simple_rest_project/backend/server /app
CMD ["/app/server"]  