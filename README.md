## Token
Darwinia token supply 

### require

1. golang >= 1.18.1



### install
```shell script
go build -o token
```

### run

```shell script
go run main.go
```


### docker

```shell script
docker-compose up -d
```

### Api doc

mainnet host: http://api.darwinia.network/

#### ring-supply

`Get /supply/ring`

### Example Response

`200 OK` and
```json
{
  "code": 0,
  "data": {
    "circulatingSupply": "449021080.3793598283033951",
    "totalSupply": "2015424739.889267952",
    "maxSupply": "10000000000",
    "details": [
      {
        "network": "Tron",
        "circulatingSupply": "42463320.5011454706097711",
        "totalSupply": "90403994.9525478491788821",
        "precision": 18,
        "type": "trc20",
        "contract": "TL175uyihLqQD656aFx3uhHYe1tyGkmXaW"
      },
      {
        "network": "Ethereum",
        "circulatingSupply": "406557759.878214357693624",
        "totalSupply": "1031251135.065152737693624",
        "precision": 18,
        "type": "erc20",
        "contract": "0x9469d013805bffb7d3debe5e7839237e535ec483"
      }
    ]
  },
  "msg": "ok"
}

```

-----


#### kton-supply

`Get /supply/kton`

### Example Response

`200 OK` and
```json
{
  "code": 0,
  "data": {
    "circulatingSupply": "53363.7233688935671044",
    "totalSupply": "68021.225215375",
    "maxSupply": "53363.7233688935671044",
    "details": [
      {
        "network": "Tron",
        "circulatingSupply": "1355.418652992761802",
        "totalSupply": "1355.418652992761802",
        "precision": 18,
        "type": "trc20",
        "contract": "TW3kTpVtYYQ5Ka1awZvLb9Yy6ZTDEC93dC"
      },
      {
        "network": "Ethereum",
        "circulatingSupply": "52008.3047159008053024",
        "totalSupply": "52008.3047159008053024",
        "precision": 18,
        "type": "erc20",
        "contract": "0x9f284e1337a815fe77d2ff4ae46544645b20c5ff"
      }
    ]
  },
  "msg": "ok"
}

```

-----



