#!/bin/bash
echo "Current node is not a validator."
rm -r ~/.decimal
echo "Restating daemon."
sudo systemctl restart decd