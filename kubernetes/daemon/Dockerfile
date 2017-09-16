# build stage
FROM golang:1.8.3-stretch AS build-env
ADD . /src
RUN go get -u k8s.io/apimachinery/pkg/apis/meta/v1 && \
    go get -u k8s.io/client-go/kubernetes && \
    go get -u k8s.io/client-go/rest && \
    go get -u github.com/John-Lin/ovsdbDriver && \
    go get -u github.com/sirupsen/logrus && \
    go get -u github.com/containernetworking/cni/pkg/types 
RUN cd /src && go build -o flaxd
 
# final stage
FROM debian:stretch-slim
WORKDIR /app
COPY --from=build-env /src/flaxd /app/
ENTRYPOINT ./flaxd
