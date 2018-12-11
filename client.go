package blackbox

import (
	"bytes"
	"context"
	"encoding/json"
	"encoding/xml"
	"io/ioutil"
	"net/http"
	"net/url"
)

type Error struct {
	HttpStatus  int
	RawResponse string
}

func (e Error) Error() string {
	return e.RawResponse
}

type client struct {
	url        string
	key        string
	signature  string
	httpClient *http.Client
}

func NewClient(url, key, signature string, httpClient *http.Client) Client {
	return &client{
		url:        url,
		key:        key,
		signature:  signature,
		httpClient: httpClient,
	}
}

type Message struct {
	Recipient string `xml:"recipient"`
	Message   string `xml:"message"`
	Sender    string `xml:"sender"`
	Keyword   string `xml:"keyword"`
}

type SendSMSRequest struct {
	XMLName  string    `xml:"request"`
	Messages []Message `xml:"sms"`
}

type SendSMSResponse struct {
	Response struct {
		Status struct {
			Code        string `json:"code"`
			Type        string `json:"type"`
			Description string `json:"description"`
			Meta        string `json:"meta"`
		} `json:"status"`
		Content struct {
			Description string `json:"description"`
			Messages    struct {
				Message string `json:"message"`
				Request struct {
					Sms []struct {
						Event             string `json:"event"`
						Recipient         string `json:"recipient"`
						Message           string `json:"message"`
						Sender            string `json:"sender"`
						Keyword           string `json:"keyword,omitempty"`
						Reference         string `json:"reference"`
						Status            string `json:"status"`
						StatusDescription string `json:"status_description"`
						Date              string `json:"date"`
						ScheduledDate     string `json:"scheduled_date,omitempty"`
					} `json:"sms"`
				} `json:"request"`
			} `json:"messages"`
		} `json:"content"`
	}
}

type Client interface {
	SendSMS(ctx context.Context, req *SendSMSRequest) (*SendSMSResponse, error)
}

func (c *client) SendSMS(ctx context.Context, req *SendSMSRequest) (*SendSMSResponse, error) {
	if req == nil {
		return nil, nil
	}

	reqXML, err := xml.Marshal(req)
	if err != nil {
		return nil, err
	}

	q := url.QueryEscape(string(reqXML))

	ret, err := c.roundTrip(ctx, http.MethodPost, "/send_sms", nil, []byte(q))
	if err != nil {
		return nil, err
	}

	resp := &SendSMSResponse{}
	err = json.Unmarshal(ret, resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (c *client) roundTrip(ctx context.Context, method, path string, query url.Values, body []byte) (ret []byte, err error) {
	u, err := url.Parse(c.url + path)
	if err != nil {
		return
	}
	u.RawQuery = query.Encode()

	req, err := http.NewRequest(method, u.String(), bytes.NewBuffer(body))
	if err != nil {
		return
	}

	if ctx != nil {
		req = req.WithContext(ctx)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	ret, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	if resp.StatusCode != http.StatusOK {
		err = &Error{
			HttpStatus:  resp.StatusCode,
			RawResponse: string(ret),
		}
	}

	return ret, err
}
