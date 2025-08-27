package api

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"

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

// HandleGetKvByKey retrieves a key-value pair from the store.
// If key not found or if soft deleted, return error
// Returns 200 with key-value
func (h *KVHandler) HandleGetKvByKey(w http.ResponseWriter, r *http.Request) {
	key := strings.TrimSpace(chi.URLParam(r, "key"))

	if key == "" {
		utils.SendErrorResponse(w, http.StatusBadRequest, "key is required")
		return
	}

	kv, err := h.KVStore.GetKvByKey(key)

	if err != nil {
		h.Logger.Printf("[HandleGetKvByKey] error getting kv by key: %v", err)

		if errors.Is(err, store.ErrKeyNotFound) {
			utils.SendErrorResponse(w, http.StatusNotFound, "key not found")
			return
		}

		utils.SendErrorResponse(w, http.StatusInternalServerError, "internal server error")
		return
	}

	utils.SendJSONResponse(w, http.StatusOK, map[string]any{
		"kv": kv,
	})
}

// HandleUpsertKv creates or updates a key-value pair in the store.
// If the key is soft deleted, it will be restored.
// Returns 200 with the created or updated key-value pair.
func (h *KVHandler) HandleUpsertKv(w http.ResponseWriter, r *http.Request) {
	var kv store.KV

	err := json.NewDecoder(r.Body).Decode(&kv)
	if err != nil {
		h.Logger.Printf("[HandleUpsertKv] error decoding request body: %v", err)
		utils.SendErrorResponse(w, http.StatusBadRequest, "invalid request body")
		return
	}

	// key must be present and trimmed before saving it
	kv.Key = strings.TrimSpace(kv.Key)
	if kv.Key == "" {
		h.Logger.Printf("[HandleUpsertKv] missing key")
		utils.SendErrorResponse(w, http.StatusBadRequest, "key must be a valid string")
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

// HandleDeleteKvByKey deletes key-value pair based on key
// If key not found or already soft deleted, return error
// Return 204
func (h *KVHandler) HandleDeleteKvByKey(w http.ResponseWriter, r *http.Request) {
	key := strings.TrimSpace(chi.URLParam(r, "key"))

	if key == "" {
		utils.SendErrorResponse(w, http.StatusBadRequest, "key is required")
		return
	}

	err := h.KVStore.DeleteKvByKey(key)

	if err != nil {
		h.Logger.Printf("[HandleDeleteKvByKey] error deleting kv: %v", err)
		if errors.Is(err, store.ErrKeyNotFound) {
			utils.SendErrorResponse(w, http.StatusNotFound, "key not found")
			return
		}

		utils.SendErrorResponse(w, http.StatusInternalServerError, "internal server error")
		return
	}

	utils.SendNoContentResponse(w)
}
