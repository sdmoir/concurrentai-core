FROM golang:latest

RUN mkdir -p /app/shared
ADD ./shared /app/shared

RUN mkdir /app/model-executor
ADD ./model-executor /app/model-executor

WORKDIR /app/model-executor
RUN go build -o main .

CMD ["/app/model-executor/main"]