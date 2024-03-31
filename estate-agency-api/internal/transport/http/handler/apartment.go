package handler

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"gilab.com/estate-agency-api/internal/entity"
	httpModel "gilab.com/estate-agency-api/internal/transport/http/model"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

const (
	apartmentsURL = "/apartments"
	apartmentURL  = "/apartments/:apartment_id"
)

type apartmentHandler struct {
	usecase  Usecase
	validate *validator.Validate
	logger   *slog.Logger
}

func NewApartmentHandler(usecase Usecase, logger *slog.Logger) *apartmentHandler {
	return &apartmentHandler{usecase: usecase, validate: validator.New(), logger: logger}
}

func (h *apartmentHandler) Register(router *gin.Engine) {
	router.GET(apartmentsURL, h.GetApartments)
	router.GET(apartmentURL, h.GetApartment)
	router.POST(apartmentsURL, h.CreateApartment)
	router.PATCH(apartmentURL, h.UpdateApartment)
	router.DELETE(apartmentURL, h.DeleteApartment)
}

func (h *apartmentHandler) GetApartments(ctx *gin.Context) {
	const op = "handler.GetApartments"

	log := h.logger.With(slog.String("op", op))

	const page_size = 10
	page, err := strconv.Atoi(ctx.DefaultQuery("page", "0"))
	if err != nil {
		log.Info("page id wrong", slog.Int("id", page))
		ctx.JSON(http.StatusBadRequest, gin.H{"err": "error id"})
		return
	}

	ct := context.WithValue(context.Background(), "logger", h.logger)
	apartments, err := h.usecase.GetAllApartment(ct, page, page_size)

	if err != nil {
		log.Info("page not found", slog.Int("id", page))
		ctx.JSON(http.StatusNotFound, gin.H{"err": "not found"})
		return
	}

	log.Info("page found", slog.Int("id", page))

	ctx.JSON(http.StatusOK, apartments)
}

func (h *apartmentHandler) GetApartment(ctx *gin.Context) {
	const op = "handler.GetApartment"

	log := h.logger.With(slog.String("op", op))

	id, err := strconv.Atoi(ctx.Param("apartment_id"))
	if err != nil {
		log.Info("id wrong", slog.Int("id", id))
		ctx.JSON(http.StatusBadRequest, gin.H{"err": "error id"})
		return
	}

	ct := context.WithValue(context.Background(), "logger", h.logger)
	apartment, realtor, err := h.usecase.GetApartmentByID(ct, id)
	if err != nil {
		log.Info("not found", slog.Int("id", id))
		ctx.JSON(http.StatusInternalServerError, gin.H{"err": "not found"})
		return
	}

	realtorView := httpModel.RealtorView{
		FirstName: realtor.FirstName,
		LastName:  realtor.LastName,
		Phone:     realtor.Phone,
		Email:     realtor.Email,
		Rating:    realtor.Rating,
	}

	apartmentView := httpModel.ApartmentView{
		Apartment:   *apartment,
		RealtorView: realtorView,
	}

	ctx.JSON(http.StatusOK, apartmentView)
	ctx.Header("Content-Length", fmt.Sprintf("%d", 64<<20))
	ctx.FileAttachment(fmt.Sprintf("./../../internal/images/apartment/%d.png", id), "photo_apartment")
	ctx.FileAttachment(fmt.Sprintf("./../../internal/images/realtor/%d.png", apartment.IDRealtor), "photo_realtor")
}

func (h *apartmentHandler) CreateApartment(ctx *gin.Context) {
	const op = "handler.CreateApartment"

	log := h.logger.With(slog.String("op", op))

	var apartment entity.Apartment

	if err := ctx.Bind(&apartment); err != nil {
		log.Info("invalid request")
		ctx.JSON(http.StatusBadRequest, gin.H{"err": "invalid request"})
		return
	}

	if err := h.validate.Struct(apartment); err != nil {
		log.Info("bad validate", slog.Any("apartment", apartment))
		ctx.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}

	file, err := ctx.FormFile("photo")
	if err != nil {
		log.Info("bad photo")
		ctx.JSON(http.StatusBadRequest, gin.H{"err": "invalid request"})
		return
	}

	file_format := strings.Split(file.Header["Content-Type"][0], "/")[1]
	if file_format != "png" {
		log.Info("not png")
		ctx.JSON(http.StatusBadRequest, gin.H{"err": "file only png"})
		return
	}

	ct := context.WithValue(context.Background(), "logger", h.logger)
	id, err := h.usecase.CreateApartment(ct, &apartment)

	if err != nil {
		log.Info("failed to create", "err", err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{"err": "failed to create"})
		return
	}

	ctx.SaveUploadedFile(file, fmt.Sprintf("./../../internal/images/apartment/%d.%s", id, file_format))

	ctx.JSON(http.StatusCreated, gin.H{"apartment_id": id})
}

func (h *apartmentHandler) UpdateApartment(ctx *gin.Context) {
	const op = "handler.UpdateApartment"

	log := h.logger.With("op", op)

	var apartment entity.Apartment
	id, err := strconv.Atoi(ctx.Param("apartment_id"))
	if err != nil {
		log.Info("invalid request")
		ctx.JSON(http.StatusBadRequest, gin.H{"err": "invalid request"})
		return
	}

	if err = ctx.Bind(&apartment); err != nil {
		log.Info("invalid request")
		ctx.JSON(http.StatusBadRequest, gin.H{"err": "invalid request"})
		return
	}

	file, err_file := ctx.FormFile("photo")
	if err_file != nil && err_file.Error() != "http: no such file" {
		ctx.JSON(http.StatusBadRequest, gin.H{"err": "invalid request"})
		return
	}

	if err = h.validate.Struct(apartment); err != nil {
		log.Info("bad validate")
		ctx.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}

	if err_file == nil {
		file_format := strings.Split(file.Header["Content-Type"][0], "/")[1]
		if file_format != "png" {
			log.Info("not png")
			ctx.JSON(http.StatusBadRequest, gin.H{"err": "file only png"})
			return
		}

		ctx.SaveUploadedFile(file, fmt.Sprintf("./../../internal/images/realtor/%d.%s", id, file_format))
	}

	ct := context.WithValue(context.Background(), "logger", h.logger)
	aff, err := h.usecase.UpdateApartment(ct, id, &apartment)
	if err != nil {
		log.Info("failed to update", "err", err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{"err": "failed to update"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"affected": aff})
}

func (h *apartmentHandler) DeleteApartment(ctx *gin.Context) {
	const op = "handler.DeleteApartment"

	log := h.logger.With("op", op)

	id, err := strconv.Atoi(ctx.Param("apartment_id"))
	if err != nil {
		log.Info("error id", slog.Int("id", id))
		ctx.JSON(http.StatusBadRequest, gin.H{"err": "error id"})
		return
	}

	ct := context.WithValue(context.Background(), "logger", h.logger)
	err = h.usecase.DeleteApartment(ct, id)
	if err != nil {
		log.Info("not deleted", slog.Int("id", id), "err", err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{"err": "not deleted"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"msg": "deleted"})
}
