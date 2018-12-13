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

const (
	StatusFailed    = "FAILED"
	StatusScheduled = "SCHEDULED"
	StatusSent      = "SENT"
	StatusRejected  = "REJECTED"
	StatusSuccess   = "SUCCESS"

	StatusDescriptionOk = "OK"
)

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

type Client interface {
	SendSMS(ctx context.Context, req *SendSMSRequest) (*SendSMSResponse, error)
}

func (c *client) SendSMS(ctx context.Context, msgReq *SendSMSRequest) (*SendSMSResponse, error) {
	if msgReq == nil {
		return nil, nil
	}

	msgXML, err := xml.Marshal(msgReq)
	if err != nil {
		return nil, err
	}

	vals := url.Values{}
	vals.Add("messages", string(msgXML))

	ret, err := c.postUrlEncoded(ctx, "/send_sms", vals)
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

func (c *client) postUrlEncoded(ctx context.Context, path string, vals url.Values) (ret []byte, err error) {
	u, err := url.Parse(c.url + path)
	if err != nil {
		return
	}

	vals.Add("api_key", c.key)
	vals.Add("api_signature", c.signature)
	vals.Add("api_format", "JSON")

	req, err := http.NewRequest(http.MethodPost, u.String(), bytes.NewBuffer([]byte(vals.Encode())))
	if err != nil {
		return
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
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
