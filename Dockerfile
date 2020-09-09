FROM golang:1.14-alpine as dev

WORKDIR /go/src/dice
COPY ./go.mod .
#COPY ./go.sum .

RUN go mod download

COPY . .
RUN go install github.com/dustinpianalto/dice

CMD [ "go", "run", "main.go"]

from alpine

WORKDIR /bin

COPY --from=dev /go/bin/dice ./dice

CMD [ "dice" ]
