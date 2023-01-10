FROM alpine:latest

WORKDIR /data

COPY script/empty_repo.sh .

RUN  apk add git git-daemon && \
     cd /data && \
     sh empty_repo.sh

CMD ["git", "daemon", "--export-all", "--verbose", "--base-path=/data", "--informative-errors", "--reuseaddr", "--listen=0.0.0.0", "--port=9418"]
