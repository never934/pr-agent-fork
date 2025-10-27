FROM golang:1.25.0

WORKDIR /app
RUN ls
COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . ./
RUN go build -v -o /prompt.tuner.api

EXPOSE 4888

CMD ["/prompt.tuner.api"]