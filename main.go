package main

import (
	"context"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/gabu/moon"
	"github.com/olekukonko/tablewriter"
	"github.com/urfave/cli"
)

var supportedExchanges []string = []string{
	"poloniex",
	"bittrex",
	"cryptopia",
	"liqui",
	"bitfinex",
	"bitgrail",
}

type Balance struct {
	Symbol   string
	Exchange string
	Balance  moon.Balance
}

func main() {
	app := cli.NewApp()
	app.Name = "coinfolio"
	app.Usage = "aggregate your coin balances from multiple coin exchanges"
	app.Version = "1.0.0"

	flags := []cli.Flag{}
	for _, exchange := range supportedExchanges {
		flags = append(flags, cli.StringFlag{
			Name:  exchange,
			Usage: "api key and secret for " + exchange + " (key:secret)",
		})
	}
	app.Flags = flags

	app.Action = func(c *cli.Context) error {
		if c.NumFlags() == 0 {
			return cli.ShowAppHelp(c)
		} else {
			return aggs(c)
		}
	}

	app.Run(os.Args)
}

func aggs(c *cli.Context) error {
	var err error
	balances := []Balance{}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	for _, exchange := range supportedExchanges {
		balances, err = getBalances(ctx, c, exchange, balances)
		if err != nil {
			return err
		}
	}

	sort.Slice(balances, func(i int, j int) bool {
		return balances[i].Balance.BtcValue > balances[j].Balance.BtcValue
	})

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Exchange", "Symbol", "Value", "BTC Value"})

	total := 0.0
	for _, b := range balances {
		v, _ := strconv.ParseFloat(b.Balance.BtcValue, 64)
		total += v
		table.Append([]string{b.Exchange, b.Symbol, b.Balance.Amount, b.Balance.BtcValue})
	}

	table.SetFooter([]string{"", "Total", "", moon.FormatFloat(total) + " BTC"})
	table.Render()

	return nil
}

func getBalances(ctx context.Context, c *cli.Context, exchange string, balances []Balance) ([]Balance, error) {
	if s := c.String(exchange); s != "" {
		key, secret, err := parseKey(s)
		if err != nil {
			return balances, err
		}
		ex := newExchange(exchange)
		bs, err := ex.GetBalances(ctx, key, secret)
		if err != nil {
			return balances, newExchangeError(exchange, err)
		}
		for s, b := range *bs {
			balances = append(balances, Balance{
				Symbol:   s,
				Exchange: exchange,
				Balance:  b,
			})
		}
	}
	return balances, nil
}

func parseKey(s string) (string, string, error) {
	ss := strings.Split(s, ":")
	if len(ss) != 2 {
		return "", "", newAPIKeySecretFormatError(s)
	}
	key, secret := ss[0], ss[1]
	if key == "" || secret == "" {
		return key, secret, newAPIKeySecretFormatError(s)
	}
	return key, secret, nil
}

func newExchange(name string) moon.Exchange {
	switch name {
	case "poloniex":
		return &moon.Poloniex{}
	case "bittrex":
		return &moon.Bittrex{}
	case "cryptopia":
		return &moon.Cryptopia{}
	case "liqui":
		return &moon.Liqui{}
	case "bitfinex":
		return &moon.Bitfinex{}
	case "bitgrail":
		return &moon.Bitgrail{}
	}
	return nil
}

func newAPIKeySecretFormatError(s string) error {
	return cli.NewExitError(s+" is invalid format (e.g. --bittrex KEY:SECRET)", 1)
}

func newExchangeError(exchange string, e error) error {
	return cli.NewExitError("Faild connect to exchange: "+exchange+", error: "+e.Error(), 1)
}
