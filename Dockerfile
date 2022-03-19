FROM golang:1.15.7-buster

COPY go.mod go.sum /go/src/github.com/gkranasinghe/letterpress/
WORKDIR /go/src/github.com/gkranasinghe/letterpress
RUN go mod download
COPY . /go/src/github.com/gkranasinghe/letterpress
RUN go build -o /usr/bin/letterpress /go/src/github.com/gkranasinghe/letterpress

EXPOSE 8080 8080
ENTRYPOINT ["/usr/bin/letterpress"]