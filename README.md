# coinfolio

Aggregate your coin balances from multiple coin exchanges.

## Supported exchanges

- poloniex
- bittrex
- cryptopia
- liqui
- bitfinex
- bitgrail

## Usage

### Build

```
git clone https://github.com/gabu/coinfolio.git
cd coinfolio
glide install
go install
```

### Run

```
$ coinfolio
NAME:
   coinfolio - aggregate your coin balances from multiple coin exchanges

USAGE:
   coinfolio [global options] command [command options] [arguments...]

VERSION:
   1.0.0

COMMANDS:
     help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --poloniex value   api key and secret for poloniex (key:secret)
   --bittrex value    api key and secret for bittrex (key:secret)
   --cryptopia value  api key and secret for cryptopia (key:secret)
   --liqui value      api key and secret for liqui (key:secret)
   --bitfinex value   api key and secret for bitfinex (key:secret)
   --bitgrail value   api key and secret for bitgrail (key:secret)
   --help, -h         show help
   --version, -v      print the version
```

### Example

```
coinfolio --bittrex YOUR_API_KEY:YOUR_API_SECRET  --liqui YOUR_API_KEY:YOUR_API_SECRET  --cryptopia YOUR_API_KEY:YOUR_API_SECRET
```

```
+-----------+--------+----------------+
| EXCHANGE  | SYMBOL |   BTC VALUE    |
+-----------+--------+----------------+
| bittrex   | WAVES  |     0.96533102 |
| bittrex   | UBQ    |     0.91293400 |
| bittrex   | GBYTE  |     0.50271003 |
| liqui     | PTOY   |     0.21913528 |
| liqui     | PLU    |     0.12250279 |
| cryptopia | INSN   |     0.09524800 |
| liqui     | SNM    |     0.06564117 |
| liqui     | MGO    |     0.05853929 |
| cryptopia | INPAY  |     0.04437000 |
| liqui     | ETH    |     0.00828395 |
| liqui     | BTC    |     0.00187128 |
| bittrex   | XLM    |     0.00050415 |
| bittrex   | BTC    |     0.00000332 |
+-----------+--------+----------------+
|             TOTAL  | 2 99707428 BTC |
+-----------+--------+----------------+
```

## LICENSE

MIT
