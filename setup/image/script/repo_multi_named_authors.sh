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

cat > go.mod << MOD
module gitchatestrepo
MOD

cat > LICENSE << '__LICENSE__'
MIT License
Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:
The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.
THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
__LICENSE__

git add go.mod
git add LICENSE

# Create 1 commit for gitcha1@gitcha.com with different author names
git commit -q -m commit1 --author="Author1 Alias1 <gitcha1@gitcha.com>" --date="2022-11-10T08:00:00-06:00"

# Create 2 commits for gitcha2@gitcha.com with different author names
touch code.go
echo 'hello world' >> code.go
git add code.go
git commit -q -m commit2 --author="Author2 Alias1 <gitcha2@gitcha.com>" --date="2022-11-10T08:10:00-06:00"

touch app.go
echo 'some code' >> app.go
git add app.go
git commit -q -m commit3 --author="Author2 Alias1 <gitcha2@gitcha.com>" --date="2022-11-10T08:20:00-06:00"

# Create 3 commits for gitcha3@gitcha.com with different author names
touch root.go
echo 'some code' >> root.go
git add root.go
git commit -q -m commit4 --author="Author3 Alias1 <gitcha3@gitcha.com>" --date="2022-11-10T08:30:00-06:00"

echo 'additional code' >> root.go
git add root.go
git commit -q -m commit5 --author="Author3 Alias2 <gitcha3@gitcha.com>" --date="2022-11-10T08:40:00-06:00"

echo 'more code' >> root.go
git add root.go
git commit -q -m commit6 --author="Author3 Alias3 <gitcha3@gitcha.com>" --date="2022-11-10T08:50:00-06:00"

# Create 4 commits for gitcha4@gitcha.com with different author names
touch root.go
echo 'some code' >> root.go
git add root.go
git commit -q -m commit7 --author="Author4 Alias1 <gitcha4@gitcha.com>" --date="2022-11-10T09:00:00-06:00"

echo 'additional code' >> root.go
git add root.go
git commit -q -m commit8 --author="Author4 Alias2 <gitcha4@gitcha.com>" --date="2022-11-10T09:10:00-06:00"

echo 'more code' >> root.go
git add root.go
git commit -q -m commit9 --author="Author4 Alias3 <gitcha4@gitcha.com>" --date="2022-11-10T09:20:00-06:00"

echo 'some code' >> root.go
git add root.go
git commit -q -m commit10 --author="Author4 Alias4 <gitcha4@gitcha.com>" --date="2022-11-10T09:30:00-06:00"
