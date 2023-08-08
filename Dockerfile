FROM golang:1.20

EXPOSE 8080

# Update packages
RUN apt-get update

# [certs]
RUN apt-get -y install ca-certificates

# [curl]
RUN apt-get -y install curl

# [copy binaries]
COPY ./bin/web /usr/local/bin/ 

WORKDIR /usr/local/bin/ 

ENTRYPOINT [ "web" ]
