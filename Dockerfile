FROM debian:wheezy
MAINTAINER Thor Sideburns
ADD thingieboab /var/server/thingieboab
ADD bobbybot.json /var/server/bobbybot.json

EXPOSE 7474
CMD ["/var/server/thingieboab"]
