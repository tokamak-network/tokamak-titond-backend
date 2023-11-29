FROM --platform=$BUILDPLATFORM golang:1.20-alpine AS builder

RUN apk add --no-cache make gcc musl-dev linux-headers git curl

WORKDIR /app

COPY . .

RUN go install github.com/swaggo/swag/cmd/swag@latest

ARG TARGETOS TARGETARCH
RUN GOOS=$TARGETOS GOARCH=$TARGETARCH make titond

RUN echo $TARGETARCH

RUN if [ "$TARGETARCH" = "amd64" ]; then \
        curl "https://awscli.amazonaws.com/awscli-exe-linux-x86_64.zip" -o "awscliv2.zip"; \
    elif [ "$TARGETARCH" = "arm64" ]; then \
        curl "https://awscli.amazonaws.com/awscli-exe-linux-aarch64.zip" -o "awscliv2.zip"; \
    fi
RUN unzip awscliv2.zip

FROM frolvlad/alpine-glibc

RUN apk add --no-cache ca-certificates jq

COPY --from=builder /app/api /root/api
COPY --from=builder /app/deployments /root/deployments
COPY --from=builder /app/build/bin/titond /usr/local/bin/
COPY --from=builder /app/aws aws

RUN ./aws/install

WORKDIR /root
ENTRYPOINT ["titond"]