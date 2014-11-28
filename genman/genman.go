package main

import (
	"encoding/json"
	"flag"
	"fmt"
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
	var root, cache string
	flag.StringVar(&root, "r", "/home/fubar/artifacts/0.0.1-SNAPSHOT", "root directory for active jars")
	flag.StringVar(&cache, "c", "/opt/cache", "cache directory for shared jars")
	flag.Usage = func() {
		fmt.Printf("Usage of %s:\n", os.Args[0])
		fmt.Printf("  %s -r /root/directory -c /op/cache\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()

	if _, err := os.Stat(root); os.IsNotExist(err) {
		fmt.Printf("No such file or directory: %s\n", root)
		return
	}

	if _, err := os.Stat(cache); os.IsNotExist(err) {
		fmt.Printf("No such file or directory: %s\n", cache)
		return
	}

	printManifest(findJars(root))
}
