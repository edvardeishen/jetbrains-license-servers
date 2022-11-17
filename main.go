package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

var availableHost []string
var testQueue = make(chan int, 10)

func getAppIDAndSecret() (string, string) {

	err := godotenv.Load()
	if err != nil {
		panic("Error when load .env file: " + err.Error())
	}

	//这个ID和Secret可能被重置了, 请到 https://search.censys.io/account/api 获取
	apiID := os.Getenv("API_ID")
	apiSecret := os.Getenv("SECRET")

	if apiID == "" || apiSecret == "" {
		panic("Please set API_ID and SECRET in .env file")
	}
	return apiID, apiSecret
}

func queryFromCensys() (string, error) {
	url := "https://search.censys.io/api/v2/hosts/search?q=service.service_name:%20HTTP%20AND%20services.http.response.headers.location:%20account.jetbrains.com/fls-auth&per_page=50&virtual_hosts=EXCLUDE"
	method := "GET"

	//避免x509: certificate signed by unknown authority
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)
		return "", err
	}
	req.Header.Add("Accept", "application/json")

	id, secret := getAppIDAndSecret()
	req.SetBasicAuth(id, secret)

	res, err := client.Do(req)
	if err != nil {

		return "", err
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(res.Body)

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	return string(body), nil
}

func TestHost(ip string, port int) {
	testQueue <- 1

	url := "http://" + ip + ":" + strconv.Itoa(port) + "/"
	//避免x509: certificate signed by unknown authority
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{
		Timeout:   time.Second * 5,
		Transport: tr,
	}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println(err)
		<-testQueue
		return
	}

	response, err := client.Do(req)
	if err != nil {
		log.Println(url, "  http get error: ", err)
		<-testQueue
		return
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Println(url, "  http close error: ", err)
		}
	}(response.Body)

	if response.StatusCode == 200 && response.Request.URL.Host == "account.jetbrains.com" && response.Request.URL.Path == "/login" {
		log.Println(url, "  is ok")
		availableHost = append(availableHost, url)
	}
	<-testQueue
}

func parseCensysResult(body string) {
	fmt.Println(body)

	fmt.Print("\n\n------------------\n\n")
	fmt.Println("耐心等待...")

	var result HostSearchResult
	err := json.Unmarshal([]byte(body), &result)
	if err != nil {
		fmt.Println("error when unmarshal json: ", err)
	}

	if result.Code != 200 && result.Status != "OK" {
		fmt.Println("error code from censys: ", result.Status)
		return
	}

	for _, hit := range result.Result.Hits {
		for _, service := range hit.Services {
			if service.ServiceName == "HTTP" {
				fmt.Printf("testing: %s %s:%d \n", hit.AutonomousSystem.Name, hit.IP, service.Port)
				go TestHost(hit.IP, service.Port)
			}
		}
	}

}

func main() {
	jsonString, err := queryFromCensys()
	if err != nil {
		fmt.Println(err)
		return
	}
	parseCensysResult(jsonString)

	for len(testQueue) > 0 {
		time.Sleep(time.Second)
	}

	fmt.Print("\n\n------------------\n\n")
	fmt.Println("可用的服务器:")
	for _, host := range availableHost {
		fmt.Println(host)
	}

	fmt.Println("\n\ndone!!!")
}
