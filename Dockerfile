FROM frolvlad/alpine-glibc
COPY kademliasrv /
ENV GIN_MODE=release
# REST port
ENV PORT=8080
EXPOSE 8080
# DHT port
EXPOSE 1200
ENTRYPOINT ["/kademliasrv"]
