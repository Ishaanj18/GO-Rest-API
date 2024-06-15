FROM golang:1.22.1-alpine3.19

WORKDIR /app

COPY ./pkg /app

RUN go build -o library_api

CMD [ "./library_api" ]

