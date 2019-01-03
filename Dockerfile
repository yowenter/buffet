FROM golang:1.10


RUN mkdir /app
RUN mkdir /go/src/app
ADD . /go/src/app
WORKDIR /go/src/app

RUN mkdir github.com/yowenter/buffet
RUN cp -rf pkg github.com/yowenter/buffet/


RUN curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh
RUN dep ensure --vendor-only

EXPOSE 5000


RUN go build -o /app/main .
CMD ["/app/main"]
