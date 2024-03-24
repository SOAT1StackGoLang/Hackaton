// Description: This is a simple client that will make requests to the API Gateway endpoint.
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
)

type Config struct {
	Endpoint        string
	AWSRegion       string
	ClientID        string
	UserPoolID      string
	Username        string
	Password        string
	NumThreads      int
	NumRequests     int
	NumGetRequests  int
	NumPostRequests int
	Token           string
	TokenExpires    time.Time
}

type Client struct {
	Config
}

func NewClient(config Config) *Client {
	return &Client{config}
}

func (c *Client) Run() {
	var wg sync.WaitGroup
	for i := 0; i < c.NumThreads; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			c.run()
		}()
	}
	wg.Wait()
}

func (c *Client) run() {
	for i := 0; i < c.NumRequests; i++ {
		if i < c.NumGetRequests {
			duration := c.get()
			// print the duration for each 100 requests
			if i%100 == 0 {
				log.Printf("GET Request %d: %s\n", i, duration)
			}
		} else {
			duration := c.post()
			// print the duration for each 100 requests
			if i%100 == 0 {
				log.Printf("POST Request %d: %s\n", i, duration)
			}
		}
	}
}

func (c *Client) get() time.Duration {
	// make a get and print the response time
	token := c.getToken()
	startTime := time.Now() // Record the start time

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/reports/daily?reference=2024-03-24", c.Endpoint), nil)
	if err != nil {
		log.Println(err)
		return time.Duration(0)
	}
	req.Header.Set("accept", "application/json")
	req.Header.Set("user_id", "testing")
	req.Header.Set("Authorization", token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println(err)
		return time.Duration(0)
	}
	defer resp.Body.Close()

	// body, err := ioutil.ReadAll(resp.Body)
	// if err != nil {
	// 	log.Println(err)
	// 	return
	// }

	// Calculate and print the response time
	responseTime := time.Since(startTime)
	// log.Printf("[GET]Response Time: %s\n", responseTime)

	return responseTime
}

func (c *Client) post() time.Duration {
	token := c.getToken()
	data := map[string]string{"entry_at": "2024-03-24T00:52:24Z"}
	payload, err := json.Marshal(data)
	if err != nil {
		log.Println(err)
		return time.Duration(0)
	}
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/clock-in", c.Endpoint), bytes.NewBuffer(payload))
	if err != nil {
		log.Println(err)
		return time.Duration(0)
	}
	req.Header.Set("accept", "application/json")
	req.Header.Set("user_id", "testing")
	req.Header.Set("Authorization", token)
	req.Header.Set("Content-Type", "application/json")

	startTime := time.Now() // Record the start time

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println(err)
		return time.Duration(0)
	}
	defer resp.Body.Close()

	// body, err := ioutil.ReadAll(resp.Body)
	// if err != nil {
	// 	log.Println(err)
	// 	return
	// }

	// Calculate and print the response time
	responseTime := time.Since(startTime)
	log.Printf("[POST] Response Time: %s\n", responseTime)

	return responseTime
}

func (c *Client) getToken() string {
	// If the token is not empty and it has not expired, return the cached token
	if c.Token != "" && time.Now().Before(c.TokenExpires) {
		return c.Token
	}

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(c.AWSRegion),
	})
	if err != nil {
		log.Println(err)
		return ""
	}

	svc := cognitoidentityprovider.New(sess)

	authInput := &cognitoidentityprovider.InitiateAuthInput{
		AuthFlow: aws.String("USER_PASSWORD_AUTH"),
		ClientId: aws.String(c.ClientID),
		AuthParameters: map[string]*string{
			"USERNAME": aws.String(c.Username),
			"PASSWORD": aws.String(c.Password),
		},
	}

	authOutput, err := svc.InitiateAuth(authInput)
	if err != nil {
		log.Println(err)
		return ""
	}

	// Cache the token and its expiration time
	c.Token = *authOutput.AuthenticationResult.AccessToken
	c.TokenExpires = time.Now().Add(time.Duration(*authOutput.AuthenticationResult.ExpiresIn) * time.Second)

	return c.Token
}

func (c *Client) RunWithRate(rate int) {
	var wg sync.WaitGroup
	for i := 0; i < c.NumThreads; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			c.runWithRate(rate)
		}()
	}
	wg.Wait()
}

