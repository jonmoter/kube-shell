##############################################################
# Golang build image
##############################################################
FROM golang:1.11-alpine3.7

WORKDIR /tmp

COPY truth.go ./
RUN go build -o /tmp/truthserver truth.go

##############################################################
# Main Docker Image
##############################################################
FROM alpine:3.7

# Explicit is better than implicit
USER root
ENV HOME /root

WORKDIR $HOME

RUN apk add --no-cache \
      bash \
      bash-completion \
      bind-tools \
      busybox-extras \
      curl \
      dnstop \
      iproute2 \
      iptraf-ng \
      jq \
      net-tools \
      openrc \
      python3 \
      tcpdump \
      tcptraceroute \
      vim

RUN pip3 install --no-cache-dir --upgrade \
      pip \
      setuptools \
      httpie==0.9.9 \
      httpstat

# Use script to install kubectl, to avoid having a 50MB+ binary in the docker image
COPY install_kubectl /usr/local/bin/

COPY dotfiles/* $HOME/

RUN mkdir -p /usr/local/bin
COPY --from=0 /tmp/truthserver /usr/local/bin/truthserver

ENTRYPOINT ["/bin/sh", "-c"]
CMD ["bash"]
