#docker build --pull --rm -f "urlshortener\build\pkg\Dockerfile" -t urlshortener "urlshortener"
FROM golang:1.17 as builder

WORKDIR /go/src/app
COPY . .

#RUN ls && exit 1
RUN GOOS=linux GOARCH=amd64 \
    go build -o /build/bin/api ./cmd

#EXPOSE 80

#CMD ["/build/bin/api"]

FROM ubuntu

COPY --from=builder /build/bin/api /build/bin/api

EXPOSE 80
ENTRYPOINT ["/build/bin/api"]