FROM golang:1.19-alpine

WORKDIR /go/ambrosia-atlas-api
COPY . /go/ambrosia-atlas-api

RUN go mod tidy
RUN go mod download
RUN go build -o /ambrosia-atlas-api

EXPOSE 8080

CMD [ "/ambrosia-atlas-api" ]
