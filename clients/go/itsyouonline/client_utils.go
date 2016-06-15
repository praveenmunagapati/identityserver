package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

func encodeBody(data interface{}) (io.Reader, error) {
	b, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(b), nil
}

func doReqWithBody(method, urlStr string, data interface{}, client http.Client, headers map[string]interface{}, qsParam string) (*http.Response, error) {
	body, err := encodeBody(data)
	if err != nil {
		return nil, err
	}

	// create the request
	req, err := http.NewRequest(method, urlStr+qsParam, body)
	if err != nil {
		return nil, err
	}

	for k, v := range headers {
		req.Header.Set(k, fmt.Sprintf("%v", v))
	}

	return client.Do(req)
}

func buildQueryString(data map[string]interface{}) string {
	if len(data) == 0 {
		return ""
	}

	baseQuery := "?"
	for k, v := range data {
		baseQuery += k + "=" + fmt.Sprint(v) + "&"
	}

	return baseQuery[:len(baseQuery)-1]
}

//Date represent RFC3399 date
type Date time.Time

//MarshalJSON override marshalJSON
func (t *Date) MarshalJSON() ([]byte, error) {
	return []byte(time.Time(*t).Format(`"` + time.RFC3339 + `"`)), nil
}

//MarshalText override marshalText
func (t *Date) MarshalText() ([]byte, error) {
	return []byte(time.Time(*t).Format(`"` + time.RFC3339 + `"`)), nil
}

//UnmarshalJSON override unmarshalJSON
func (t *Date) UnmarshalJSON(b []byte) error {
	ts, err := time.Parse(`"`+time.RFC3339+`"`, string(b))
	if err != nil {
		return err
	}

	*t = Date(ts)
	return nil
}

//UnmarshalText override unmarshalText
func (t *Date) UnmarshalText(b []byte) error {
	ts, err := time.Parse(`"`+time.RFC3339+`"`, string(b))
	if err != nil {
		return err
	}

	*t = Date(ts)
	return nil
}

func (t *Date) String() string {
	return time.Time(*t).String()
}
