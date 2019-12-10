FROM registry.access.redhat.com/ubi8/ubi as production
WORKDIR $HOME/go/src/fcos-pinger-backend
COPY . .

ENV GOPATH=$HOME/go \
    GOBIN=$HOME/go/bin \
    GOPINGER=$HOME/go/src/fcos-pinger-backend \
    PATH=$HOME/go/bin:$PATH

# Install golang
RUN yum install -y golang curl git && \
    mkdir -p $GOPATH && \
    yum clean all

# Install dep - the dependency management tool for golang
RUN mkdir -p $GOBIN $GOPINGER && curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh

# Install mongodb-go-driver
RUN dep ensure

# Compile the code
RUN go build ./main.go

# Run the server
CMD ./main --debug
