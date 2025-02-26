FROM nginx

COPY nginx.conf /etc/nginx/conf.d/default.conf

COPY *.html /usr/share/nginx/html
COPY *.css /usr/share/nginx/html
COPY *.js /usr/share/nginx/html
COPY img /usr/share/nginx/html/img

EXPOSE 3003
