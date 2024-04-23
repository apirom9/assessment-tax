FROM golang:1.21-alpine as build-base
WORKDIR /app
COPY go.mod .
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go test -v
RUN go build -o ./out/tax-api .

FROM alpine:3.16.2
COPY --from=build-base /app/out/tax-api /app/tax-api
CMD ["/app/tax-api"]