FROM golang:1.21-alpine AS builder

RUN apk add --no-cache git ca-certificates

WORKDIR /app

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o c255t19r .

# run
FROM scratch
# FROM gcr.io/distroless/static:nonroot

COPY --from=builder /app/c255t19r /c255t19r

USER nonroot:nonroot

ENV TZ=America/Mexico_City

EXPOSE 8080

ENTRYPOINT [ "/c255t19r" ]
