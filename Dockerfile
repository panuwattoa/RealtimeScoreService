FROM golang:alpine

RUN apk add bash ca-certificates git gcc g++ libc-dev
WORKDIR /go/rangkingserver

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN go install

RUN GOOS=linux GOARCH=amd64 go build

EXPOSE 12400 8444

CMD [ "gamerangkingserver" ]
