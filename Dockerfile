# Fase di compilazione
FROM golang:latest AS build

WORKDIR /web-server
COPY . .
RUN CGO_ENABLED=0 go build -ldflags "-X main.Version=1.0.0" -o /web-server/web-server

# Fase finale
FROM scratch
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
WORKDIR /web-server
COPY --from=build /web-server/media /web-server/media
COPY --from=build /web-server/cert /web-server/cert
COPY --from=build /web-server/templates /web-server/templates
COPY --from=build /web-server/file /web-server/file 
COPY --from=build /web-server/web-server /web-server/web-server
EXPOSE 443
ENTRYPOINT ["/web-server/web-server"]