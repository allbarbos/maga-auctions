package controller

import (
	"api-facade/api/handler"
	"api-facade/entity"
	"api-facade/vehicle"
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type response struct {
	Vehicle entity.Vehicle `json:"vehicle"`
	Links   []Link         `json:"links"`
}

type Link struct {
	URI          string `json:"uri"`
	Relation     string `json:"rel"`
	RelationType string `json:"type"`
}

// VehicleController contract
type VehicleController interface {
	Create(c *gin.Context)
	ByID(c *gin.Context)
}

type vehicleCtrl struct {
	srv vehicle.Service
}

// NewVehicle controller
func NewVehicle(srv vehicle.Service) VehicleController {
	return &vehicleCtrl{
		srv: srv,
	}
}

// Create vehicle
func (v vehicleCtrl) Create(c *gin.Context) {
	var ve entity.Vehicle
	err := c.BindJSON(&ve)

	if err != nil {
		handler.ResponseError(
			handler.BadRequest{
				Message: "body is invalid",
			},
			c,
		)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	registered, err := v.srv.Create(ctx, ve)

	if err != nil {
		handler.ResponseError(
			handler.InternalServer{
				Message: "deu ruim",
			},
			c,
		)
		return
	}

	c.Header("Location", fmt.Sprintf("%s%s/%d", c.Request.Host, c.Request.RequestURI, registered.ID))
	uri := fmt.Sprintf("%s/%d", c.Request.RequestURI, registered.ID)
	links := []Link{
		{
			Relation:     "self",
			RelationType: "GET",
			URI:          uri,
		},
		{
			Relation:     "self",
			RelationType: "PUT",
			URI:          uri,
		},
		{
			Relation:     "self",
			RelationType: "DELETE",
			URI:          uri,
		},
	}

	res := response{Vehicle: *registered, Links: links}

	handler.ResponseSuccess(201, res, c)
}

func (v vehicleCtrl) ByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 32)

	if err != nil {
		handler.ResponseError(
			handler.BadRequest{
				Message: "id is invalid",
			},
			c,
		)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	ve, err := v.srv.ByID(ctx, int(id))

	if err != nil {
		handler.ResponseError(err, c)
		return
	}

	res := response{
		Vehicle: *ve,
		Links: []Link{
			{
				Relation:     "self",
				RelationType: "PUT",
				URI:          c.Request.RequestURI,
			},
			{
				Relation:     "self",
				RelationType: "DELETE",
				URI:          c.Request.RequestURI,
			},
		},
	}

	handler.ResponseSuccess(200, res, c)
}
