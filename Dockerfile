FROM alpine:3

WORKDIR /app
# Now just add the binary
ADD fairservice /app
ENTRYPOINT ["./fairservice"]
EXPOSE 3000
ENV DB_HOSTPORT localhost:5432
ENV DB_PWD secret
ENV DB_NAME hgkforms
ENV DB_USR hgkforms
