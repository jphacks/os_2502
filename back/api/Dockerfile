FROM golang:1.24-alpine

WORKDIR /work

RUN apk add --no-cache make

COPY go.mod go.sum ./
RUN go mod download
COPY . .

RUN go install github.com/cosmtrek/air@v1.40.4
RUN go install github.com/aarondl/sqlboiler/v4@latest
RUN go install github.com/aarondl/sqlboiler/v4/drivers/sqlboiler-mysql@latest

CMD ["make", "start"]
