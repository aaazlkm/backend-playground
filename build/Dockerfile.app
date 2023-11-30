# buildステージ
FROM golang:1.21.4-alpine3.18 AS build
WORKDIR /server
COPY server .
RUN go mod download && go build -o main .

# 実行ステージ
FROM scratch
WORKDIR /server
COPY --from=build server/main .
ENTRYPOINT [ "./main" ]
