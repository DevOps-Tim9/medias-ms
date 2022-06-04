FROM golang:latest as build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o src ./src

FROM alpine:latest as run

WORKDIR /root/

COPY --from=build /app/src .

CMD ["./src"]



