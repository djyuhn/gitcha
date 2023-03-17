FROM alpine:latest

WORKDIR /data

COPY script/multi_lang_repo.sh .
COPY language_samples ./language_samples

RUN  apk add git git-daemon && \
     cd /data && \
     sh multi_lang_repo.sh

CMD ["git", "daemon", "--export-all", "--verbose", "--base-path=/data", "--informative-errors", "--reuseaddr", "--listen=0.0.0.0", "--port=9418"]
