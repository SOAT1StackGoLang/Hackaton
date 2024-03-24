// Description: This is a simple client that will make requests to the API Gateway endpoint.
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/joho/godotenv"
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
	mu              sync.Mutex
	requestCount    int
}

type Client struct {
	Config
}

func NewClient(config *Config) *Client {
	return &Client{
		Config: Config{
			Endpoint:        config.Endpoint,
			AWSRegion:       config.AWSRegion,
			ClientID:        config.ClientID,
			UserPoolID:      config.UserPoolID,
			Username:        config.Username,
			Password:        config.Password,
			NumThreads:      config.NumThreads,
			NumRequests:     config.NumRequests,
			NumGetRequests:  config.NumGetRequests,
			NumPostRequests: config.NumPostRequests,
			Token:           config.Token,
			TokenExpires:    config.TokenExpires,
		},
	}
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
			c.get()
		} else {
			c.post()
		}
	}
}

func (c *Client) get() {
	// make a get and print the response time
	token := c.getToken()

	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/reports/daily?reference=2024-03-24", c.Endpoint), nil)
	if err != nil {
		log.Println(err)
	}
	req.Header.Set("accept", "application/json")
	req.Header.Set("user_id", "testing")
	req.Header.Set("Authorization", token)

	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
		log.Fatal("Error making request")
	}
	defer resp.Body.Close()

	// body, err := ioutil.ReadAll(resp.Body)
	// if err != nil {
	// 	log.Println(err)
	// 	return
	// }

	// log.Printf("[GET]Response Time: %s\n", responseTime)

}

func (c *Client) post() {
	token := c.getToken()
	// add timeout to http client
	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	data := map[string]string{"entry_at": "2024-03-24T00:52:24Z"}
	payload, err := json.Marshal(data)
	if err != nil {
		log.Println(err)
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/entries", c.Endpoint), bytes.NewBuffer(payload))
	if err != nil {
		log.Println(err)
	}
	req.Header.Set("accept", "application/json")
	req.Header.Set("user_id", "testing")
	req.Header.Set("Authorization", token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
		log.Fatal("Error making request")
	}
	defer resp.Body.Close()

	// body, err := ioutil.ReadAll(resp.Body)
	// if err != nil {
	// 	log.Println(err)
	// 	return
	// }

	// Calculate and print the response time
	// log.Printf("[POST] Response Time: %s\n", responseTime)

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
		startTime := time.Now()
		c.run()
		duration := time.Since(startTime)
		log.Printf("Execution time: %s", duration)
	}
}

func (c *Client) runWithRateAndDuration(threadNum, rate int, duration time.Duration) {
	ticker := time.NewTicker(time.Second / time.Duration(rate))
	defer ticker.Stop()
	timer := time.NewTimer(duration)
	defer timer.Stop()
	progressTicker := time.NewTicker(10 * time.Second)
	defer progressTicker.Stop()
	startTime := time.Now()
	for {
		select {
		case <-ticker.C:
			c.run()
			c.mu.Lock()
			c.requestCount++
			c.mu.Unlock()
		case <-timer.C:
			elapsedTime := time.Since(startTime)
			log.Printf("Thread %d: Total requests: %d", threadNum, c.requestCount)
			log.Printf("Thread %d: Elapsed time: %s", threadNum, elapsedTime)
			return
		case <-progressTicker.C:
			c.mu.Lock()
			//elapsedTime := time.Since(startTime)
			//rate := float64(c.requestCount) / elapsedTime.Seconds()
			c.mu.Unlock()
			//log.Printf("Thread %d: Current rate: %f requests per second", threadNum, rate)
			//log.Printf("Thread %d: Current duration: %s", threadNum, elapsedTime)
		}
	}
}

