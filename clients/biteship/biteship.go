package biteship

import (
	"book-nest/config"
	"encoding/json"
	"errors"
	"fmt"
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
	fmt.Println("KEY: ", b.ServerKey)
	fmt.Println("BODY: ", string(body))

	if !biteshipCourierResponse.Success {
		return nil, errors.New(biteshipCourierResponse.Error)
	}

	return biteshipCourierResponse, nil
}
