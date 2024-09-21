FROM golang:1-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN go build -o /app/verbilobot .

FROM alpine:3

RUN apk add ffmpeg
COPY --from=builder /app/verbilobot /usr/local/bin

CMD ["verbilobot"]
