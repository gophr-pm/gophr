FROM cassandra:3.7

RUN set -x \
	&& apt-get update \
	&& apt-get install -y --no-install-recommends \
			wget maven unzip openjdk-8-jdk \
	&& \
		(mkdir lucene-index-tmp \
	    && cd lucene-index-tmp \
	    && wget https://github.com/Stratio/cassandra-lucene-index/archive/3.7.1.zip -O archive.zip \
	    && unzip archive.zip \
	    && cd cassandra* \
	    && update-ca-certificates -f \
	    && mvn clean package \
	    && cp plugin/target/cassandra-lucene-index-plugin-*.jar /usr/share/cassandra/lib \
	    && cd ../.. \
	    && rm -rf lucene-index-tmp) \
	&& apt-get purge -y --auto-remove \
			wget maven unzip openjdk-8-jdk \
	&& rm -rf /var/lib/apt/lists/*
