FROM golang:latest

RUN mkdir -p /app/shared
ADD ./shared /app/shared

RUN mkdir /app/rendezvous-collector
ADD ./rendezvous-collector /app/rendezvous-collector

WORKDIR /app/rendezvous-collector
RUN go build -o main .

CMD ["/app/rendezvous-collector/main"]