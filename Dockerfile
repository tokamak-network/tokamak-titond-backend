FROM --platform=$BUILDPLATFORM golang:1.20-alpine AS builder

RUN apk add --no-cache make gcc musl-dev linux-headers git

WORKDIR /app

COPY . .

ARG TARGETOS TARGETARCH
RUN GOOS=$TARGETOS GOARCH=$TARGETARCH make titond

FROM alpine

RUN apk add --no-cache ca-certificates jq curl
COPY --from=builder /app/api /root/api
COPY --from=builder /app/deployments /root/deployments
COPY --from=builder /app/build/bin/titond /usr/local/bin/

WORKDIR /root
ENTRYPOINT ["titond"]