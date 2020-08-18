FROM golang:latest

RUN mkdir -p /app/shared
ADD ./shared /app/shared

RUN mkdir /app/rendezvous-api
ADD ./rendezvous-api /app/rendezvous-api

WORKDIR /app/rendezvous-api
RUN go build -o main .

CMD ["/app/rendezvous-api/main"]