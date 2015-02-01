# Start from a Debian image with the latest version of Go installed
# and a workspace (GOPATH) configured at /go.
FROM golang

# Copy the local package files to the container's workspace.
ADD . /go/src/github.com/ninetwentyfour/go-imago

# Download and install wkhtmltopdf
RUN sed 's/main$/main universe/' -i /etc/apt/sources.list
RUN apt-get update
RUN apt-get upgrade -y

RUN apt-get install -y build-essential xorg libssl-dev libxrender-dev wget
RUN wget http://wkhtmltopdf.googlecode.com/files/wkhtmltoimage-0.11.0_rc1-static-amd64.tar.bz2
RUN tar xvjf wkhtmltoimage-0.11.0_rc1-static-amd64.tar.bz2
RUN install wkhtmltoimage-amd64 /usr/bin/wkhtmltoimage

# Build the outyet command inside the container.
# (You may fetch or manage dependencies here,
# either manually or with a tool like "godep".)
RUN go get github.com/asaskevich/govalidator
RUN go get github.com/gorilla/mux
RUN go get github.com/ninetwentyfour/go-wkhtmltoimage
RUN go get github.com/zenazn/goji/graceful
RUN go get gopkg.in/amz.v1/s3
RUN go get github.com/garyburd/redigo/redis
RUN go get github.com/nfnt/resize
RUN go install github.com/ninetwentyfour/go-imago

# Run the outyet command by default when the container starts.
ENTRYPOINT /go/bin/go-imago

# Document that the service listens on port 6001.
EXPOSE 6001
