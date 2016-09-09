FROM node:latest

# Install nginx and build tools.
ENV NGINX_VERSION 1.11.3-1~jessie
RUN apt-key adv --keyserver hkp://pgp.mit.edu:80 --recv-keys 573BFD6B3D8FBC641079A6ABABF5BD827BD9BF62 \
	&& echo "deb http://nginx.org/packages/mainline/debian/ jessie nginx" >> /etc/apt/sources.list \
	&& apt-get update \
	&& apt-get install --no-install-recommends --no-install-suggests -y \
		ca-certificates \
		nginx=${NGINX_VERSION} \
		nginx-module-geoip \
		nginx-module-image-filter \
		nginx-module-perl \
		nginx-module-njs \
		gettext-base \
		gitweb \
		autoconf \
		build-essential \
		automake \
		libtool \
		libfcgi-dev \
		apache2-utils

# Build and install fcgiwrap
RUN git clone https://github.com/gnosek/fcgiwrap.git \
		&& cd fcgiwrap/ \
		&& autoreconf -i \
		&& ./configure \
		&& make \
		&& make install \
		&& cp fcgiwrap /usr/bin

# Copy in gitweb and nginx configurations
COPY ./files/nginx/gitweb.conf /etc/gitweb.conf
COPY ./files/nginx/git.conf /etc/nginx/conf.d/git.conf
COPY ./files/nginx/htpasswd /etc/nginx/

# Copy in spawn-fcgi scripts
COPY ./files/bin/spawn-fcgi /usr/bin/spawn-fcgi
RUN chmod +x /usr/bin/spawn-fcgi

# Copy in spawn-fcgi scripts
COPY ./files/init.d/spawn-fcgi /etc/init.d/spawn-fcgi
RUN chmod +x /etc/init.d/spawn-fcgi

COPY ./files/hosts /etc/hosts

CMD service nginx start

EXPOSE 80 443
