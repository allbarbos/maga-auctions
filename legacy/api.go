package legacy

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"maga-auctions/entity"
	"maga-auctions/env"
	"net/http"
)

var (
	// Client is the web client
	Client HTTPClient
	// APIURI is legacy service url
	APIURI, method string
)

func init() {
	Client = &http.Client{}
	APIURI = env.Vars.APILegacy.URI
	method = "POST"
}

// API contract
type API interface {
	Get(ctx context.Context) ([]VehicleLegacy, error)
	Create(ctx context.Context, vehicle entity.Vehicle) (*VehicleLegacy, error)
	Update(ctx context.Context, vehicle entity.Vehicle) (*VehicleLegacy, error)
	// Delete(ctx context.Context, id int) (*http.Response, error)
}

// HTTPClient is the web client
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type srv struct{}

type VehicleLegacy struct {
	ID             int     `json:"ID,omitempty"`
	DataLance      string  `json:"DATALANCE"`
	Lote           string  `json:"LOTE"`
	CodigoControle string  `json:"CODIGOCONTROLE"`
	Marca          string  `json:"MARCA"`
	Modelo         string  `json:"MODELO"`
	AnoFabricacao  int     `json:"ANOFABRICACAO"`
	AnoModelo      int     `json:"ANOMODELO"`
	ValorLance     float32 `json:"VALORLANCE"`
	UsuarioLance   string  `json:"USUARIOLANCE"`
}

type body struct {
	Operacao string        `json:"OPERACAO,omitempty"`
	Veiculo  VehicleLegacy `json:"VEICULO,omitempty"`
}

// NewAPI returns a planet service instance
func NewAPI() API {
	return &srv{}
}

func makeRequest(body interface{}) (*http.Request, error) {
	b, err := json.Marshal(body)
	payload := bytes.NewReader(b)

	req, err := http.NewRequest(method, APIURI, payload)

	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")

	return req, nil
}

func (s srv) Get(ctx context.Context) ([]VehicleLegacy, error) {
	req, err := makeRequest(body{Operacao: "consultar"})

	if err != nil {
		return nil, err
	}

	res, err := Client.Do(req)

	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	items := []VehicleLegacy{}
	err = json.NewDecoder(res.Body).Decode(&items)

	if err != nil {
		return nil, err
	}

	return items, nil
}

func (s srv) Create(ctx context.Context, vehicle entity.Vehicle) (*VehicleLegacy, error) {
	b := body{
		Operacao: "criar",
		Veiculo: VehicleLegacy{
			Lote:           vehicle.Lot.ID,
			CodigoControle: vehicle.Lot.VehicleLotID,
			Marca:          vehicle.Brand,
			Modelo:         vehicle.Model,
			AnoFabricacao:  vehicle.ManufacturingYear,
			AnoModelo:      vehicle.ModelYear,
			DataLance:      "-",
			UsuarioLance:   "-",
		},
	}

	req, err := makeRequest(b)

	if err != nil {
		return nil, err
	}

	res, err := Client.Do(req)

	if err != nil {
		return nil, err
	}

	body, _ := ioutil.ReadAll(res.Body)

	v := VehicleLegacy{}
	err = json.Unmarshal(body, &v)

	if err != nil {
		return nil, err
	}

	return &v, nil
}

func (s srv) Update(ctx context.Context, vehicle entity.Vehicle) (*VehicleLegacy, error) {
	b := body{
		Operacao: "alterar",
		Veiculo: VehicleLegacy{
			ID:             vehicle.ID,
			Marca:          vehicle.Brand,
			Modelo:         vehicle.Model,
			AnoFabricacao:  vehicle.ManufacturingYear,
			AnoModelo:      vehicle.ModelYear,
			Lote:           vehicle.Lot.ID,
			CodigoControle: vehicle.Lot.VehicleLotID,
			DataLance:      "-",
			UsuarioLance:   "-",
		},
	}

	req, err := makeRequest(b)

	if err != nil {
		return nil, err
	}

	res, err := Client.Do(req)

	if err != nil {
		return nil, err
	}

	if res.StatusCode == 200 {
		return &b.Veiculo, nil
	}

	return nil, errors.New("error when updating in legacy api")
}

// func (s srv) Delete(ctx context.Context, id int) (*http.Response, error) {
// 	b := body{
// 		Operacao: "apagar",
// 		Veiculo:  VehicleLegacy{ID: id},
// 	}

// 	req, err := makeRequest(b)

// 	if err != nil {
// 		return nil, err
// 	}

// 	res, err := Client.Do(req)

// 	if err != nil {
// 		return nil, err
// 	}

// 	return res, nil
// }
