FROM ubuntu:20.04 as builder

RUN apt-get update

ENV PATH="/usr/local/go/bin:${PATH}"
ENV DEBIAN_FRONTEND='noninteractive'

RUN apt-get install -y libzmq3-dev \
                        wget \
                        pkg-config \
                        #gcc \
                        build-essential \
                        git

RUN wget https://golang.org/dl/go1.15.7.linux-amd64.tar.gz ; \
    tar -C /usr/local -xzf go1.15.7.linux-amd64.tar.gz ; \
    rm go1.15.7.linux-amd64.tar.gz

WORKDIR project

COPY . .

RUN make build

FROM ubuntu:20.04

RUN apt-get update
RUN apt-get install -y libzmq3-dev

COPY ./resources/ ./resources
COPY --from=builder ./project/signernode ./
COPY --from=builder ./project/bootstrap ./

CMD ["./signernode"]