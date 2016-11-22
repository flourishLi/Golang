package httputil

type ResultBase struct {
	Code         int    `json:"code"`
	ErrorMessage string `json:"errorMessage"`
	ResponseTime int64  `json:"responseTime"`
}

type RequestBase struct {
	Command     string //`json:"command"`
	UserId      int    //`json:"userId"`
	Tocken      string //`json:"tocken"`
	RequestTime int64  //`json:"requestTime"`
}
