FROM nginx:stable-alpine

MAINTAINER Alexander Teplov <teploff.aa@gmail.com>
LABEL maintainer="teploff.aa@gmail.com"

ADD nginx.conf /etc/nginx/nginx.conf
COPY sites-enabled /etc/nginx/sites-enabled

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]