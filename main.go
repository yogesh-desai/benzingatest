package main

import (
	"log"
	"net/http"

	"github.com/yogesh-desai/benzingatest/api"
	"github.com/yogesh-desai/benzingatest/utils"
)

func init() {

	api.Init()
	utils.Init()
}

func main() {

	log.Println("[main] Started listening at :9000 ")
	http.ListenAndServe(":9000", nil)

}
