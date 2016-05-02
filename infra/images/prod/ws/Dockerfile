FROM nginx:alpine
COPY files/gophr.pm.csr /etc/nginx/
COPY files/gophr.pm.key /etc/nginx/
COPY files/gophr.pm.crt /etc/nginx/
VOLUME ["/usr/share/nginx/html"]
