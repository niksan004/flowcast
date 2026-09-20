# this container is used for testing only
FROM alpine:latest

RUN apk add --no-cache openssh

RUN ssh-keygen -A

COPY key.pub /root/.ssh/authorized_keys

ENTRYPOINT ["/usr/sbin/sshd", "-D"]
