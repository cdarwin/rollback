package main

import (
	"encoding/json"
	"fmt"
	"github.com/codegangsta/cli"
	"os"
	"path/filepath"
)

func findJars(path string) []string {
	var slice []string
	filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		f, err := filepath.Glob(path + "/*exec-war.jar")
		if err != nil {
			return nil
		}
		for _, i := range f {
			slice = append(slice, i)
		}
		return nil
	})
	return slice
}

func matchCache(artifact string) string {
  return "war.cache"
}

func printManifest(services []string) {
	for _, jar := range services {
		base := filepath.Base(filepath.Dir(filepath.Dir(jar)))
		m := make(map[string][]string)
		m[base] = append(m[base], jar)
		// find cache
		m[base] = append(m[base], matchCache(jar))
		s, _ := json.MarshalIndent(m, "", " ")
		fmt.Println(string(s))
	}
}

func main() {
	app := cli.NewApp()
	app.Name = "genman"
	app.Usage = "Generate a manifest of artifacts"
	app.Version = "0.0.1"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "root, r",
			Value: "/home/fubar/artifacts/0.0.1-SNAPSHOT",
			Usage: "root directory for active jars",
		},
		cli.StringFlag{
			Name:  "cache, c",
			Value: "/opt/cache",
			Usage: "cache directory for shared jars",
		},
	}
	app.Action = func(c *cli.Context) {
		root := c.String("root")
		if _, err := os.Stat(root); os.IsNotExist(err) {
			fmt.Printf("No such file or directory: %s\n", root)
		}
		cache := c.String("cache")
		if _, err := os.Stat(cache); os.IsNotExist(err) {
			fmt.Printf("No such file or directory: %s\n", cache)
		}
		printManifest(findJars(root))
	}
	app.Run(os.Args)
}
