FROM nginx:alpine
COPY files/cert.crt /etc/nginx/
COPY files/cert.csr /etc/nginx/
COPY files/cert.key /etc/nginx/
COPY files/nginx.conf /etc/nginx/
COPY assets/* /usr/share/nginx/html/
