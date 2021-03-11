ARG GO_VERSION=1.16
ARG UI_VERSION=latest

FROM golang:${GO_VERSION}-alpine AS binary-builder

RUN apk add --no-cache ca-certificates git

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY ./ ./

RUN CGO_ENABLED=0 go build \
    -ldflags="-s -w" \
    -installsuffix 'static' \
    -o /app .



FROM node:12.7-alpine AS ui-builder
WORKDIR /usr/src/app
COPY ui/ .
RUN npm install -g @angular/cli
RUN npm install
RUN ng build --prod --base-href /ui/


FROM alpine AS final

COPY --from=binary-builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

COPY --from=binary-builder /app /app

RUN mkdir -p /ui
COPY --from=ui-builder /usr/src/app/dist /ui

COPY --from=binary-builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

ENTRYPOINT ["/app"]