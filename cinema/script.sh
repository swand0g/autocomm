#!/bin/bash

set -e

export PATH=$PATH_ADDITION:$PATH

PROMPT_PREFIX="%F{220}⚡>%f "

function type() {
  printf $PROMPT_PREFIX
  echo -n " "
  echo $* | node ./type-letters.js
}

type 'git status'
git status

type 'git add .'
git add .

type 'autocomm'
autocomm

type 'git status'
git status

type git push
git push

sleep 2
echo ""