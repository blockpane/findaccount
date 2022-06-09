# 1st stage, build app
FROM golang:1.18 as builder
RUN apt-get update && apt-get -y upgrade
COPY . /build/app
WORKDIR /build/app

RUN bash ./chains/fetch.sh
RUN go get ./... && go build -ldflags "-s -w" -o findaccount-server cmd/findaccount-server/main.go

# 2nd stage, create a user to copy, and install libraries needed if connecting to upstream TLS server
FROM debian:10 AS ssl
ENV DEBIAN_FRONTEND noninteractive
RUN apt-get update && apt-get -y upgrade && apt-get install -y ca-certificates && \
    addgroup --gid 26656 --system findaccount && adduser -uid 26656 --ingroup findaccount --system --home /var/lib/findaccount findaccount

# 3rd and final stage, copy the minimum parts into a scratch container, is a smaller and more secure build.
FROM scratch
COPY --from=ssl /etc/ca-certificates /etc/ca-certificates
COPY --from=ssl /etc/ssl /etc/ssl
COPY --from=ssl /usr/share/ca-certificates /usr/share/ca-certificates
COPY --from=ssl /usr/lib /usr/lib
COPY --from=ssl /lib /lib
COPY --from=ssl /lib64 /lib64

COPY --from=ssl /etc/passwd /etc/passwd
COPY --from=ssl /etc/group /etc/group
COPY --from=ssl --chown=findaccount:findaccount /var/lib/findaccount /var/lib/findaccount

COPY --from=builder /build/app/findaccount-server /findaccount-server

EXPOSE 8080
USER findaccount
WORKDIR /var/lib/findaccount

ENTRYPOINT ["/findaccount-server"]
