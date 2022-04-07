FROM golang:alpine AS builder
ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64 
RUN mkdir -p /data/www/shop
WORKDIR /data/www/shop
RUN apk --update --no-cache add ca-certificates gcc libtool make musl-dev protoc
COPY Makefile go.mod go.sum ./
RUN make init
COPY . .
RUN make proto tidy build

FROM scratch
COPY --from=builder /data/www/shop/shopcategory /shopcategory
ENTRYPOINT ["/shopcategory"]
CMD []
