package rtgo

import "fmt"
import "io/ioutil"
import "net/http"
import "strconv"
import "strings"
import "net/url"

type RT struct {
	// URL of the request tracker REST API
	URL string

	// Username to authenticate with
	Username string

	// Password to authenticate with
	Password string

	Client *http.Client
}

// Constructs a new RT instance for API requests.
// The URL is the base URL to the request tracker installation, without a trailing slash.
func NewRT(url string, username string, password string) *RT {
	return &RT{
		URL: url + "/REST/1.0/",
		Username: username,
		Password: password,
		Client: &http.Client{},
	}
}

func (rt *RT) request(path string, request map[string]string, response interface{}, long bool) (string, error) {
	postData := make(url.Values)
	postData.Set("user", rt.Username)
	postData.Set("pass", rt.Password)

	if request != nil {
		content := ""
		for key, value := range request {
			content += fmt.Sprintf("%s: %s\n", key, strings.Replace(value, "\n", "\n ", -1))
		}
		postData.Set("content", content)
	}

	resp, err := rt.Client.PostForm(rt.URL + path, postData)
	if err != nil {
		return "", fmt.Errorf("error performing HTTP request: %v", err)
	}

	if resp.StatusCode == 401 {
		return "", fmt.Errorf("server refused the provided user credentials")
	} else if resp.StatusCode != 200 {
		return "", fmt.Errorf("server returned status code %d", resp.StatusCode)
	}

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading HTTP response: %v", err)
	}

	responseString := string(bytes)
	responseParts := strings.SplitN(responseString, "\n", 3)
	if len(responseParts) != 3 {
		return "", fmt.Errorf("response body contains less than three lines")
	} else if !strings.Contains(responseParts[0], "RT") {
		return "", fmt.Errorf("first line of response body does not contain \"RT\"")
	}

	if response != nil {
		if long {
			// ...
		} else {
			err := UnmarshalShort(responseParts[2], response)
			if err != nil {
				return "", fmt.Errorf("error decoding with short form: %v", err)
			}
		}
	}

	return responseString, nil
}

func (rt *RT) CreateTicket(queue string, requestor string, subject string, text string) (int, error) {
	request := map[string]string{
		"Queue": queue,
		"Requestor": requestor,
		"Subject": subject,
		"Text": text,
		"id": "ticket/new",
	}
	response, err := rt.request("ticket/new", request, nil, false)
	if err != nil {
		return 0, err
	}

	parts1 := strings.Split(response, "# Ticket ")
	if len(parts1) < 2 {
		return 0, fmt.Errorf("response does not include ticket number")
	}
	parts2 := strings.Split(parts1[1], " created")
	if len(parts2) < 2 {
		return 0, fmt.Errorf("response does not include ticket number")
	}

	ticketNumberStr := parts2[0]
	ticketNumber, err := strconv.Atoi(ticketNumberStr)
	if err != nil {
		return 0, fmt.Errorf("failed to convert ticket number %s to int", ticketNumberStr)
	}

	return ticketNumber, nil
}
