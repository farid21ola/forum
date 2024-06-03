FROM golang:1.20

RUN go version
ENV GOPATH=/

COPY ./ ./

RUN apt-get update
RUN apt-get -y install postgresql-client

RUN go mod download
RUN go build -o forum ./server.go

CMD ["./forum"]