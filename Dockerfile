FROM debian:wheezy
MAINTAINER Thor Sideburns
RUN add-apt-repository ppa:duh/golang && apt-get update && apt-get install -y golang
ADD thingieboab /var/server/thingieboab.go

EXPOSE 7474
CMD ["go", "run", "/var/server/thingieboab.go