func (c *Client) runWithRate(rate int) {
	ticker := time.NewTicker(time.Second / time.Duration(rate))
	defer ticker.Stop()
	for range ticker.C {
		c.run()
	}
}

func (c *Client) RunWithRateAndDuration(rate int, duration time.Duration) {
	var wg sync.WaitGroup
	for i := 0; i < c.NumThreads; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			c.runWithRateAndDuration(rate, duration)
		}()
	}
	wg.Wait()
}

func (c *Client) runWithRateAndDuration(rate int, duration time.Duration) {
	ticker := time.NewTicker(time.Second / time.Duration(rate))
	defer ticker.Stop()
	timer := time.NewTimer(duration)
	defer timer.Stop()
	for {
		select {
		case <-ticker.C:
			c.run()
		case <-timer.C:
			return
		}
	}
}

func (c *Client) RunWithRateAndDurationAndWarmUp(rate int, duration time.Duration, warmUp time.Duration) {
	var wg sync.WaitGroup
	for i := 0; i < c.NumThreads; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			c.runWithRateAndDurationAndWarmUp(rate, duration, warmUp)
		}()
	}
	wg.Wait()
}

func (c *Client) runWithRateAndDurationAndWarmUp(rate int, duration time.Duration, warmUp time.Duration) {
	ticker := time.NewTicker(time.Second / time.Duration(rate))
	defer ticker.Stop()
	timer := time.NewTimer(duration)
	defer timer.Stop()
	warmUpTimer := time.NewTimer(warmUp)
	defer warmUpTimer.Stop()
	for {
		select {
		case <-ticker.C:
			c.run()
		case <-timer.C:
			return
		case <-warmUpTimer.C:
			ticker = time.NewTicker(time.Second / time.Duration(rate))
		}
	}
}

func (c *Client) RunWithRateAndDurationAndWarmUpAndCooldown(rate int, duration time.Duration, warmUp time.Duration, cooldown time.Duration) {
	var wg sync.WaitGroup
	for i := 0; i < c.NumThreads; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			c.runWithRateAndDurationAndWarmUpAndCooldown(rate, duration, warmUp, cooldown)
		}()
	}
	wg.Wait()
}

func (c *Client) runWithRateAndDurationAndWarmUpAndCooldown(rate int, duration time.Duration, warmUp time.Duration, cooldown time.Duration) {
	ticker := time.NewTicker(time.Second / time.Duration(rate))
	defer ticker.Stop()
	timer := time.NewTimer(duration)
	defer timer.Stop()
	warmUpTimer := time.NewTimer(warmUp)
	defer warmUpTimer.Stop()
	cooldownTimer := time.NewTimer(cooldown)
	defer cooldownTimer.Stop()
	for {
		select {
		case <-ticker.C:
			c.run()
		case <-timer.C:
			return
		case <-warmUpTimer.C:
			ticker = time.NewTicker(time.Second / time.Duration(rate))
		case <-cooldownTimer.C:
			ticker.Stop()
		}
	}
}

// lets create a test that will run for 10 minutes and will generate 1000 requests per second.
// we will have a warm up time of 10 seconds and a cooldown of 15 seconds.
// must only print the rate of requests per second.
// first lets create the test base

func main() {
	config := Config{
		Endpoint:        "https://qin95h1ccj.execute-api.us-east-1.amazonaws.com",
		AWSRegion:       "us-east-1",
		ClientID:        "3f63krirom0trasn2duulhfktu",
		UserPoolID:      "us-east-1_2Gtre7AbT",
		Username:        "11122233300",
		Password:        "F@ap1234",
		NumThreads:      50,
		NumRequests:     10,
		NumGetRequests:  7,
		NumPostRequests: 3,
	}
	client := NewClient(config)
	//client.RunWithRateAndDurationAndWarmUpAndCooldown(1000, 10*time.Minute, 10*time.Second, 15*time.Second)

	// lets do a controlled test to reach 1000 transactions per second
	// lets do a test that will run for 10 minutes and will generate 1000 requests per second.
	start := time.Now()
	client.RunWithRateAndDuration(1000, 10*time.Minute)
	elapsed := time.Since(start)
	log.Printf("Elapsed time: %s\n", elapsed)

}
