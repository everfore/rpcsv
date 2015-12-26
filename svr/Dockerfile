FROM golang

ENV GOPATH /app/gopath

RUN go get github.com/shurcooL/github_flavored_markdown
RUN go get github.com/everfore/rpcsv

WORKDIR /wkdir
RUN git clone --depth 1 git://github.com/everfore/rpcsv.git .
WORKDIR /wkdir/svr
RUN go build main.go

EXPOSE 8800

CMD ["/wkdir/svr/svr"]
