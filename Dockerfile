FROM scratch

LABEL maintainer="lukapiske <lukapiske@gmail.com>"

ADD aloha /

EXPOSE 8090

CMD ["./aloha"]