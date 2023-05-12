package handler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"server/internal/entity"
	"server/internal/executor"
)

type Handler struct {
	execCPP executor.Executor
}

func NewHandler(e executor.Executor) *Handler {
	return &Handler{
		execCPP: e,
	}
}

func (h *Handler) HandleSubmit(w http.ResponseWriter, r *http.Request) {
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request body"})
		return
	}

	var requestBody entity.ExecuteRequest
	json.Unmarshal(reqBody, &requestBody)

	if requestBody.Content == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "content was not provided"})
		return
	}
	if requestBody.LangSlug == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "lang_slug was not provided"})
		return
	}

	switch requestBody.LangSlug {
	case "cpp":
		{
			_, err := h.execCPP.ExecuteFromSource(requestBody.Content)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
				return
			}
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]string{"Result": ""})
			return
		}
	default:
		{
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": fmt.Sprintf("Language '%s' is not supported right now", requestBody.LangSlug)})

		}
	}
}
