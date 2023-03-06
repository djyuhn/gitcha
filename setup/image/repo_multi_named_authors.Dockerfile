FROM alpine:latest

WORKDIR /data

COPY script/repo_multi_named_authors.sh .

RUN  apk add git git-daemon && \
     cd /data && \
     sh repo_multi_named_authors.sh

CMD ["git", "daemon", "--export-all", "--verbose", "--base-path=/data", "--informative-errors", "--reuseaddr", "--listen=0.0.0.0", "--port=9418"]
