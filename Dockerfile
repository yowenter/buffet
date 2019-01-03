FROM golang:1.10



RUN mkdir -p /go/src/github.com/yowenter/buffet
ADD . /go/src/github.com/yowenter/buffet
WORKDIR /go/src/github.com/yowenter/buffet




RUN curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh
RUN dep ensure -v

EXPOSE 5000


RUN go build -o /main .
CMD ["/main"]
