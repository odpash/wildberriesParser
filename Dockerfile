FROM golang:latest
RUN mkdir /app
ADD . /app/
WORKDIR /app
RUN go build -o main .
RUN chmod +rwx /app/main
CMD ["/app/main"]
