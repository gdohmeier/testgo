FROM golang:1.23-alpine AS build
WORKDIR /src
COPY go.mod main.go ./
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/testgo .

FROM alpine:3.20
RUN apk add --no-cache wget && adduser -D -H -u 10001 app
USER app
WORKDIR /home/app
COPY --from=build /out/testgo /usr/local/bin/testgo
EXPOSE 8080
ENV PORT=8080
HEALTHCHECK --interval=15s --timeout=3s --retries=3 CMD wget -qO- http://127.0.0.1:8080/up || exit 1
CMD ["testgo"]
