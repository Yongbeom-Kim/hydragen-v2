FROM golang:1.25-alpine AS dev
WORKDIR /app
COPY go.mod ./
RUN go mod download
COPY . .
EXPOSE 8080
CMD ["go", "tool", "air"]

FROM golang:1.25-alpine AS build
WORKDIR /app
COPY go.mod ./
RUN go mod download
COPY . .
RUN go build -o /bin/server ./main.go

FROM alpine:3.20 AS runtime
WORKDIR /app
COPY --from=build /bin/server /usr/local/bin/server
EXPOSE 8080
CMD ["server"]
