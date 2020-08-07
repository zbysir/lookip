package main

import (
	"github.com/urfave/cli/v2"
	"github.com/zbysir/lookip/internal/lib/public_ip"
	"github.com/zbysir/lookip/internal/pkg/signal"
	"github.com/zbysir/lookip/internal/worker"
	"log"
	"os"
	"sync"
)

func main() {
	c := cli.NewApp()
	c.Name = "lookip"
	c.Usage = ""
	c.Version = ""
	c.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:     "domain",
			Usage:    "domain, e.g. baidu.com",
			EnvVars:  []string{"DOMAIN"},
			Required: true,
		},
		&cli.StringFlag{
			Name:     "rr",
			Usage:    "e.g. www / *",
			EnvVars:  []string{"RR"},
			Required: true,
		},
		&cli.StringFlag{
			Name:     "access-key-id",
			Usage:    "ali AccessKeyID",
			EnvVars:  []string{"ACCESS_KEY_ID"},
			Required: true,
		},
		&cli.StringFlag{
			Name:     "access-key-secret",
			Usage:    "ali AccessKeySecret",
			EnvVars:  []string{"ACCESS_KEY_SECRET"},
			Required: true,
		},
		&cli.StringFlag{
			Name:     "region-id",
			Usage:    "regionId, e.g. zh-hangzhou",
			EnvVars:  []string{"REGION_ID"},
			Required: false,
			Value:    "zh-hangzhou",
		},
		&cli.StringFlag{
			Name:     "ip-getter",
			Usage:    "ip-getter, Can use values from the following array: [httpbin(httpbin.org}, 3322(3322.net)]",
			EnvVars:  []string{"IP_GETTER"},
			Required: false,
			Value:    "httpbin",
		},
	}
	c.Action = func(c *cli.Context) error {
		ctx, _ := signal.NewTermContext()

		var wg sync.WaitGroup

		wg.Add(1)
		go func() {
			defer wg.Done()
			regionId := c.String("region-id")
			domain := c.String("domain")
			key := c.String("access-key-id")
			secret := c.String("access-key-secret")
			rr := c.String("rr")
			ipGetter := c.String("ip-getter")

			g := public_ip.Factory(ipGetter)

			log.Printf("use `%s` to get ip", g.Name())

			w := worker.NewIpWorker(regionId, key, secret, domain, rr, g)
			w.LoopUpdateIp(ctx)
		}()

		wg.Wait()

		return nil
	}

	err := c.Run(os.Args)
	if err != nil {
		panic(err)
	}
}
