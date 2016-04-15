FROM golang:1.6

# prepare working dir
RUN mkdir -p /go/src/app
WORKDIR /go/src/app

# copy code
COPY . /go/src/app
RUN go-wrapper download
RUN go-wrapper install

# set up cron dir
RUN mkdir -p /etc/gocron.d
VOLUME /etc/gocron.d
ADD example.yml /etc/gocron.d/example.yml

CMD ["go-wrapper", "run"]
