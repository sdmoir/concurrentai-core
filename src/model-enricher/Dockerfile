FROM golang:latest

RUN mkdir -p /app/shared
ADD ./shared /app/shared

RUN mkdir /app/model-enricher
ADD ./model-enricher /app/model-enricher

WORKDIR /app/model-enricher
RUN go build -o main .

CMD ["/app/model-enricher/main"]