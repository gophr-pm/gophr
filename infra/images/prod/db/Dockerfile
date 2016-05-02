# vim:set ft=dockerfile:
FROM debian:jessie-backports

# explicitly set user/group IDs
RUN groupadd -r cassandra --gid=999 && useradd -r -g cassandra --uid=999 cassandra

ENV GOSU_VERSION 1.7
ENV CASSANDRA_DATA /cassandra_data/data
ENV CASSANDRA_CONFIG /etc/cassandra
ENV CASSANDRA_VERSION 3.0.5
ENV CASSANDRA_DATA_HOME /cassandra_data

# grab gosu for easy step-down from root
RUN set -x \
  && apt-key adv --keyserver ha.pool.sks-keyservers.net --recv-keys 514A2AD631A57A16DD0047EC749D6EEC0353B12C \
  && echo 'deb http://www.apache.org/dist/cassandra/debian 30x main' >> /etc/apt/sources.list.d/cassandra.list \
	&& apt-get update \
  && apt-get install -y --no-install-recommends ca-certificates wget \
  && rm -rf /var/lib/apt/lists/* \
	&& wget -O /usr/local/bin/gosu "https://github.com/tianon/gosu/releases/download/$GOSU_VERSION/gosu-$(dpkg --print-architecture)" \
	&& wget -O /usr/local/bin/gosu.asc "https://github.com/tianon/gosu/releases/download/$GOSU_VERSION/gosu-$(dpkg --print-architecture).asc" \
	&& export GNUPGHOME="$(mktemp -d)" \
	&& gpg --keyserver ha.pool.sks-keyservers.net --recv-keys B42F6819007F00F88E364FD4036A9C25BF357DD4 \
	&& gpg --batch --verify /usr/local/bin/gosu.asc /usr/local/bin/gosu \
	&& rm -r "$GNUPGHOME" /usr/local/bin/gosu.asc \
	&& chmod +x /usr/local/bin/gosu \
	&& gosu nobody true \
  && apt-get install -y cassandra="$CASSANDRA_VERSION" \
  && mkdir -p "$CASSANDRA_DATA" "$CASSANDRA_CONFIG" \
  && chown -R cassandra.cassandra "$CASSANDRA_CONFIG" "$CASSANDRA_DATA_HOME" \
  && chmod o+w -R "$CASSANDRA_CONFIG" "$CASSANDRA_DATA_HOME" \
	&& apt-get purge -y --auto-remove ca-certificates wget \
  && rm -rf /var/lib/apt/lists/*

COPY ./files/cassandra-lucene-index-plugin*  /usr/share/cassandra/lib
COPY ./files/cassandra.yaml 					       /etc/cassandra/cassandra.yaml
COPY ./files/logback.xml 							       /etc/cassandra/logback.xml
COPY ./files/kubernetes-cassandra.jar        /kubernetes-cassandra.jar
COPY ./files/docker-entrypoint.sh 		       /docker-entrypoint.sh
COPY ./files/create-schema.sh 				       /create-schema.sh
COPY ./files/run-cassandra.sh 				       /run-cassandrae.sh
COPY ./files/schema.cql 							       /schema.cql

USER cassandra
VOLUME ["/cassandra_data/data"]
ENTRYPOINT ["/docker-entrypoint.sh"]

# 7000: intra-node communication
# 7001: TLS intra-node communication
# 7199: JMX
# 9042: CQL
# 9160: thrift service
EXPOSE 7000 7001 7199 9042 9160
CMD /run-cassandra.sh
