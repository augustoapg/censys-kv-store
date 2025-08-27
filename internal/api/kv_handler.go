package api

import (
	"net/http"

	"github.com/augustoapg/censysKvStore/internal/utils"
	"github.com/go-chi/chi/v5"
)

type KVHandler struct{}

func NewKVHandler() *KVHandler {
	return &KVHandler{}
}

func (h *KVHandler) HandleGetKvByKey(w http.ResponseWriter, r *http.Request) {
	key := chi.URLParam(r, "key")

	if key == "" {
		utils.SendErrorResponse(w, http.StatusBadRequest, "key is required")
		return
	}

	value := "mock_" + key

	data := map[string]string{
		"key":   key,
		"value": value,
	}

	utils.SendJSONResponse(w, http.StatusOK, data)
}

func (h *KVHandler) HandleSetKv(w http.ResponseWriter, r *http.Request) {
	utils.SendJSONResponse(w, http.StatusCreated, map[string]string{
		"message": "kv set successfully",
	})
}
