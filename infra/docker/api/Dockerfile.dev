FROM golang:1.7

# Set the default timezone to EST.
ENV TZ=America/New_York
RUN echo $TZ | tee /etc/timezone \
	&& dpkg-reconfigure --frontend noninteractive tzdata

RUN apt-get update && apt-get -q -y install \
	git openssl apt-transport-https ca-certificates curl g++ gcc libc6-dev make pkg-config \
	libssl-dev cmake

# Build libssh2 from source
RUN cd $HOME && curl -fsSL https://github.com/libssh2/libssh2/archive/libssh2-1.7.0.tar.gz -o libssh2.tar.gz \
  && mkdir libgit2 \
 	&& tar xvf libssh2.tar.gz -C libgit2 \
	&& ls -la libgit2 \
	&& cd libgit2/libssh2-libssh2-1.7.0 \
	&& cmake -DBUILD_SHARED_LIBS=ON . \
	&& cmake --build . \
	&& make \
	&& make install \
	&& ldconfig

# Build libgit2 from source
RUN cd $HOME && curl -fsSL https://github.com/libgit2/libgit2/archive/v0.24.1.tar.gz -o v0.24.1.tar.gz \
 	&& tar xvf v0.24.1.tar.gz -C libgit2 \
	&& cd libgit2/libgit2-0.24.1 \
	&& cmake -DCURL=OFF . \
	&& cmake --build . \
	&& make \
	&& make install \
	&& ldconfig \
	&& rm -rf $HOME/libgit2

# Copy in requisite files.
COPY ./api /go/src/github.com/gophr-pm/gophr/api
COPY ./infra /go/src/github.com/gophr-pm/gophr/infra
COPY ./lib /go/src/github.com/gophr-pm/gophr/lib

# Build source and move things around.
RUN cd /go/src/github.com/gophr-pm/gophr/api \
    && echo -e "\nFetching dependencies...\n" \
    && go get -d -v \
    && echo -e "\nBuilding the binary...\n" \
    && go build -v -o gophr-api-binary \
    && chmod +x ./gophr-api-binary \
    && echo -e "\nMoving things around...\n" \
    && mkdir /gophr \
    && mv ./gophr-api-binary /gophr/gophr-api-binary \
    && mv ../infra/scripts/wait-for-it.sh /gophr/wait-for-it.sh \
    && cd /gophr \
    && rm -rf /go

# Since its no longer necessary, remove git.
RUN apt-get -y remove git

# Set the environment variables
ENV PORT="3000"
ENV GOPHR_ENV="dev"
ENV GOPHR_DB_ADDR="db-svc"
ENV GOPHR_SECRETS_PATH="/secrets"

VOLUME ["/secrets"]

# Set the generated directory as the work directory.
WORKDIR /gophr

# Wait for db:9042 with no timeout, then start the binary.
CMD ./wait-for-it.sh -h "$GOPHR_DB_ADDR" -p 9042 -t 0 \
    -- \
    ./gophr-api-binary --port "$PORT"
