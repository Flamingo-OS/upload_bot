# Build Stage: Build bot using the alpine image, also install doppler in it
FROM golang:1.19 AS builder
RUN apt-get update && apt-get upgrade -y && apt-get install build-essential -y
WORKDIR /app
COPY . .
RUN GOOS=`go env GOHOSTOS` GOARCH=`go env GOHOSTARCH` go build -o out/upload-bot -ldflags="-w -s" .

# Run Stage: Run bot using the bot and doppler binary copied from build stage
FROM golang:1.19
COPY --from=builder /app/out/upload-bot /app/upload-bot
CMD ["/app/upload-bot"]