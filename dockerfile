FROM golang:1.21 AS builder

WORKDIR /usr/local/fakeSSH
ADD . /usr/local/fakeSSH
RUN go build

FROM golang:1.21

WORKDIR /usr/local/fakeSSH
COPY --from=builder /usr/local/fakeSSH/fakeSSH /usr/local/bin

CMD fakeSSH
