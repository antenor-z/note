FROM nginx

COPY nginx.conf /etc/nginx/conf.d/default.conf

COPY *.html /usr/share/nginx/html
COPY *.css /usr/share/nginx/html
COPY *.js /usr/share/nginx/html
COPY js/3p/*.js /usr/share/nginx/html/js/3p/
COPY js/*.js /usr/share/nginx/html/js/
COPY img /usr/share/nginx/html/img
COPY favicon /usr/share/nginx/html/favicon

EXPOSE 3003
