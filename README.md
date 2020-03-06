# Decimal Go Node

## Installation

- Install [Golang](https://golang.org/doc/install)

- Set Golang environment variables:
```shell script
export GO111MODULES=on
export GOPATH=<wherever you want>
```

- Add `$GOPATH/bin` to your PATH:
```shell script
export PATH="$PATH:$GOPATH/bin"
```

- Cloning repository
```shell script
mkdir -p $GOPATH/src/bitbucket.org/decimalteam
cd $GOPATH/src/bitbucket.org/decimalteam
git clone https://bitbucket.org/decimalteam/go-node.git
cd go-node
```
- Run `make all`

- `decd start`