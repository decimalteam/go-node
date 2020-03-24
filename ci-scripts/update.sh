#!/bin/sh

DECIMALDIR="$HOME/go/src/bitbucket.org/decimalteam/go-node"
DECIMALGIT="git@bitbucket.org:decimalteam/go-node.git"
BRANCH="develop"
SERVICEPATH="/etc/systemd/system/decd.service"
DECDBIN="$HOME/go/bin/decd"
printf "Decimal folder: %s\n" "$DECIMALDIR"

if [ ! -d $DECIMALDIR ] ; then
  printf "No Decimal found. Cloning...\n"
  mkdir -p $DECIMALDIR
  cd $DECIMALDIR || exit
  git clone --branch $BRANCH $DECIMALGIT .
  printf "Cloned. Building...\n"
  make all
else
  printf "Pulling new version..."
  cd $DECIMALDIR || exit
  git pull origin $BRANCH
  printf "Pulled. Building...\n"
  make all
fi

printf "Restarting service...\n"

if [ ! -d $SERVICEPATH ] ; then
  printf "No service file found. Creating.\n"
  sudo touch $SERVICEPATH
  echo "[Unit]
Description=Decimal daemon

[Service]
Type=simple
ExecStart=$DECDBIN start" | sudo tee $SERVICEPATH
  sudo systemctl daemon-reload
  sudo systemctl start decd
else
  sudo systemctl restart decd
fi

