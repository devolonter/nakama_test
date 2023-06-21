FROM heroiclabs/nakama-pluginbuilder:3.16.0 AS go-builder

ENV GO111MODULE on
ENV CGO_ENABLED 1

WORKDIR /test

COPY go.mod .
COPY main.go .
COPY vendor/ vendor/

RUN go build --trimpath --mod=vendor --buildmode=plugin -o ./test.so

FROM registry.heroiclabs.com/heroiclabs/nakama:3.16.0

COPY --from=go-builder /test/test.so /nakama/data/modules/
COPY local.yml /nakama/data/
COPY data/ /nakama/data/modules/data/