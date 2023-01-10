#!/bin/bash
set -eu -o pipefail

mkdir testdata && cd testdata

git init -q

# Name and Email necessary for recognition
git config --local --add "committer.name" "gitcha-committer-name"
git config --local --add "committer.email" "gitcha-committer-email@gitcha.com"

git config --local --add "user.name" "gitcha-user-name"
git config --local --add "user.email" "gitcha-user-email@gitcha.com"

git config --local --add "author.name" "gitcha-author-name"
git config --local --add "author.email" "gitcha-author-email@gitcha.com"