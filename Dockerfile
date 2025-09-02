FROM golang:latest AS base

WORKDIR /app/

COPY . .

RUN go mod download
RUN go mod verify

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -o slimserve .

FROM scratch

COPY --from=base /app/slimserve /usr/bin/