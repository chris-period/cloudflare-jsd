# Overview
Retreives **"cf_clearance"** cookie from cloudflare "jsd" challenge

### Constants
- **BrowserConfiguration**: The payload of the challenge (old payload, `06/26/2023`)

### LZString
- LZString, Lempel-Ziv (LZ) compression algorithm.
- Originally retrieved from: [daku10/go-lz-string](https://github.com/daku10/go-lz-string)

### Classes
`cloudflare.Client`
A class that creates a "*Client" struct for the solver.

- **Constructor**: Initializes the client.
- **CreatePayload(userAgent string)**: Creates a payload struct that modifies payload to consist of the username provided, and update the time to current time.
- **Solve()**: Sends the modified payload to the respective endpoint.

`cloudflare.payload`
A class that creates a "*payload" struct for payload construction.

- **Constructor**: Initializes the client with the userAgent
- **constructPayload()**: Copies the "BrowserConfiguration" constant, to it's own "payload.Data" map.
- **buildPayload()**: Modifies copied payload,"payload.data", with new useragent information and current time.

# Parser
  Finds LZStringKey to compress with, and the secretKey within the javascript file with regex.
  ```golang
  var lzKeyRegex = regexp.MustCompile(`[^\s,]*\$[^\s,]*\+?[^\s,]*`)
  var sKeyRegex = regexp.MustCompile(`\d+\.\d+:\d+:[^\s,]+`)
  ```
