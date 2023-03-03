package main

import (
	"github.com/urfave/cli/v2"
	"github.com/zbysir/lookip/internal/lib/public_ip"
	dns2 "github.com/zbysir/lookip/internal/pkg/dns"
	"github.com/zbysir/lookip/internal/pkg/dns/alidns"
	"github.com/zbysir/lookip/internal/pkg/dns/cloudflare"
	"github.com/zbysir/lookip/internal/pkg/signal"
	"github.com/zbysir/lookip/internal/worker"
	"log"
	"os"
	"strings"
	"sync"
)

func main() {
	c := cli.NewApp()
	c.Name = "lookip"
	c.Usage = ""
	c.Version = ""
	c.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:     "name",
			Usage:    "e.g. '*.domain.com' or 'domain.com'",
			EnvVars:  []string{"NAME"},
			Required: true,
		},
		&cli.StringFlag{
			Name:     "domain",
			Usage:    "e.g. domain.com",
			EnvVars:  []string{"DOMAIN"},
			Required: false,
		},
		&cli.StringFlag{
			Name:     "dns",
			Usage:    "'ali' or 'cloudflare'",
			EnvVars:  []string{"DNS"},
			Required: false,
		},
		&cli.StringFlag{
			Name:     "access-key-id",
			Usage:    "ali AccessKeyID",
			EnvVars:  []string{"ACCESS_KEY_ID"},
			Required: false,
		},
		&cli.StringFlag{
			Name:     "access-key-secret",
			Usage:    "ali AccessKeySecret",
			EnvVars:  []string{"ACCESS_KEY_SECRET"},
			Required: false,
		},
		&cli.StringFlag{
			Name:     "cf_token",
			Usage:    "cloudflare token",
			EnvVars:  []string{"CF_TOKEN"},
			Required: false,
		},
		&cli.StringFlag{
			Name:     "cf_zone_id",
			Usage:    "cloudflare zone id",
			EnvVars:  []string{"CF_ZONE_ID"},
			Required: false,
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
			cfZoneId := c.String("cf_zone_id")
			cfToken := c.String("cf_token")
			key := c.String("access-key-id")
			secret := c.String("access-key-secret")
			name := c.String("name")
			ipGetter := c.String("ip-getter")
			dnsType := c.String("dns")

			var dnss []dns2.DNS

			if dnsType != "" {
				for _, d := range strings.Split(dnsType, ",") {
					if d != "ali" && d != "cloudflare" {
						log.Fatalf("dns type `%s` not support", d)
					}

					switch d {
					case "cloudflare":
						dnss = append(dnss, cloudflare.NewDNS(cfToken, cfZoneId, name))
					default:
						dnss = append(dnss, alidns.NewAliDns(regionId, key, secret, domain, name))
					}
				}
			} else {
				dnss = append(dnss, alidns.NewAliDns(regionId, key, secret, domain, name))
			}

			g := public_ip.Factory(ipGetter)

			log.Printf("use `%s` to get ip", g.Name())
			log.Printf("use `%s` to update dns", dnsType)

			w := worker.NewIpWorker(dns2.NewCombinedDNS(dnss), g)
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
