#
# For production:
#   docker build --target production -t phantom .
#   docker run -it --rm -v $HOME/.phantom:/home/phantom/.phantom phantom 
#
# For unit testing:
#   docker build --target test .
#   docker build --target test --build-arg GO_TESTS_FLAGS=--skip-generate .
#

# STAGE: base
## Compiles Phantom for use
FROM golang:latest AS base

### Base packages
RUN apt-get update --fix-missing && apt-get -y install \
    git build-essential zlib1g zlib1g-dev wget zip unzip

### Add phantom user
RUN groupadd -g 999 phantom && useradd -r -u 999 -g phantom phantom
RUN mkdir -p /home/phantom/ && chown -R phantom:phantom /home/phantom

### Build phantom:
RUN mkdir -p /go/src/github.com/cryptdefender3232/phantom
WORKDIR /go/src/github.com/cryptdefender3232/phantom
ADD . /go/src/github.com/cryptdefender3232/phantom/
RUN make
RUN cp -vv phantom-server /opt/phantom-server 

# STAGE: test
## Run unit tests against the compiled instance
## Use `--target test` in the docker build command to run this stage
FROM base AS test

ARG GO_TESTS_FLAGS=""
ENV GO_TESTS_FLAGS="${GO_TESTS_FLAGS}"

RUN apt-get update --fix-missing \
    && apt-get -y upgrade \
    && apt-get -y install \
    curl

RUN /opt/phantom-server unpack --force 

### Run unit tests
RUN /go/src/github.com/cryptdefender3232/phantom/go-tests.sh ${GO_TESTS_FLAGS}

# STAGE: production
## Final dockerized form of Phantom
FROM debian:bookworm-slim AS production

### Install production packages
RUN apt-get update --fix-missing \
    && apt-get -y upgrade \
    && apt-get -y install \
    libxml2 libxml2-dev libxslt-dev locate gnupg \
    libreadline6-dev libcurl4-openssl-dev git-core \
    libssl-dev libyaml-dev openssl autoconf libtool \
    ncurses-dev bison curl xsel postgresql \
    postgresql-contrib postgresql-client libpq-dev \
    curl libapr1 libaprutil1 libsvn1 \
    libpcap-dev libsqlite3-dev libgmp3-dev \
    nasm

### Install MSF for stager generation
RUN curl https://raw.githubusercontent.com/rapid7/metasploit-omnibus/master/config/templates/metasploit-framework-wrappers/msfupdate.erb > msfinstall \
    && chmod 755 msfinstall \
    && ./msfinstall \
    && mkdir -p ~/.msf4/ \
    && touch ~/.msf4/initial_setup_complete 

### Cleanup unneeded packages
RUN apt-get remove -y curl gnupg \
    && apt-get autoremove -y \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*

### Add phantom user
RUN groupadd -g 999 phantom \
    && useradd -r -u 999 -g phantom phantom \
    && mkdir -p /home/phantom/ \
    && chown -R phantom:phantom /home/phantom \
    && su -l phantom -c 'mkdir -p ~/.msf4/ && touch ~/.msf4/initial_setup_complete'

### Copy compiled binary
COPY --from=base /opt/phantom-server  /opt/phantom-server

### Unpack Phantom:
USER phantom
RUN /opt/phantom-server unpack --force 

WORKDIR /home/phantom/
VOLUME [ "/home/phantom/.phantom" ]
ENTRYPOINT [ "/opt/phantom-server" ]


# STAGE: production-slim (about 1Gb smaller)
FROM debian:bookworm-slim AS production-slim

### Install production packages
RUN apt-get update --fix-missing \
    && apt-get -y upgrade

### Cleanup unneeded packages
RUN apt-get autoremove -y \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*

### Add phantom user
RUN groupadd -g 999 phantom \
    && useradd -r -u 999 -g phantom phantom \
    && mkdir -p /home/phantom/ \
    && chown -R phantom:phantom /home/phantom

### Copy compiled binary
COPY --from=base /opt/phantom-server  /opt/phantom-server

### Unpack Phantom:
USER phantom
RUN /opt/phantom-server unpack --force 

WORKDIR /home/phantom/
VOLUME [ "/home/phantom/.phantom" ]
ENTRYPOINT [ "/opt/phantom-server" ]
