FROM golang:1.13-alpine as builder
WORKDIR /
RUN pwd
RUN ls
COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o app .

FROM bysir/alpine-shanghai
COPY --from=0 /app /
ENTRYPOINT ["./app"]
