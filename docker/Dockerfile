FROM golang:1.22-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN go build -o /app/verbilobot .

FROM jrottenberg/ffmpeg:7-alpine

COPY --from=builder /app/verbilobot /usr/local/bin

ENTRYPOINT ["verbilobot"]
CMD [""]