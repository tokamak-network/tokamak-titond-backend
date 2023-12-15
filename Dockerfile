FROM --platform=$BUILDPLATFORM golang:1.20-alpine AS builder

RUN apk add --no-cache make gcc musl-dev linux-headers git curl

WORKDIR /app

COPY . .

RUN go install github.com/swaggo/swag/cmd/swag@latest
RUN make install

ENTRYPOINT ["/bin/sh"]