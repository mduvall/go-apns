package apns

type Payload struct {
	Type    string
	Message string
}

func (p *Payload) MarshalJSON() ([]byte, error) {
	return []byte(`{"aps":{"` + p.Type + `": "` + p.Message + `"}}`), nil
}