func (c *Client) RunWithRateAndDuration(rate int, duration time.Duration) {
	var wg sync.WaitGroup
	c.requestCount = 0 // Reset the request count before starting
	summaryTicker := time.NewTicker(10 * time.Second)
	defer summaryTicker.Stop()
	startTime := time.Now()

	go func() {
		for range summaryTicker.C {
			c.mu.Lock()
			elapsedTime := time.Since(startTime)
			rate := float64(c.requestCount) / elapsedTime.Seconds()
			c.mu.Unlock()
			log.Printf("Overall summary: Current rate: %f requests per second", rate)
			log.Printf("Overall summary: Current duration: %s", elapsedTime)
		}
	}()

	for i := 0; i < c.NumThreads; i++ {
		wg.Add(1)
		go func(threadNum int) {
			defer wg.Done()
			c.runWithRateAndDuration(threadNum, rate, duration)
		}(i)
	}
	wg.Wait()
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

func loadConfig() Config {
	config := Config{
		Endpoint:        "https://ttqe3moh10.execute-api.us-east-1.amazonaws.com",
		AWSRegion:       "us-east-1",
		ClientID:        "4hojr2kivbpeuis0so7f4liqih",
		UserPoolID:      "us-east-1_eagKRqRru",
		Username:        "11122233300",
		Password:        "F@ap1234",
		NumThreads:      10,
		NumRequests:     100,
		NumGetRequests:  70,
		NumPostRequests: 30,
	}

	err := godotenv.Load(".testenv")
	if err != nil {
		log.Println("Error loading .env file:", err)
		return config
	}

	// Override the configuration values from the .env file
	if endpoint := os.Getenv("ENDPOINT"); endpoint != "" {
		config.Endpoint = endpoint
	}
	if awsRegion := os.Getenv("AWS_REGION"); awsRegion != "" {
		config.AWSRegion = awsRegion
	}
	if clientID := os.Getenv("CLIENT_ID"); clientID != "" {
		config.ClientID = clientID
	}
	if userPoolID := os.Getenv("USER_POOL_ID"); userPoolID != "" {
		config.UserPoolID = userPoolID
	}
	if username := os.Getenv("USERNAME"); username != "" {
		config.Username = username
	}
	if password := os.Getenv("PASSWORD"); password != "" {
		config.Password = password
	}
	if numThreads := os.Getenv("NUM_THREADS"); numThreads != "" {
		// Parse the value as an integer
		if threads, err := strconv.Atoi(numThreads); err == nil {
			config.NumThreads = threads
		} else {
			log.Println("Invalid value for NUM_THREADS:", err)
		}
	}
	if numRequests := os.Getenv("NUM_REQUESTS"); numRequests != "" {
		// Parse the value as an integer
		if requests, err := strconv.Atoi(numRequests); err == nil {
			config.NumRequests = requests
		} else {
			log.Println("Invalid value for NUM_REQUESTS:", err)
		}
	}
	if numGetRequests := os.Getenv("NUM_GET_REQUESTS"); numGetRequests != "" {
		// Parse the value as an integer
		if getRequests, err := strconv.Atoi(numGetRequests); err == nil {
			config.NumGetRequests = getRequests
		} else {
			log.Println("Invalid value for NUM_GET_REQUESTS:", err)
		}
	}
	if numPostRequests := os.Getenv("NUM_POST_REQUESTS"); numPostRequests != "" {
		// Parse the value as an integer
		if postRequests, err := strconv.Atoi(numPostRequests); err == nil {
			config.NumPostRequests = postRequests
		} else {
			log.Println("Invalid value for NUM_POST_REQUESTS:", err)
		}
	}

	// log each configuration value
	log.Println("Endpoint:", config.Endpoint)
	log.Println("AWS Region:", config.AWSRegion)
	log.Println("Client ID:", config.ClientID)
	log.Println("User Pool ID:", config.UserPoolID)
	log.Println("Username:", config.Username)
	log.Println("Password:", config.Password)
	log.Println("Number of Threads:", config.NumThreads)
	log.Println("Number of Requests:", config.NumRequests)
	log.Println("Number of GET Requests:", config.NumGetRequests)
	log.Println("Number of POST Requests:", config.NumPostRequests)

	return config
}

func main() {
	config := loadConfig()

	client := NewClient(&config)

	start := time.Now()
	log.Println("Starting test, please wait...")
	client.RunWithRateAndDuration(1000, 10*time.Minute)
	elapsed := time.Since(start)
	log.Printf("Elapsed time: %s\n", elapsed)

}
