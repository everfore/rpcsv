FROM golang

WORKDIR /app/gopath
ENV GOPATH /app/gopath
ADD . /app/gopath/

RUN go get github.com/shurcooL/github_flavored_markdown
RUN go build -o rpcsv

EXPOSE 80

CMD ["/app/gopath/rpcsv"]