FROM ubuntu:latest

RUN apt-get update

ENV PATH="/usr/local/go/bin:${PATH}"

RUN DEBIAN_FRONTEND='noninteractive' apt-get install -y libzmq3-dev \
                        wget \
                        pkg-config \
                        gcc \
                        git

RUN wget https://golang.org/dl/go1.15.7.linux-amd64.tar.gz ; \
    tar -C /usr/local -xzf go1.15.7.linux-amd64.tar.gz ; \
    rm go1.15.7.linux-amd64.tar.gz

WORKDIR project

COPY . .

RUN git config --global url."https://jffp113:d635c14a90e10ab8ddbf3a09a8ef8c16507ea08e@github.com".insteadOf https://github.com

RUN go build signernode.go
RUN go build main_bootstrap.go

CMD ["./signernode"]