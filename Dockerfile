FROM golang:alpine as builder

RUN apk add --no-cache make

COPY . /kaas
WORKDIR /kaas

RUN make build

FROM alpine

COPY --from=builder /kaas/build/manager /bin/kaas-manager

EXPOSE 8080

ENTRYPOINT [ "/bin/kaas-manager" ]