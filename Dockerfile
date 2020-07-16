FROM ubuntu:latest

VOLUME /go/src/github.com/mikemackintosh/wonka
WORKDIR /go/src/github.com/mikemackintosh/wonka

RUN apt-get update \
  && apt-get install -fy\
    build-essential \
    autoconf \
    git \
    gettext \
    libtool \
    autopoint \
    libconfig-dev \
    login \
    pamtester \
    ca-certificates \
    curl \
    xsltproc \
    bison \
    flex \
    ruby-serverspec

RUN curl -s https://storage.googleapis.com/golang/go1.13.4.linux-amd64.tar.gz | tar -v -C /usr/local -xz

# Let's people find our Go binaries
ENV PATH $PATH:/usr/local/go/bin
ENV GOPATH /Users/splug/go

COPY testing/Gemfile Gemfile
RUN gem install bundler && \
    bundler install --jobs=3 --path=/vendor && \
    cp -R .bundle /

CMD ["echo", "Please choose a command to run"]
