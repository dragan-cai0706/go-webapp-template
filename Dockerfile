FROM golang:1.24 AS build

WORKDIR /workspace

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o freecharge-server ./cmd/server

FROM gcr.io/distroless/base-debian12

WORKDIR /app

COPY --from=build /workspace/freecharge-server .

EXPOSE 8080

ENTRYPOINT ["/app/freecharge-server"]

