FROM golang:1.9-alpine as builder
RUN CGO_ENABLED=0 GOOS=linux go build -o app .

FROM bysir/alpine-shanghai
COPY --from=0 app /
ENTRYPOINT ["./app"]
