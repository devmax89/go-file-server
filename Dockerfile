# Fase di compilazione
FROM golang:latest AS build

WORKDIR /src
COPY . .
RUN CGO_ENABLED=0 go build -o /bin/app

# Fase finale
FROM scratch
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /src/cert /cert
COPY --from=build /src/templates /templates
COPY --from=build /src/file /file 
COPY --from=build /bin/app /bin/app
EXPOSE 443
ENTRYPOINT ["/bin/app"]