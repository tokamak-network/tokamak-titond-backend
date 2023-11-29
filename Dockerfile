FROM --platform=$BUILDPLATFORM golang:1.20-alpine AS builder

RUN apk add --no-cache make gcc musl-dev linux-headers git curl

WORKDIR /app

COPY . .

RUN go install github.com/swaggo/swag/cmd/swag@latest

ARG TARGETOS TARGETARCH
RUN GOOS=$TARGETOS GOARCH=$TARGETARCH make titond

FROM alpine:3.18.4
RUN apk add --no-cache aws-cli ca-certificates jq

COPY --from=builder /app/api /root/api
COPY --from=builder /app/deployments /root/deployments
COPY --from=builder /app/build/bin/titond /usr/local/bin/

WORKDIR /root
ENTRYPOINT ["titond"]