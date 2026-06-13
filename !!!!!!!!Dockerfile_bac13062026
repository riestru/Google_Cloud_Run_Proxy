FROM nginx:1.27-alpine

RUN apk add --no-cache gettext

COPY nginx.conf.template /etc/nginx/nginx.conf.template

RUN mkdir -p /etc/nginx/templates

EXPOSE 8080

CMD ["/bin/sh", "-c", "envsubst '$V2RAY_SERVER_IP1 $V2RAY_SERVER_IP2 $V2RAY_SERVER_IP3 $V2RAY_SERVER_IP5' < /etc/nginx/nginx.conf.template > /etc/nginx/nginx.conf && nginx -t && nginx -g 'daemon off;'"]
