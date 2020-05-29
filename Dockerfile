FROM golang:1.14-alpine as build
WORKDIR /opt/beargrabz/

RUN apk add libpcap-dev gcc libc-dev

ADD main.go .
COPY go.mod go.sum ./
RUN go mod download

RUN go build -o beargrabz main.go


FROM alpine
WORKDIR /opt/beargrabz/

RUN apk add libpcap
COPY --from=build /opt/beargrabz/beargrabz .

CMD ["/opt/beargrabz/beargrabz"]