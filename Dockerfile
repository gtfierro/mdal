FROM ubuntu:xenial
MAINTAINER Gabe Fierro <gtfierro@eecs.berkeley.edu>

RUN apt-get -y update && apt-get install -y git libraptor2-dev libssl-dev
RUN apt-get clean && rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*

ENV TZ=America/Los_Angeles
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

#ADD hod /bin/hod
#ADD hodconfig.yaml /etc/hod/
#ADD entrypoint.sh /bin/
#ADD server /server
#ADD Brick.ttl BrickFrame.ttl /
ADD mdal /bin/

ENTRYPOINT [ "/bin/mdal", "start", "-c", "/etc/mdal/config.yml" ]

