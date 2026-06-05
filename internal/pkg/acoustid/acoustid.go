package acoustid

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/google/go-querystring/query"
	"github.com/mjdevelops/tunes/internal/pkg/acoustid/models"
)

type AcoustIdApi struct {
	apiKey string
}

const acoustIdLookupBase string = "https://api.acoustid.org/v2/lookup"

var ErrNoApiKey = errors.New("no api key provided")

func NewAcoustIdApi(apiKey string) *AcoustIdApi {
	return &AcoustIdApi{
		apiKey,
	}
}

func (a *AcoustIdApi) Lookup(req *models.AcoustIdLookupParams) (*models.AcoustIdResponse, error) {
	if a.apiKey == "" {
		return nil, ErrNoApiKey
	}

	v, err := query.Values(req)
	if err != nil {
		return nil, err
	}

	reqUrl, err := url.Parse(acoustIdLookupBase)
	if err != nil {
		return nil, err
	}

	v.Add("client", a.apiKey)

	reqUrl.RawQuery = v.Encode()

	res, err := http.Get(reqUrl.String())
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	acoustRes := &models.AcoustIdResponse{}

	err = json.NewDecoder(res.Body).Decode(acoustRes)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	return acoustRes, nil
}
