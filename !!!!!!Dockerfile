FROM golang:1.22-alpine AS build
WORKDIR /app
COPY main.go .
RUN go mod init proxy && go build -ldflags="-s -w" -o proxy main.go

FROM gcr.io/distroless/static
COPY --from=build /app/proxy /proxy
EXPOSE 8080
ENTRYPOINT ["/proxy"]
