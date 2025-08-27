package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/augustoapg/censysKvStore/internal/store"
	"github.com/augustoapg/censysKvStore/internal/utils"
	"github.com/go-chi/chi/v5"
)

type KVHandler struct {
	KVStore store.KVStore
	Logger  *log.Logger
}

func NewKVHandler(kvStore store.KVStore, logger *log.Logger) *KVHandler {
	return &KVHandler{
		KVStore: kvStore,
		Logger:  logger,
	}
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

// HandleUpsertKv creates or updates a key-value pair in the store.
// If the key is soft deleted, it will be restored.
// Returns the created or updated key-value pair.
func (h *KVHandler) HandleUpsertKv(w http.ResponseWriter, r *http.Request) {
	var kv store.KV

	err := json.NewDecoder(r.Body).Decode(&kv)
	if err != nil {
		h.Logger.Printf("[HandleUpsertKv] error decoding request body: %v", err)
		utils.SendErrorResponse(w, http.StatusBadRequest, "invalid request body")
		return
	}

	kv, err = h.KVStore.UpsertKv(&kv)
	if err != nil {
		h.Logger.Printf("[HandleUpsertKv] error upserting kv: %v", err)
		utils.SendErrorResponse(w, http.StatusInternalServerError, "internal server error")
		return
	}

	utils.SendJSONResponse(w, http.StatusOK, map[string]any{
		"kv": kv,
	})
}

func (h *KVHandler) HandleDeleteKvByKey(w http.ResponseWriter, r *http.Request) {
	key := chi.URLParam(r, "key")

	if key == "" {
		utils.SendErrorResponse(w, http.StatusBadRequest, "key is required")
		return
	}

	utils.SendNoContentResponse(w)
}
