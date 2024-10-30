package v1

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"gilab.com/estate-agency-api/internal/domain/entity"
	"gilab.com/estate-agency-api/internal/transport/http/dto"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

const (
	apartmentsURL = "/apartments"
	apartmentURL  = "/apartments/:apartment_id"
)

type ApartmentUsecase interface {
	GetAllApartment(ctx context.Context, page int, pageSize int) (apartments []*entity.Apartment, err error)
	GetApartmentByID(ctx context.Context, id int) (apartment *entity.Apartment, realtor *entity.Realtor, err error)
	CreateApartment(ctx context.Context, apartment *entity.Apartment) (id int64, err error)
	UpdateApartment(ctx context.Context, id int, apartment *entity.Apartment) (aff int64, err error)
	DeleteApartment(ctx context.Context, id int) error
}

type apartmentHandler struct {
	apartmentUsecase ApartmentUsecase
	validate         *validator.Validate
}

func NewApartmentHandler(apartmentUsecase ApartmentUsecase) *apartmentHandler {
	return &apartmentHandler{apartmentUsecase: apartmentUsecase, validate: validator.New()}
}

func (h *apartmentHandler) Register(router *gin.Engine) {
	router.GET(apartmentsURL, h.GetApartments)
	router.GET(apartmentURL, h.GetApartment)
	router.POST(apartmentsURL, h.CreateApartment)
	router.PATCH(apartmentURL, h.UpdateApartment)
	router.DELETE(apartmentURL, h.DeleteApartment)
}

func (h *apartmentHandler) GetApartments(ctx *gin.Context) {
	const page_size = 10
	page, err := strconv.Atoi(ctx.DefaultQuery("page", "0"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}

	apartments, err := h.apartmentUsecase.GetAllApartment(context.Background(), page, page_size)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, apartments)
}

func (h *apartmentHandler) GetApartment(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("apartment_id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}

	apartment, realtor, err := h.apartmentUsecase.GetApartmentByID(context.Background(), id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}

	realtorView := dto.RealtorView{
		FirstName: realtor.FirstName,
		LastName:  realtor.LastName,
		Phone:     realtor.Phone,
		Email:     realtor.Email,
		Rating:    realtor.Rating,
	}

	apartmentView := dto.ApartmentView{
		Apartment:   *apartment,
		RealtorView: realtorView,
	}

	ctx.JSON(http.StatusOK, apartmentView)
	ctx.Header("Content-Length", fmt.Sprintf("%d", 64<<20))
	ctx.FileAttachment(fmt.Sprintf("./../../internal/images/aprtment/%d.png", id), "photo_apartment")
	ctx.FileAttachment(fmt.Sprintf("./../../internal/images/realtor/%d.png", apartment.IDRealtor), "photo_realtor")
}

func (h *apartmentHandler) CreateApartment(ctx *gin.Context) {
	var apartment entity.Apartment

	if err := ctx.Bind(&apartment); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}

	if err := h.validate.Struct(apartment); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}

	file, err := ctx.FormFile("photo")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}

	file_format := strings.Split(file.Header["Content-Type"][0], "/")[1]
	if file_format != "png" {
		ctx.JSON(http.StatusBadRequest, gin.H{"err": "file only png"})
		return
	}

	id, err := h.apartmentUsecase.CreateApartment(context.Background(), &apartment)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}

	ctx.SaveUploadedFile(file, fmt.Sprintf("./../../internal/images/apartment/%d.%s", id, file_format))

	ctx.JSON(http.StatusCreated, gin.H{"apartment_id": id})
}

func (h *apartmentHandler) UpdateApartment(ctx *gin.Context) {
	var apartment entity.Apartment
	id, err := strconv.Atoi(ctx.Param("apartment_id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}

	if err = ctx.Bind(&apartment); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}

	if err = h.validate.Struct(apartment); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}

	aff, err := h.apartmentUsecase.UpdateApartment(context.Background(), id, &apartment)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}

	file, err := ctx.FormFile("photo")
	if err != nil && err.Error() != "http: no such file" {
		ctx.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}

	if err == nil {
		file_format := strings.Split(file.Header["Content-Type"][0], "/")[1]
		if file_format != "png" {
			ctx.JSON(http.StatusBadRequest, gin.H{"err": "file only png"})
			return
		}

		ctx.SaveUploadedFile(file, fmt.Sprintf("./../../internal/images/realtor/%d.%s", id, file_format))
	}

	ctx.JSON(http.StatusOK, gin.H{"affected": aff})
}

func (h *apartmentHandler) DeleteApartment(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("apartment_id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}

	err = h.apartmentUsecase.DeleteApartment(context.Background(), id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"msg": "deleted"})
}
