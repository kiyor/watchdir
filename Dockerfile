FROM golang
ADD . /go/src/github.com/kiyor/watchdir
RUN cd /go/src/github.com/kiyor/watchdir && \
	go get && \
	go install github.com/kiyor/watchdir

ENTRYPOINT /go/bin/watchdir
