package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hungp29/x-word/internal/model"
	"github.com/hungp29/x-word/internal/service"
)

type WordService interface {
	GetWord(word string, dict service.Dictionary) (*model.Word, error)
	GetWords(words []string, dict service.Dictionary) ([]*model.Word, error)
}

type WordHandler struct {
	service WordService
}

func NewWordHandler(svc WordService) *WordHandler {
	return &WordHandler{service: svc}
}

// parseDictionary reads the optional ?dict= query param and defaults to DictionaryEnglish.
// Returns 400 if the value is provided but unknown.
func parseDictionary(c *gin.Context) (service.Dictionary, bool) {
	raw := c.DefaultQuery("dict", string(service.DictionaryEnglish))
	dict := service.Dictionary(raw)
	if _, ok := service.DictionaryURLTemplates[dict]; !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "unknown dictionary: " + raw})
		return "", false
	}
	return dict, true
}

func (h *WordHandler) GetWord(c *gin.Context) {
	dict, ok := parseDictionary(c)
	if !ok {
		return
	}

	word := c.Param("word")
	result, err := h.service.GetWord(word, dict)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

type WordsRequest struct {
	Words []string           `json:"words"`
	Dict  service.Dictionary `json:"dict"`
}

func (h *WordHandler) GetWords(c *gin.Context) {
	var req WordsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	dict := req.Dict
	if dict == "" {
		dict = service.DictionaryEnglish
	}
	if _, ok := service.DictionaryURLTemplates[dict]; !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "unknown dictionary: " + string(dict)})
		return
	}

	results, err := h.service.GetWords(req.Words, dict)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, results)
}
