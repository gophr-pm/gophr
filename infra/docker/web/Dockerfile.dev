FROM node:latest

# Set the default timezone to EST.
ENV TZ=America/New_York
RUN echo $TZ | tee /etc/timezone \
	&& dpkg-reconfigure --frontend noninteractive tzdata

# Install nginx and yarn.
ENV NGINX_VERSION 1.11.6-1~jessie
RUN apt-key adv --keyserver hkp://pgp.mit.edu:80 --recv-keys 573BFD6B3D8FBC641079A6ABABF5BD827BD9BF62 \
	&& echo "deb http://nginx.org/packages/mainline/debian/ jessie nginx" >> /etc/apt/sources.list \
	&& curl -sS https://dl.yarnpkg.com/debian/pubkey.gpg | apt-key add - \
	&& echo "deb http://dl.yarnpkg.com/debian/ stable main" | tee /etc/apt/sources.list.d/yarn.list \
	&& apt-get update \
	&& apt-get install --no-install-recommends --no-install-suggests -y \
		ca-certificates \
		nginx=${NGINX_VERSION} \
		nginx-module-geoip \
		nginx-module-image-filter \
		nginx-module-perl \
		nginx-module-njs \
		gettext-base \
		yarn \
	&& rm -rf /var/lib/apt/lists/*

# Forward nginx request and error logs to docker log collector
RUN ln -sf /dev/stdout /var/log/nginx/access.log \
	&& ln -sf /dev/stderr /var/log/nginx/error.log

COPY ./web /web
COPY ./infra/docker/web/files/start.dev.sh /start.sh
COPY ./infra/docker/web/files/nginx.conf /etc/nginx/nginx.conf

RUN cd /web \
	&& yarn \
	&& yarn run build \
	&& mkdir -p /usr/share/nginx/html \
	&& mv ./build/* /usr/share/nginx/html \
	&& cd /etc/nginx/ \
	&& rm -rf /web

RUN chmod +x /start.sh

VOLUME ["/secrets"]

EXPOSE 80 443

ENTRYPOINT /start.sh
