FROM alpine:3.21.2 as certs

FROM scratch
COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY wms /
ENTRYPOINT ["/wms"]