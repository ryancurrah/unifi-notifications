# ---- Build container
FROM golang:alpine AS builder
WORKDIR /unifi-notifications
COPY . .
RUN apk add --no-cache git
RUN go build -v ./...

# ---- App container
FROM alpine:latest as unifi-notifications
ENV NOTIFCATION_SERVICES=slack
ENV UNIFI_URL=
ENV UNIFI_SITES=
ENV UNIFI_USERNAME=
ENV UNIFI_PASSWORD=
ENV SLACK_ALARMS_WEBHOOK=
ENV SLACK_EVENTS_WEBHOOK=
RUN apk --no-cache add ca-certificates
COPY --from=builder unifi-notifications/unifi-notifications /
ENTRYPOINT ./unifi-notifications
LABEL Name=unifi-notifications Version=0.0.1
