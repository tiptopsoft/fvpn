FROM alpine:3.15
## ADD, resolve: sh: xxx not found
RUN mkdir /lib64 && ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2
ADD fvpn /usr/local/bin/
CMD ["fvpn"]
