# Start from a Debian image with the latest version of Go installed
# and a workspace (GOPATH) configured at /go.
FROM golang

# Copy the local package files to the container's workspace.
ADD . /go/src/github.com/ninetwentyfour/go-imago

# Download and install wkhtmltopdf
RUN DEBIAN_FRONTEND=noninteractive apt-get update && apt-get install -y fontconfig libjpeg62-turbo xfonts-base xfonts-75dpi xfonts-100dpi libx11-6 libxext6 libxrender1
ADD http://deis-deps.s3.amazonaws.com/wkhtmltox-0.12.2.1_linux-jessie-amd64.deb /tmp/
RUN dpkg -i /tmp/wkhtmltox-0.12.2.1_linux-jessie-amd64.deb

# Build the outyet command inside the container.
# (You may fetch or manage dependencies here,
# either manually or with a tool like "godep".)
RUN go get github.com/gorilla/mux
RUN go get github.com/ninetwentyfour/go-wkhtmltoimage
RUN go get github.com/zenazn/goji/graceful
RUN go get gopkg.in/amz.v1/s3
RUN go get github.com/garyburd/redigo/redis
RUN go get github.com/nfnt/resize
RUN go get github.com/yvasiyarov/gorelic
RUN go get github.com/yvasiyarov/newrelic_platform_go
RUN go get github.com/yvasiyarov/go-metrics
RUN go install github.com/ninetwentyfour/go-imago

# Run the outyet command by default when the container starts.
ENTRYPOINT /go/bin/go-imago

# Document that the service listens on port 6001.
EXPOSE 6001
