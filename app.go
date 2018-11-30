package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/url"
	"regexp"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

type config struct {
	Secure    bool     `json:"secure"`
	CertFile  string   `json:"certFile"`
	KeyFile   string   `json:"keyFile"`
	Upstream  []string `json:"upstream"`
	Prohibido []string `json:"prohibido"`
	Port      int      `json:"port"`
}

func main() {
	e := echo.New()
	path := flag.String("path", "/etc/proxy-keos/proxy-keos.conf", "ruta del archivo de configuraciòn")
	flag.Parse()

	dat, err := ioutil.ReadFile(*path)

	if err != nil {
		e.Logger.Fatal(err)
	}
	config := config{}

	if err := json.Unmarshal(dat, &config); err != nil {
		e.Logger.Fatal(err)
	}

	targets := []*middleware.ProxyTarget{}

	for _, upstream := range config.Upstream {
		url1, err := url.Parse(upstream)
		if err != nil {
			e.Logger.Fatal(err)
		}
		targets = append(targets, &middleware.ProxyTarget{URL: url1})
	}

	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(context echo.Context) error {
			url := context.Request().URL.Path
			if config.Prohibido == nil {
				return next(context)
			}

			for _, pattern := range config.Prohibido {
				exp, err := regexp.Compile(pattern)
				if err != nil {
					e.Logger.Error(err)
				}
				if exp.Match([]byte(url)) {
					return echo.ErrNotFound
				}
			}

			return next(context)
		}
	})

	e.Use(middleware.Proxy(middleware.NewRoundRobinBalancer(targets)))

	port := fmt.Sprintf(":%v", config.Port)

	if config.Secure {
		certFile, keyFile := config.CertFile, config.KeyFile
		e.Logger.Fatal(e.StartTLS(port, certFile, keyFile))
		return
	}
	e.Logger.Fatal(e.Start(port))
}
