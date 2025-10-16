package main

import (
	"flag"
	"fmt"
	"github.com/go-resty/resty/v2"
	"math/rand"
	"sync"
	"time"
)

func main() {
	numReq := flag.Int("numReqs", 100, "number of requests to run")
	maxCon := flag.Int("maxCons", 3, "number of concurrent go routines to run")

	flag.Parse()

	totalReq := *numReq * *maxCon

	fmt.Printf("total %d requests will be made\n", *numReq**maxCon)

	client := resty.New()

	wg := sync.WaitGroup{}
	wg.Add(totalReq)

	startTime := time.Now()

	for i := 0; i < *maxCon; i++ {
		go func(n int) {
			for j := 0; j < n; j++ {
				_, err := client.R().
					SetHeader("Content-Type", "application/json").
					SetBody(randomReqBody()).
					Post("http://localhost:9000/gamers/points/credit")

				if err != nil {
					//panic(err)
				}
				wg.Done()
			}
		}(*numReq)
	}

	wg.Wait()

	endtime := time.Now()

	fmt.Println("test time in milliseconds", endtime.Sub(startTime).Milliseconds())
	fmt.Println("test time in seconds", endtime.Sub(startTime).Seconds())

}

type reqBody struct {
	Game    string `json:"game"`
	Points  int    `json:"points"`
	GamerID string `json:"gamer_id"`
}

func randomReqBody() reqBody {
	bodies := []reqBody{
		{
			Game:    "call of duty",
			Points:  1987,
			GamerID: "gunExpert",
		},
		{
			Game:    "clash royale",
			Points:  1200,
			GamerID: "goku",
		},
		{
			Game:    "candy crush",
			Points:  200,
			GamerID: "sugar",
		},
		{
			Game:    "hay day",
			Points:  123,
			GamerID: "farm_animal",
		},
		{
			Game:    "",
			Points:  0,
			GamerID: "",
		},
	}

	return bodies[rand.Int()%len(bodies)]
}
