FROM golang:1.22.5-alpine

RUN apk add --no-cache gcc musl-dev


WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY *.go ./

RUN go get golang-assignment

RUN CGO_ENABLED=1 GOOS=linux go build -o /golang-assignment

EXPOSE 8080
CMD [ "/golang-assignment" ]