package v1

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"gilab.com/estate-agency-api/internal/domain/entity"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

const (
	realtorsURL = "/realtors"
	realtorURL  = "/realtors/:realtor_id"
)

type RealtorUsecase interface {
	GetAllRealtor(ctx context.Context, page int, pageSize int) (realtors []*entity.Realtor, err error)
	GetRealtorByID(ctx context.Context, id int) (realtor *entity.Realtor, err error)
	CreateRealtor(ctx context.Context, realtor *entity.Realtor) (id int64, err error)
	UpdateRealtor(ctx context.Context, id int, realtor *entity.Realtor) (aff int64, err error)
	DeleteRealtor(ctx context.Context, id int) error
}

type realtorHandler struct {
	realtorUsecase RealtorUsecase
	validate       *validator.Validate
}

func NewRealtorHandler(realtorUsecase RealtorUsecase) *realtorHandler {
	return &realtorHandler{realtorUsecase: realtorUsecase, validate: validator.New()}
}

func (h *realtorHandler) Register(router *gin.Engine) {
	router.GET(realtorsURL, h.GetRealtors)
	router.GET(realtorURL, h.GetRealtor)
	router.POST(realtorsURL, h.CreateRealtor)
	router.PATCH(realtorURL, h.UpdateRealtor)
	router.DELETE(realtorURL, h.DeleteRealtor)
}

func (h *realtorHandler) GetRealtors(ctx *gin.Context) {
	const page_size = 10
	page, err := strconv.Atoi(ctx.DefaultQuery("page", "0"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}

	realtors, err := h.realtorUsecase.GetAllRealtor(context.Background(), page, page_size)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, realtors)
}

func (h *realtorHandler) GetRealtor(ctx *gin.Context) {
	var realtor *entity.Realtor
	id, err := strconv.Atoi(ctx.Param("realtor_id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}

	realtor, err = h.realtorUsecase.GetRealtorByID(context.Background(), id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, realtor)
	ctx.Header("Content-Length", fmt.Sprintf("%d", 32<<20))
	ctx.FileAttachment(fmt.Sprintf("./../../internal/images/realtor/%d.png", id), "photo")
}

func (h *realtorHandler) CreateRealtor(ctx *gin.Context) {
	var realtor entity.Realtor

	if err := ctx.Bind(&realtor); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}

	if err := h.validate.Struct(realtor); err != nil {
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

	id, err := h.realtorUsecase.CreateRealtor(context.Background(), &realtor)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}

	ctx.SaveUploadedFile(file, fmt.Sprintf("./../../internal/images/realtor/%d.%s", id, file_format))

	ctx.JSON(http.StatusCreated, gin.H{"realtor_id": id})
}

func (h *realtorHandler) UpdateRealtor(ctx *gin.Context) {
	var realtor entity.Realtor
	id, err := strconv.Atoi(ctx.Param("realtor_id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}

	if err = ctx.Bind(&realtor); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}

	if err = h.validate.Struct(realtor); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}

	aff, err := h.realtorUsecase.UpdateRealtor(context.Background(), id, &realtor)
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

func (h *realtorHandler) DeleteRealtor(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("realtor_id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}

	err = h.realtorUsecase.DeleteRealtor(context.Background(), id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"msg": "deleted"})
}
