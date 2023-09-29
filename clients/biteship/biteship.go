package biteship

import (
	"book-nest/config"
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

type Biteship struct {
	BaseUrl   string
	ServerKey string
}

func NewBiteshipClient() *Biteship {
	return &Biteship{
		BaseUrl:   config.Cfg.BiteshipConf.BaseUrl,
		ServerKey: config.Cfg.BiteshipConf.ServerKey,
	}
}

func (b *Biteship) GetListCourier() (*BiteshipCourierResponse, error) {
	url := b.BaseUrl + "/v1/couriers"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", b.ServerKey)
	req.Header.Add("content-type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var biteshipCourierResponse *BiteshipCourierResponse
	err = json.Unmarshal(body, &biteshipCourierResponse)
	if err != nil {
		return nil, err
	}

	if !biteshipCourierResponse.Success {
		return nil, errors.New(biteshipCourierResponse.Error)
	}

	return biteshipCourierResponse, nil
}

func (b *Biteship) CheckRates(payload *BiteshipCheckRatesRequest) (*BiteshipCheckRatesResponse, error) {
	url := b.BaseUrl + "/v1/rates/couriers"

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", b.ServerKey)
	req.Header.Add("accept", "application/json")
	req.Header.Add("content-type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var biteshipCheckRatesResponse *BiteshipCheckRatesResponse
	err = json.Unmarshal(body, &biteshipCheckRatesResponse)
	if err != nil {
		return nil, err
	}

	if !biteshipCheckRatesResponse.Success {
		return nil, errors.New(biteshipCheckRatesResponse.Error)
	}

	return biteshipCheckRatesResponse, nil
}
