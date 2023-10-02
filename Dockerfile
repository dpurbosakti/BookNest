FROM alpine:3.18
RUN apk update && apk add --no-cache coreutils openssl1.1-compat procps
RUN addgroup -g 1001 -S go
RUN adduser -S go -u 1001

USER go

# ganti dengan APP_PORT sesuai dengan config.yaml
EXPOSE 8080

ARG APP_NAME=booknest
ARG PATH=/var/opt/${APP_NAME}

WORKDIR ${PATH}
COPY config.yaml .
COPY ${APP_NAME} .

CMD ["./booknest", "serve"]

#? how to run this dockerfile
# ./start_dockerfile.sh
