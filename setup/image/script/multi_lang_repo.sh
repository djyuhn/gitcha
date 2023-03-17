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

git remote add origin https://github.com/someuser/repo.git

# Create main branch and file
git checkout -b main

# Add C# Files
mv ../language_samples/C#/ .

git add .

git commit -q -m "Adding C# files" --author="FirstName1 LastName1 <gitcha1@gitcha.com>" --date="2022-11-10T08:00:00-06:00"

# Add Clojure Files
mv ../language_samples/Clojure/ .

git add .

git commit -q -m "Adding Clojure files" --author="FirstName1 LastName1 <gitcha1@gitcha.com>" --date="2022-11-10T08:10:00-06:00"

# Add Elixir Files
mv ../language_samples/Elixir/ .

git add .

git commit -q -m "Adding Elixir files" --author="FirstName2 LastName2 <gitcha2@gitcha.com>" --date="2022-11-20T08:20:00-06:00"

# Add Go Files
mv ../language_samples/Go/ .

git add .

git commit -q -m "Adding Go files" --author="FirstName2 LastName2 <gitcha2@gitcha.com>" --date="2022-11-20T08:30:00-06:00"

# Add Rust Files
mv ../language_samples/Rust/ .

git add .

git commit -q -m "Adding Rust files" --author="FirstName3 LastName3 <gitcha3@gitcha.com>" --date="2022-11-20T08:40:00-06:00"