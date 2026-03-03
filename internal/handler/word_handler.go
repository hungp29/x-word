package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hupham/x-word/internal/model"
)

type WordService interface {
	GetWord(string) (*model.Word, error)
	GetWords([]string) ([]*model.Word, error)
}

type WordHandler struct {
	service interface {
		GetWord(string) (*model.Word, error)
		GetWords([]string) ([]*model.Word, error)
	}
}

func NewWordHandler(svc interface {
	GetWord(string) (*model.Word, error)
	GetWords([]string) ([]*model.Word, error)
}) *WordHandler {
	return &WordHandler{service: svc}
}

func (h *WordHandler) GetWord(c *gin.Context) {

	word := c.Param("word")
	result, err := h.service.GetWord(word)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

type WordsRequest struct {
	Words []string `json:"words"`
}

func (h *WordHandler) GetWords(c *gin.Context) {

	var req WordsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	results, err := h.service.GetWords(req.Words)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, results)
}
