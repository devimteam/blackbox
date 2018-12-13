package blackbox

import (
	"context"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func defaultHandler(response string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(response))
	}
}

func Test_SendSMS(t *testing.T) {
	type test struct {
		name    string
		handler http.HandlerFunc
		req     *SendSMSRequest
	}

	cases := []test{
		{
			name: "t1",
			handler: defaultHandler(`
{
  "response": {
    "status": {
      "code": "0000",
      "type": "SUCCESS",
      "description": "REQUEST_SUCCESSFUL",
      "meta": "SMS_QUEUE"
    },
    "content": {
      "description": "2 messages queued. 4 SMS deducted, 2504256 SMS available.",
      "messages": {
        "message": "2504256",
        "request": {
          "sms": [
            {
              "event": "OUTBOX",
              "recipient": "+25472212356",
              "message": "This is an unscheduled message",
              "sender": "",
              "keyword": "",
              "reference": "",
              "status": "FAILED",
              "status_description": "INVALID_SERVICE",
              "date": "2015-01-28 18:45:00"
            },
            {
              "event": "OUTBOX",
              "recipient": "+25472212356",
              "message": "This is a shortcode message via music keyword",
              "sender": "20717",
              "keyword": "MUSIC",
              "reference": "bb0b30ea2208c99b1164c965ade9dfb8cbb23fdaa1",
              "status": "SENT",
              "status_description": "OK",
              "date": "2015-01-28 18:45:00"
            },
            {
              "event": "OUTBOX",
              "recipient": "+25472212356",
              "message": "This is a shortcode message without keyword",
              "sender": "20333",
              "reference": "bb0b30ea2208c99b1164c965ade9dfb8cbb23fdaa2",
              "status": "SENT",
              "status_description": "OK",
              "date": "2015-01-28 18:45:00"
            },
            {
              "event": "OUTBOX",
              "recipient": "+25473312356",
              "message": "This is a scheduled message",
              "sender": "SENDER_123",
              "scheduled_date": "2015-01-28 18:45:00",
              "reference": "bb0b30ea2208c99b1164c965ade9dfb8cbb23fdaa3",
              "status": "SCHEDULED",
              "status_description": "OK",
              "date": "2015-01-28 18:45:00"
            }
          ]
        }
      }
    }
  }
}`),
			req: &SendSMSRequest{
				Messages: []Message{
					{
						Recipient: "w1",
						Message:   "w2",
						Sender:    "h1",
						Keyword:   "h2",
					},
				},
			},
		},
	}

	for _, tc := range cases {
		ts := httptest.NewServer(tc.handler)
		c := NewClient(ts.URL, "", "", http.DefaultClient)
		t.Run(tc.name, func(t *testing.T) {
			resp, err := c.SendSMS(context.Background(), tc.req)
			assert.NoError(t, err)
			assert.NotNil(t, resp)
			assert.Equal(t, resp.Response.Content.Description, "2 messages queued. 4 SMS deducted, 2504256 SMS available.")
		})
		ts.Close()
	}
}
