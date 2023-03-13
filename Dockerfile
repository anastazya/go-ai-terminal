FROM golang:1.19
WORKDIR /app
COPY ./ai.go .
RUN go env -w GO111MODULE=auto && go build -o ai

CMD ["./ai"]