#!/bin/bash
set -eu -o pipefail

shopt -s expand_aliases
source ~/.bashrc

SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )

docker build -f "$SCRIPT_DIR"/image/basic_repo_single_author.Dockerfile -t gitcha/basic_repo_single_author ./setup/image/
docker build -f "$SCRIPT_DIR"/image/basic_repo_multiple_authors.Dockerfile -t gitcha/basic_repo_multiple_authors ./setup/image/
