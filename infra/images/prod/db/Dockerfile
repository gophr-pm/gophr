# vim:set ft=dockerfile:
FROM google/debian:jessie

COPY ./files/cassandra.list /etc/apt/sources.list.d/cassandra.list
COPY ./files/run-cassandra.sh /run-cassandra.sh

ENV CASSANDRA_DATA /cassandra_data/data
ENV CASSANDRA_CONFIG /etc/cassandra
ENV CASSANDRA_VERSION 3.0.5
ENV CASSANDRA_DATA_HOME /cassandra_data
# The below is important for the kubernetes extension of cassandra. This
# variabled **must** match the name of the db service in kubernetes verbatim.
ENV CASSANDRA_SERVICE gophr-db-svc

RUN gpg --keyserver pgp.mit.edu --recv-keys F758CE318D77295D && \
  gpg --export --armor F758CE318D77295D | apt-key add - && \
  gpg --keyserver pgp.mit.edu --recv-keys 2B5C1B00 && \
  gpg --export --armor 2B5C1B00 | apt-key add - && \
  gpg --keyserver pgp.mit.edu --recv-keys 0353B12C && \
  gpg --export --armor 0353B12C | apt-key add - && \
  apt-get update && \
  apt-get -qq -y install procps  cassandra openjdk-8-jre-headless && \
  chmod a+rx /run-cassandra.sh && \
  mkdir -p "$CASSANDRA_DATA" && \
  chown -R cassandra.cassandra /etc/cassandra /cassandra_data && \
  chmod o+w -R /etc/cassandra "$CASSANDRA_DATA_HOME" && \
  rm -rf /var/lib/apt/lists/* && \
  rm -rf /usr/share/doc/ && \
  rm -rf /usr/share/doc-base/ && \
  rm -rf /usr/share/man/ && \
  rm -rf /tmp/*

COPY ./files/cassandra-lucene-index-plugin*  /usr/share/cassandra/lib/
COPY ./files/cassandra.yaml                  /etc/cassandra/cassandra.yaml
COPY ./files/logback.xml 							       /etc/cassandra/logback.xml
COPY ./files/kubernetes-cassandra.jar        /kubernetes-cassandra.jar
COPY ./files/docker-entrypoint.sh 		       /docker-entrypoint.sh
COPY ./files/create-schema.sh 				       /create-schema.sh
COPY ./files/schema.cql 							       /schema.cql

VOLUME ["/cassandra_data/data"]

USER cassandra

ENTRYPOINT /docker-entrypoint.sh
CMD /run-cassandra.sh
