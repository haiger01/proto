FROM centos:latest
RUN yum update -y && \
    yum install -y vim git wget
RUN wget https://dl.google.com/go/go1.13.6.linux-amd64.tar.gz 
RUN tar -C /usr/local/ -zxf go1.13.6.linux-amd64.tar.gz && \
    rm go1.13.6.linux-amd64.tar.gz && \
    # mkdir -p /usr/local/router && \
    mkdir -p /root/go/src/github.com/spectrex02/router-shakyo-go && \
    echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bash_profile && \
    echo 'export PATH=$PATH:$(go env GOPATH)/bin' >> ~/.bash_profile
RUN ["/bin/bash", "-c", "source ~/.bash_profile"]
WORKDIR /root/go/src/github.com/spectrex02/router-shakyo-go/
CMD ["/bin/bash"]
