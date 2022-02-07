FROM golang:1.17-alpine AS build

WORKDIR /build
COPY . /build
RUN CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -o dns-updater

FROM alpine

LABEL maintener "Atomys<contact@atomys.fr>"
LABEL repository "https://gitlab.com/Atomys/dns-updater"

COPY --from=build /build/dns-updater /dns-updater
COPY config/ config/

CMD [ "/dns-updater" ]