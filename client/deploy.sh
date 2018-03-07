#!/bin/bash

NODE_ENV=production node_modules/.bin/webpack
git add public/bundle.js
git commit -m "Deploying"
git push
cd $GOPATH/src/github.com/ubclaunchpad/bumper
git subtree push --prefix client/public origin gh-pages
