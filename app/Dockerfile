FROM golang:1.22

WORKDIR /song_library
COPY go.mod ./
COPY go.sum ./
RUN go mod download
COPY . .
RUN go build -o app

CMD ["/song_library/app"]