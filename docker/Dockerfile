FROM --platform=$BUILDPLATFORM library/golang:1.22-alpine as builder
ARG BUILDPLATFORM
WORKDIR /build/src
COPY . .
RUN GOOS=linux GOARCH=$TARGETARCH go build -o release/$BUILDPLATFORM/docker-cache cmd/main.go

FROM --platform=$BUILDPLATFORM library/alpine:3
LABEL maintainer="kit101 <qkssk1711@163.com>"
LABEL repository="https://github.com/kit101/drone-plugin-docker-cache"

ARG BUILDPLATFORM
RUN apk add jq
#COPY release/$BUILDPLATFORM/docker-cache /usr/local/bin/docker-cache
COPY --from=builder /build/src/release/$BUILDPLATFORM/docker-cache /usr/local/bin/docker-cache
ENTRYPOINT ["/usr/local/bin/docker-cache"]