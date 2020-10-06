package controller

import (
	"context"
	"errors"
	"fmt"
	"maga-auctions/api/handler"
	"maga-auctions/api/helper/filters"
	"maga-auctions/entity"
	"maga-auctions/vehicle"
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
	All(c *gin.Context)
	ByID(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
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
		handler.ResponseError(err, c)
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

func buildFilters(c *gin.Context, fs *[]filters.Filter) error {
	fb := c.Query("brand")
	if fb != "" {
		*fs = append(*fs, filters.NewVehicleBrand(fb))
	}

	md := c.Query("model")
	if md != "" {
		*fs = append(*fs, filters.NewVehicleModel(md))
	}

	mMin, err := strconv.ParseInt(c.DefaultQuery("manufacturingYearMin", "0"), 10, 32)
	if err != nil {
		return errors.New("manufacturing year min is invalid")
	}

	mMax, err := strconv.ParseInt(c.DefaultQuery("manufacturingYearMax", "0"), 10, 32)
	if err != nil {
		return errors.New("manufacturing year max is invalid")
	}

	if mMin > 0 && mMax > 0 {
		f, err := filters.NewVehicleYearBetween(int(mMin), int(mMax))

		if err != nil {
			return err
		}

		*fs = append(*fs, f)
	}

	mfy, err := strconv.ParseInt(c.DefaultQuery("manufacturingYear", "0"), 10, 32)
	if err != nil {
		return errors.New("manufacturing year is invalid")
	}

	mdy, err := strconv.ParseInt(c.DefaultQuery("modelYear", "0"), 10, 32)
	if err != nil {
		return errors.New("model year is invalid")
	}

	if mfy > 0 && mdy > 0 {
		f, err := filters.NewVehicleYear(int(mdy), int(mfy))

		if err != nil {
			return errors.New(err.Error())
		}

		*fs = append(*fs, f)
	}

	return nil
}

func (v vehicleCtrl) All(c *gin.Context) {
	var fs []filters.Filter
	err := buildFilters(c, &fs)

	if err != nil {
		handler.ResponseError(handler.BadRequest{Message: err.Error()}, c)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	items, err := v.srv.All(ctx, fs, c.Query("bidOrder"))

	if err != nil {
		handler.ResponseError(err, c)
		return
	}

	handler.ResponseSuccess(200, items, c)
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

func (v vehicleCtrl) Update(c *gin.Context) {
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

	var ve entity.Vehicle
	err = c.BindJSON(&ve)

	if err != nil {
		handler.ResponseError(
			handler.BadRequest{
				Message: "body is invalid",
			},
			c,
		)
		return
	}

	ve.ID = int(id)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	err = v.srv.Update(ctx, &ve)

	if err != nil {
		handler.ResponseError(err, c)
		return
	}

	res := response{
		Vehicle: ve,
		Links: []Link{
			{
				Relation:     "self",
				RelationType: "GET",
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

func (v vehicleCtrl) Delete(c *gin.Context) {
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

	err = v.srv.Delete(ctx, int(id))

	if err != nil {
		handler.ResponseError(err, c)
		return
	}

	handler.ResponseSuccess(200, nil, c)
}
