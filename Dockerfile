FROM alpine:3.21

COPY venstar-monitor /bin/

WORKDIR /app

EXPOSE 9872

ENTRYPOINT [ "/bin/venstar-monitor" ]
