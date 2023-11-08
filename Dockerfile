FROM golang:1.20-alpine AS builder

RUN apk add --no-cache make gcc musl-dev linux-headers git

WORKDIR /app

COPY . .

RUN make titond

FROM alpine

RUN apk add --no-cache ca-certificates jq curl
COPY --from=builder /app/deployments /usr/local/bin/
COPY --from=builder /app/api /usr/local/bin/
COPY --from=builder /app/build/bin/titond /usr/local/bin/

WORKDIR /usr/local/bin/
ENTRYPOINT ["titond"]