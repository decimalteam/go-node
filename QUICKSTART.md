# Установка и запуск Decimal Go Node

## 0. Переменные окружения

```
export GOPATH=$HOME/go
export PATH=$PATH:$GOPATH/bin
export GO111MODULE=on
```

## 1. Загрузка

```
mkdir -p $GOPATH/src/bitbucket.org/decimalteam/go-node
cd $GOPATH/src/bitbucket.org/decimalteam/go-node
git clone git@bitbucket.org:decimalteam/go-node.git $GOPATH/src/bitbucket.org/decimalteam/go-node
```

## 2. Установка

```
make all
```

## 3. Конфигурация

```
cd 
mkdir -p .decimal/daemon/config/
```

Копируем в созданную папку `genesis.json` и `config.toml`

## 4. Запуск

```
decd start
```