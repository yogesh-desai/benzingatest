package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/yogesh-desai/benzingatest/utils"
)

// GetHealth is a simple health check endpoint
func GetHealth(w http.ResponseWriter, r *http.Request) {

	w.Write([]byte("OK")) // This call will automatically sets HTTP 200-OK as header response and body will be OK

}

// HandleLog handles the logs as per env variables
func HandleLog(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Request Method is not POST", http.StatusBadRequest)
		return
	}

	buf, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var payload Payload = Payload{}

	rdr := ioutil.NopCloser(bytes.NewBuffer(buf))

	err = json.NewDecoder(rdr).Decode(&payload)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var (
		statusCode int
		duration   time.Duration
	)

	// store in cache
	if utils.GetCount().Key < utils.Cfg.BatchSize {

		_ = utils.Cache.Add(fmt.Sprintf("%v", utils.GetCount().Key), payload)
		utils.IncrementCacheCount()
	} else {

		//post to endpoint and clear cache.
		log.Println("[INFO][HandleLog] Posting to external endpoint. BatchSize: ", utils.Cfg.BatchSize)
		statusCode, duration, err = postToExternalEndpoint()

		log.Println("[INFO][HandleLog] StatusCode: ", statusCode, "\tduration: ", duration, "\terr: ", err)
	}

	// Cache is cleared so make sure to store current payload to cache.
	if statusCode == 200 && utils.Cache.Len() == 0 {
		_ = utils.Cache.Add(fmt.Sprintf("%v", utils.GetCount().Key), payload)
		utils.IncrementCacheCount()
	}
	// start goroutine to post after each interval
	go doPostAtInterval(time.Duration(utils.Cfg.BatchInterval)*time.Second, postToExternalEndpoint)
}

// postToExternalEndpoint post data to external endpoint
func postToExternalEndpoint() (int, time.Duration, error) {

	var postDuration time.Duration
	var batchLoad []interface{} = make([]interface{}, utils.Cfg.BatchSize)

	// Make sure not to post while cache is empty. specially for interval routine
	if utils.Cache.Len() == 0 {
		return 0, postDuration, nil
	}
	for i := 0; i < utils.Cfg.BatchSize; i++ {

		subLoad, _ := utils.Cache.Get(fmt.Sprintf("%v", i))
		batchLoad[i] = subLoad //store it in array
	}

	load, err := json.Marshal(batchLoad)
	if err != nil {
		log.Println("[ERROR][postToExternalEndpoint] Couldn't marshal payload. ", err.Error())
		return 0, postDuration, err
	}

	noOpBody := ioutil.NopCloser(bytes.NewBuffer(load))
	// New request to post
	req, err := http.NewRequest(http.MethodPost, strings.TrimSpace(utils.Cfg.Endpoint), noOpBody)
	if err != nil {
		log.Println("[ERROR][postToExternalEndpoint] Couldn't make request object. ", err.Error())
		return 0, postDuration, err
	}
	req.Header.Add("Content-Type", "application/json")
	log.Println("[INFO][postToExternalEndpoint] The current config values are: batchSize: ", utils.Cfg.BatchSize, "\tBatchInterval: ", utils.Cfg.BatchInterval, "\tEndpoint: ", utils.Cfg.Endpoint)

	client := &http.Client{}

	start := time.Now()
	res, err := client.Do(req)
	if err != nil { // This is just bare minimum logic. We can use any standard libs in actual project.
		req.Body = noOpBody
		res, err = client.Do(req)
		if err != nil {
			req.Body = noOpBody
			res, err = client.Do(req)
			if err != nil {
				req.Body = noOpBody
				res, err = client.Do(req)
				if err != nil {
					log.Fatal("[FATAL][postToExternalEndpoint] Retried 3 times and all failed, Exiting application. ", err.Error())
				}
			}
		}
	}

	postDuration = time.Since(start)
	defer res.Body.Close()

	//Clear cache contents
	utils.Cache.Purge()
	utils.ResetCacheCount() // as cache is cleared, need to reset the counts.

	return res.StatusCode, postDuration, nil
}

func doPostAtInterval(d time.Duration, f func() (int, time.Duration, error)) {
	for range time.Tick(d) {
		f()
	}
}

type Payload struct {
	UserID    string  `json:"user_id,omitempty"`
	Total     float32 `json:"total,omitempty"`
	Title     string  `json:"title,omitempty"`
	Meta      Meta    `json:"meta,omitempty"`
	Completed bool    `json:"completed,omitempty"`
}

type Meta struct {
	Logins       []LoginInfo  `json:"logins,omitempty"`
	PhoneNumbers PhoneNumbers `json:"phone_numbers,omitempty"`
}

type LoginInfo struct {
	Time string `json:"time,omitempty"`
	Ip   string `json:"ip,omitempty"`
}
type PhoneNumbers struct {
	Home   string `json:"home,omitempty"`
	Mobile string `json:"mobile,omitempty"`
}
