FROM golang:1-alpine

LABEL maintainer="Shaikhzidhin <sinuzidin@gmail.com>"

WORKDIR /app

COPY . .

RUN  go build

RUN go build -o main

EXPOSE 3000

CMD [ "./main" ]