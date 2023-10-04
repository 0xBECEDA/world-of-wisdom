package tcp_message

type Type int

const (
	ChallengeReq Type = iota
	ChallengeResp
	QuoteReq
	QuoteResp
)
