package cloudflare

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"maps"
	"net/http"
	"strings"
	"time"
)

type payload struct {
	UserAgent string
	Data      map[string]interface{}
}

func (p *payload) constructPayload() {
	maps.Copy(p.Data, BrowserConfiguration)
}

func (p *payload) buildPayload() {
	p.Data["APPVERSION"] = strings.TrimPrefix(p.UserAgent, "Mozilla/")
	p.Data["USERAGENT"] = p.UserAgent
	p.Data["06/26/2023 06:47:34"] = time.Now().Format("01/02/2006 15:04:05")
}

type Client struct {
	HTTP    *http.Client
	Payload string
}

func (c *Client) getScript() (string, []byte) {
	req, err := http.NewRequest("GET", "https://crypto.com/cdn-cgi/challenge-platform/scripts/jsd/main.js", nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("accept", "*/*")
	req.Header.Set("accept-language", "en-US,en;q=0.9")
	req.Header.Set("cache-control", "no-cache")
	req.Header.Set("dnt", "1")
	req.Header.Set("pragma", "no-cache")
	req.Header.Set("sec-ch-ua", `"Not.A/Brand";v="8", "Chromium";v="114", "Google Chrome";v="114"`)
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("sec-ch-ua-platform", `"Windows"`)
	req.Header.Set("sec-fetch-dest", "script")
	req.Header.Set("sec-fetch-mode", "no-cors")
	req.Header.Set("sec-fetch-site", "same-origin")
	req.Header.Set("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36")
	resp, err := c.HTTP.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	bodyText, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	return resp.Header.Get("Cf-Ray"), bodyText
}

func (c *Client) sendPayload(rayID string, windowProperties string, secretKey string) {
	data := strings.NewReader(fmt.Sprintf(`{"wp":"%s","s":"%s"}`, windowProperties, secretKey))
	req, err := http.NewRequest("POST", "https://crypto.com/cdn-cgi/challenge-platform/h/g/jsd/r/"+rayID, data)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("accept", "*/*")
	req.Header.Set("accept-language", "en-US,en;q=0.9")
	req.Header.Set("cache-control", "no-cache")
	req.Header.Set("content-type", "application/json")
	req.Header.Set("dnt", "1")
	req.Header.Set("origin", "https://crypto.com")
	req.Header.Set("pragma", "no-cache")
	req.Header.Set("sec-ch-ua", `"Not.A/Brand";v="8", "Chromium";v="114", "Google Chrome";v="114"`)
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("sec-ch-ua-platform", `"Windows"`)
	req.Header.Set("sec-fetch-dest", "empty")
	req.Header.Set("sec-fetch-mode", "cors")
	req.Header.Set("sec-fetch-site", "same-origin")
	req.Header.Set("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36")
	resp, err := c.HTTP.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	fmt.Println(`[cf_clearance] Status: `, resp.Status)
	fmt.Println(resp.Cookies())
}

func (c *Client) CreatePayload(userAgent string) {
	payload := &payload{
		UserAgent: userAgent,
		Data:      make(map[string]interface{}),
	}
	payload.constructPayload()
	payload.buildPayload()

	pBytes, err := json.Marshal(payload.Data)
	if err != nil {
		log.Fatal(err)
	}
	c.Payload = string(pBytes)
}

func (c Client) Solve() {
	rayID, script := c.getScript()
	rayID = strings.Split(rayID, "-")[0]
	lzStringKey, secretKey := parse(string(script))

	windowProperties, err := new(LZString).CompressToEncodedURIComponent(c.Payload, lzStringKey)
	if err != nil {
		panic(err)
	}
	c.sendPayload(rayID, windowProperties, secretKey)
}

func CreateClient(client *http.Client) *Client {
	return &Client{
		HTTP:    client,
		Payload: "",
	}
}
