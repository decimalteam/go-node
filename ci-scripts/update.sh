#!/bin/sh
DECIMALDIR="$HOME/go/src/bitbucket.org/decimalteam/go-node"
DECIMALGIT="git@bitbucket.org:decimalteam/go-node.git"
BRANCH="develop"
printf "Decimal folder: %s\n" "$DECIMALDIR"
if [ ! -d $DECIMALTEAM ] ; then
  printf "No Decimal found. Cloning...\n"
  mkdir -p "$DECIMALDIR"
  cd "$DECIMALDIR" || exit
  git clone --branch "$BRANCH" "$DECIMALGIT" .
else
  printf "Pulling new version..."
  cd "$DECIMALDIR" || exit
  git pull origin "$BRANCH"
fi