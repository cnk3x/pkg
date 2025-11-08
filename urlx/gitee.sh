#!/usr/bin/env bash

set -eu

cd $(dirname $0)

rm -rf .temp/gitee
mkdir -p .temp/gitee
cd .temp/gitee
git clone git@gitee.com:k3x/urlx.git .
rm -rf *
rm -rf .vscode .gitignore

cd ../../
cp -r * .temp/gitee
cp .gitignore .temp/gitee

cd .temp/gitee
rm gitee.sh
find . \( -name '*.go' -o -name 'go.mod' -o -name 'go.sum' \) -exec sed -i '' 's@github.com/cnk3x@gitee.com/k3x@g' {} \;

git add .
git commit -m 'sync'
git push

cd ../../
rm -rf .temp/gitee
