package main

import (
	"log"
	"net/http"
	//"github.com/plentymarkets/status/pkg/status"
	"github.com/plentymarkets/version/pkg/version"
)

func main() {
	http.HandleFunc("/plugin-git-version/version", version.GetVersions)
	log.Fatal(http.ListenAndServe(":5026", nil))
}
