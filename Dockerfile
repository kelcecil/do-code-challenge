FROM ubuntu:latest
RUN apt-get update && \
  apt-get install -y curl tar make && \
  rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*
COPY . /gopath/src/github.com/kelcecil/do-code-challenge
ENV GOPATH /gopath
RUN curl -O https://storage.googleapis.com/golang/go1.6.linux-amd64.tar.gz && \
  tar -xvf go1.6.linux-amd64.tar.gz && \
  mv go /usr/local && \
  export PATH=${PATH}:/usr/local/go/bin/ && \
  cd /gopath/src/github.com/kelcecil/do-code-challenge && \
  make && make install && \
  rm -rf /usr/local/go
CMD ["index_server"]
