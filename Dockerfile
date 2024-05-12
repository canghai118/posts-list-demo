FROM golang:1.21

WORKDIR /build

ADD . /build

RUN make release

FROM alpine:3.12

COPY --from=0 /build/posts-list /bin/posts-list

ENTRYPOINT ["/bin/posts-list"]