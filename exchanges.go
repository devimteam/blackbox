package blackbox

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
