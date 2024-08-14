package handler

type Request struct {
	Value string `json:"value"`
	TTL   int64  `json:"ttl"`
}

type Response struct {
	Value string `json:"value"`
}
