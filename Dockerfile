FROM golang:1.19-alpine

WORKDIR /go/ambrosia-atlas-api
COPY . /go/ambrosia-atlas-api

RUN go mod tidy
RUN go mod download

EXPOSE 5000

CMD [ "go", "run", "." ]
