FROM golang:latest

RUN mkdir -p /app/shared
ADD ./shared /app/shared

RUN mkdir /app/collector
ADD ./collector /app/collector

WORKDIR /app/collector
RUN go build -o main .

CMD ["/app/collector/main"]