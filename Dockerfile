FROM golang:1.24
LABEL authors="tomato"


WORKDIR /usr/src/app/magic-server


COPY go.mod go.sum ./

RUN go mod download
COPY . .
RUN go build -v cmd/server

CMD ["server"]