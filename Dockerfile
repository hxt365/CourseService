FROM golang:1.14

RUN mkdir /code

WORKDIR /code

COPY go.mod go.sum /code/

RUN go mod download

COPY . /code/

EXPOSE 8080