FROM ubuntu:latest
EXPOSE 8080
RUN apt-get update && \
  apt-get install -y curl tar make && \
  rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*
COPY . /src/
RUN curl -O https://storage.googleapis.com/golang/go1.6.linux-amd64.tar.gz && \
  tar -xvf go1.6.linux-amd64.tar.gz && \
  mv go /usr/local && \
  export PATH=${PATH}:/usr/local/go/bin/ && \
  cd /src && \
  make && make install && \
  rm -rf /usr/local/go
RUN useradd server
USER server
ENTRYPOINT ["index_server"]
