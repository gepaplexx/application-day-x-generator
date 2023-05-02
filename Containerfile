###############################################################################
################################# build stage #################################
###############################################################################
FROM golang:1.20.4-alpine3.16 as builder

ENV GO111MODULE=on

WORKDIR /app

COPY src/go.mod .
COPY src/go.sum .

RUN go mod download

COPY src/ .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build

###############################################################################
################################# final stage #################################
###############################################################################
FROM alpine

ARG KUBESEAL_VERSION=v0.16.0

RUN apk add --update --no-cache curl

# Install kubeseal
RUN curl -sL https://github.com/bitnami-labs/sealed-secrets/releases/download/${KUBESEAL_VERSION}/kubeseal-linux-amd64 -o kubeseal && \
    mv kubeseal /usr/bin/kubeseal && \
    chmod +x /usr/bin/kubeseal

WORKDIR /app

COPY --from=builder /app/day-x-generator /app/
COPY ./src/templates /app/templates

ENTRYPOINT ["/app/day-x-generator"]