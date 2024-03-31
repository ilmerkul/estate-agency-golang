package handler

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"gilab.com/estate-agency-api/internal/entity"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

const (
	realtorsURL = "/realtors"
	realtorURL  = "/realtors/:realtor_id"
)

type realtorHandler struct {
	usecase  Usecase
	validate *validator.Validate

	logger *slog.Logger
}

func NewRealtorHandler(usecase Usecase, logger *slog.Logger) *realtorHandler {
	return &realtorHandler{usecase: usecase, validate: validator.New(), logger: logger}
}

func (h *realtorHandler) Register(router *gin.Engine) {
	router.GET(realtorsURL, h.GetRealtors)
	router.GET(realtorURL, h.GetRealtor)
	router.POST(realtorsURL, h.CreateRealtor)
	router.PATCH(realtorURL, h.UpdateRealtor)
	router.DELETE(realtorURL, h.DeleteRealtor)
}

func (h *realtorHandler) GetRealtors(ctx *gin.Context) {
	const op = "handler.GetRealtors"

	const page_size = 10
	page, err := strconv.Atoi(ctx.DefaultQuery("page", "0"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"err": "error page"})
		return
	}

	ct := context.WithValue(context.Background(), "logger", h.logger)
	realtors, err := h.usecase.GetAllRealtor(ct, page, page_size)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"err": "failed to get"})
		return
	}

	ctx.JSON(http.StatusOK, realtors)
}

func (h *realtorHandler) GetRealtor(ctx *gin.Context) {
	const op = "handler.GetRealtor"

	var realtor *entity.Realtor
	id, err := strconv.Atoi(ctx.Param("realtor_id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"err": "error id"})
		return
	}

	ct := context.WithValue(context.Background(), "logger", h.logger)
	realtor, err = h.usecase.GetRealtorByID(ct, id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"err": "failed to get"})
		return
	}

	ctx.JSON(http.StatusOK, realtor)
	ctx.Header("Content-Length", fmt.Sprintf("%d", 32<<20))
	ctx.FileAttachment(fmt.Sprintf("./../../internal/images/realtor/%d.png", id), "photo")
}

func (h *realtorHandler) CreateRealtor(ctx *gin.Context) {
	const op = "handler.CreateRealtor"

	var realtor entity.Realtor

	if err := ctx.Bind(&realtor); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"err": "invalid request"})
		return
	}

	if err := h.validate.Struct(realtor); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}

	file, err := ctx.FormFile("photo")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"err": "error photo"})
		return
	}

	file_format := strings.Split(file.Header["Content-Type"][0], "/")[1]
	if file_format != "png" {
		ctx.JSON(http.StatusBadRequest, gin.H{"err": "file only png"})
		return
	}

	ct := context.WithValue(context.Background(), "logger", h.logger)
	id, err := h.usecase.CreateRealtor(ct, &realtor)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"err": "failed to create"})
		return
	}

	ctx.SaveUploadedFile(file, fmt.Sprintf("./../../internal/images/realtor/%d.%s", id, file_format))

	ctx.JSON(http.StatusCreated, gin.H{"realtor_id": id})
}

func (h *realtorHandler) UpdateRealtor(ctx *gin.Context) {
	const op = "handler.UpdateRealtor"

	var realtor entity.Realtor
	id, err := strconv.Atoi(ctx.Param("realtor_id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"err": "error id"})
		return
	}

	if err = ctx.Bind(&realtor); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"err": "invalid request"})
		return
	}

	if err = h.validate.Struct(realtor); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}

	file, err := ctx.FormFile("photo")
	if err != nil && err.Error() != "http: no such file" {
		ctx.JSON(http.StatusBadRequest, gin.H{"err": "error photo"})
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

	ct := context.WithValue(context.Background(), "logger", h.logger)
	aff, err := h.usecase.UpdateRealtor(ct, id, &realtor)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"err": "failed to update"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"affected": aff})
}

func (h *realtorHandler) DeleteRealtor(ctx *gin.Context) {
	const op = "handler.DeleteRealtor"

	id, err := strconv.Atoi(ctx.Param("realtor_id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"err": "error id"})
		return
	}

	ct := context.WithValue(context.Background(), "logger", h.logger)
	err = h.usecase.DeleteRealtor(ct, id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"err": "failed to delete"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"msg": "deleted"})
}
