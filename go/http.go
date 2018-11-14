package util

import (
	"io/ioutil"
	"net/url"
	"net/http"
	"strconv"
	"strings"
)

type Any interface{}

type HttpResp struct {
	Url string `json:"url"`
	Method string `json:"method"`
	Params Any `json:"params, omitempty"`
	HttpCode int `json:"httpCode"`
	HttpError string `json:"httpError, omitempty"`
	Raw string `json:"raw, omitempty"`
	ReqHeaders *http.Header `json:"reqHeaders, omitempty"`
	ReqCookies *map[string]string `json:"reqHookies, omitempty"`
	RespCookies []*http.Cookie `json:"respCookies, omitempty"`
}

func HttpGet (strUrl string, params map[string]Any, headers map[string]string, cookies map[string]string) *HttpResp {
	client := &http.Client{}
	ret := &HttpResp{}
    httpParams := url.Values{}
	for k, v := range params {
		switch realValue := v.(type) {
			case string:
				httpParams.Set(k, realValue)
			case int:
				httpParams.Set(k, strconv.Itoa(realValue))
			case []byte:
				httpParams.Set(k, string(realValue))
			case byte:
				httpParams.Set(k, string(realValue))
		}
	}
	//request, err := http.NewRequest("Post", strUrl, strings.NewReader(httpParams.Encode()))
	if len(httpParams) > 0 {
		strUrl += "?" + httpParams.Encode()
	}
	request, err := http.NewRequest("GET", strUrl, nil)
	if err != nil {
		ret.HttpError = err.Error()
		return ret
	}
	if nil != headers {
		for k, v := range headers {
			request.Header.Add(k, v)	
		}	
	}
	if 0 >= len(request.Header.Get("User-Agent")) {
		request.Header.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/61.0.3163.100 Safari/537.36")	
	}

	strCookie := ""
	if nil != cookies {
		for k, v := range cookies {
			strCookie += k + "=" + v + ";"
		}
	}
	request.Header.Add("Cookie", strCookie)
	ret.Params = &params
	ret.ReqHeaders = &request.Header
	ret.ReqCookies = &cookies
	response, err := client.Do(request)
	if nil != err {
		ret.HttpError = err.Error()
		return ret
	} else {
		ret.HttpCode = response.StatusCode
		ret.Url = strUrl
		ret.Method = "GET"
		ret.RespCookies = response.Cookies()
		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			ret.HttpError = err.Error()
			return ret
		} else {
			response.Body.Close()
			ret.Raw = string(body)
			return ret
		}
	}


}


func HttpPost (strUrl string, p Any, headers map[string]string, cookies map[string]string) *HttpResp {
	client := &http.Client{}
	ret := &HttpResp{}
	httpParams := url.Values{}
	var request *http.Request
	var err error
	if params, ok := p.(map[string]Any) ; ok {
		for k, v := range params {
			switch realValue := v.(type) {
			case string:
				httpParams.Set(k, realValue)
			case int:
				httpParams.Set(k, strconv.Itoa(realValue))
			case []byte:
				httpParams.Set(k, string(realValue))
			case byte:
				httpParams.Set(k, string(realValue))
			}
		}
		request, err = http.NewRequest("POST", strUrl, strings.NewReader(httpParams.Encode()))
	} else {
		if strParam, ok := p.(string) ; ok {
			request, err = http.NewRequest("POST", strUrl, strings.NewReader(strParam))
		} else {
			if byteParam, ok := p.([]byte) ; ok {
				request, err = http.NewRequest("POST", strUrl, strings.NewReader(string(byteParam)))
			}
		}
	}
	if err != nil {
		ret.HttpError = err.Error()
		return ret
	}
	if nil != headers {
		for k, v := range headers {
			request.Header.Add(k, v)
		}
	}
	if 0 >= len(request.Header.Get("User-Agent")) {
		request.Header.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/61.0.3163.100 Safari/537.36")
	}

	strCookie := ""
	if nil != cookies {
		for k, v := range cookies {
			strCookie += k + "=" + v + ";"
		}
	}
	request.Header.Add("Cookie", strCookie)
	ret.Params = &p
	ret.ReqHeaders = &request.Header
	ret.ReqCookies = &cookies
	response, err := client.Do(request)
	if nil != err {
		ret.HttpError = err.Error()
		return ret
	} else {
		ret.HttpCode = response.StatusCode
		ret.Url = strUrl
		ret.Method = "GET"
		ret.RespCookies = response.Cookies()
		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			ret.HttpError = err.Error()
			return ret
		} else {
			response.Body.Close()
			ret.Raw = string(body)
			return ret
		}
	}


}