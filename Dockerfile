FROM golang:1.9-alpine as builder
WORKDIR /builder
RUN CGO_ENABLED=0 GOOS=linux go build -o app .

FROM bysir/alpine-shanghai
COPY --from=0 /builder/app /
ENTRYPOINT ["./app"]
