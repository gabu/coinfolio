package main

import (
	"context"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"

	"golang.org/x/sync/errgroup"

	"github.com/gabu/moon"
	"github.com/olekukonko/tablewriter"
	"github.com/urfave/cli"
)

var supportedExchanges = []string{
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

	flags = append(flags, cli.StringSliceFlag{
		Name:  "sort",
		Usage: `sort method, "exchange", "symbol", btc", "value" are available. (default: btc)`,
	})

	for _, exchange := range supportedExchanges {
		flags = append(flags, cli.StringSliceFlag{
			Name:  exchange,
			Usage: "api key and secret for " + exchange + " (key:secret)",
		})
	}
	app.Flags = flags

	app.Action = func(c *cli.Context) error {
		if c.NumFlags() == 0 {
			return cli.ShowAppHelp(c)
		}
		return aggs(c)
	}

	app.Run(os.Args)
}

func aggs(c *cli.Context) error {
	var mu sync.Mutex
	balances := []Balance{}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	wg, ctx := errgroup.WithContext(ctx)

	for _, exchange := range supportedExchanges {
		exchange := exchange
		wg.Go(func() error {
			b, err := getBalances(ctx, c, exchange)
			mu.Lock()
			defer mu.Unlock()
			balances = append(balances, b...)
			return err
		})
	}

	if err := wg.Wait(); err != nil {
		return err
	}

	sortm := c.StringSlice("sort")
	sort.Slice(balances, func(i int, j int) bool {
		ib, jb := balances[i], balances[j]
		for _, m := range sortm {
			if m == "btc" && ib.Balance.BtcValue != jb.Balance.BtcValue {
				return ib.Balance.BtcValue > jb.Balance.BtcValue
			}
			if m == "exchange" && ib.Exchange != jb.Exchange {
				return ib.Exchange > jb.Exchange
			}
			if m == "symbol" && ib.Symbol != jb.Symbol {
				return ib.Symbol > jb.Symbol
			}
			if m == "value" && ib.Balance.Amount != jb.Balance.Amount {
				return ib.Balance.Amount > jb.Balance.Amount
			}
		}
		return false
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

func getBalances(ctx context.Context, c *cli.Context, exchange string) ([]Balance, error) {
	ss := c.StringSlice(exchange)
	balances := make([]Balance, 0, len(ss))

	for i, s := range ss {
		if s == "" {
			continue
		}
		var exSuffix string
		if len(ss) > 1 {
			exSuffix = " #" + strconv.FormatInt(int64(i+1), 10)
		}
		key, secret, err := parseKey(s)
		if err != nil {
			return balances, err
		}
		ex := newExchange(exchange)
		bs, err := ex.GetBalances(ctx, key, secret)
		if err != nil {
			return balances, newExchangeError(exchange+exSuffix, err)
		}
		for s, b := range *bs {
			balances = append(balances, Balance{
				Symbol:   s,
				Exchange: exchange + exSuffix,
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
