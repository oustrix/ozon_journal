FROM golang:1.22.3-alpine3.19 as base
WORKDIR /builder
RUN apk add upx make

ENV GO111MODULE=on CGO_ENABLED=0

COPY go.mod go.sum /builder/
RUN go mod download

COPY . .
RUN make gen
RUN go build -o /builder/main /builder/cmd/app/main.go
RUN upx -9 /builder/main

# runner image
FROM gcr.io/distroless/static:latest
WORKDIR /app
COPY --from=base /builder/main main
COPY --from=base /builder/config/config.yml config/config.yml
COPY --from=base /builder/migrations migrations/


CMD ["/app/main"]