FROM golang:alpine as builder

RUN apk add --no-cache make

COPY . /kaas
WORKDIR /kaas

RUN make build

FROM alpine

WORKDIR /app
COPY --from=builder /kaas/build/manager /app/kaas-manager
COPY --from=builder /kaas/docs  /app/docs

EXPOSE 8080

ENTRYPOINT [ "/app/kaas-manager" ]