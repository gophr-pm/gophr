FROM cassandra:3

RUN set -x \
    && apt-get update \
    && apt-get install -y wget maven unzip openjdk-8-jdk --no-install-recommends \
    && (mkdir lucene-index-tmp \
      && cd lucene-index-tmp \
      && wget https://github.com/Stratio/cassandra-lucene-index/archive/3.0.4.0.zip \
      && unzip 3.0.4.0.zip \
      && cd cassandra-lucene-index-3.0.4.0 \
      && update-ca-certificates -f \
      && mvn clean package \
      && cp plugin/target/cassandra-lucene-index-plugin-*.jar /usr/share/cassandra/lib \
      && cd ../.. \
      && rm -rf lucene-index-tmp) \
    && apt-get purge -y --auto-remove wget maven unzip openjdk-8-jdk \
    && rm -rf /var/lib/apt/lists/*

COPY ./files/docker-entrypoint.sh /docker-entrypoint.sh
COPY ./files/create-schema.sh /create-schema.sh
COPY ./files/schema.cql /schema.cql
ENTRYPOINT ["/docker-entrypoint.sh"] 
CMD ["cassandra", "-f"]
