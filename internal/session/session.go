package session

type Session struct {
	Day       string
	Date      string
	Available int
	Limit     int
	Time      string
}

type Subscription struct {
	Name          string
	Remaining     string
	Date          string
	PostRequestId string
}
