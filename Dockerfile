FROM golang:1.26.1

WORKDIR /usr/src/server
COPY . .

RUN make build
CMD ["./server"]
