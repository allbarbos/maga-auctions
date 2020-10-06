package legacy

import (
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"maga-auctions/entity"
	"maga-auctions/utils"
	"net/http"
	"regexp"
	"time"
)

var (
	// Client is the web client
	Client utils.HTTPClient
	// APIURI is legacy service url
	APIURI, method string
)

func init() {
	Client = &http.Client{}
	APIURI = utils.EnvVars.APILegacy.URI
	method = "POST"
}

// API contract
type API interface {
	Get(ctx context.Context) ([]entity.Vehicle, error)
	Create(ctx context.Context, vehicle *entity.Vehicle) error
	Update(ctx context.Context, vehicle *entity.Vehicle) error
	Delete(ctx context.Context, id int) error
}

type srv struct{}

// VehicleLegacy is legacy entity
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

func (s srv) Get(ctx context.Context) ([]entity.Vehicle, error) {
	req, err := utils.MakeRequest(method, APIURI, body{Operacao: "consultar"})

	if err != nil {
		return nil, err
	}

	res, err := Client.Do(req)

	if err != nil {
		return nil, err
	}

	if res.StatusCode != 200 {
		return nil, errors.New("an error occurred while requesting the legacy api")
	}

	defer res.Body.Close()

	items := []VehicleLegacy{}
	err = json.NewDecoder(res.Body).Decode(&items)

	if err != nil {
		return nil, err
	}

	v := []entity.Vehicle{}

	for i := 0; i < len(items); i++ {
		l := items[i]

		layout := "02/01/2006 - 15:04"
		bidDate, err := time.Parse(layout, l.DataLance)

		if err != nil {
			continue
		}

		item := entity.Vehicle{
			ID:                l.ID,
			Brand:             l.Marca,
			Model:             l.Modelo,
			ModelYear:         l.AnoModelo,
			ManufacturingYear: l.AnoFabricacao,
			Lot: entity.Lot{
				ID:           l.Lote,
				VehicleLotID: l.CodigoControle,
			},
			Bid: entity.Bid{
				Date:  bidDate,
				User:  l.UsuarioLance,
				Value: l.ValorLance,
			},
		}

		v = append(v, item)
	}

	return v, nil
}

func (s srv) Create(ctx context.Context, vehicle *entity.Vehicle) error {
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

	req, err := utils.MakeRequest(method, APIURI, b)

	if err != nil {
		return err
	}

	res, err := Client.Do(req)

	if err != nil {
		return err
	}

	if res.StatusCode != 200 {
		return errors.New("an error occurred while requesting the legacy api")
	}

	body, _ := ioutil.ReadAll(res.Body)

	err = json.Unmarshal(body, &vehicle)

	if err != nil {
		return err
	}

	return nil
}

func (s srv) Update(ctx context.Context, vehicle *entity.Vehicle) error {
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

	req, err := utils.MakeRequest(method, APIURI, b)

	if err != nil {
		return err
	}

	res, err := Client.Do(req)

	if err != nil {
		return err
	}

	if res.StatusCode != 200 {
		return errors.New("an error occurred while requesting the legacy api")
	}

	body, _ := ioutil.ReadAll(res.Body)
	isError, _ := regexp.MatchString("nao encontrado", string(body))

	if isError {
		return errors.New("id not found")
	}

	return nil
}

func (s srv) Delete(ctx context.Context, id int) error {
	b := body{
		Operacao: "apagar",
		Veiculo:  VehicleLegacy{ID: id},
	}

	req, err := utils.MakeRequest(method, APIURI, b)

	if err != nil {
		return err
	}

	res, err := Client.Do(req)

	if err != nil {
		return err
	}

	if res.StatusCode != 200 {
		return errors.New("an error occurred while requesting the legacy api")
	}

	body, _ := ioutil.ReadAll(res.Body)
	isError, _ := regexp.MatchString("nao encontrado", string(body))

	if isError {
		return errors.New("id not found")
	}

	return nil
}
