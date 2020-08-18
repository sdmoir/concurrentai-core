FROM golang:latest

RUN mkdir -p /app/shared
ADD ./shared /app/shared

RUN mkdir /app/api
ADD ./api /app/api

WORKDIR /app/api
RUN go build -o main .

CMD ["/app/api/main"]