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

	if r.Method == "POST" {
		w.Header().Set("Content-Type", "application/json")

		reqBody, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error: %s", err.Error()), http.StatusInternalServerError)
			return
		}

		var requestBody entity.ExecuteRequest
		json.Unmarshal(reqBody, &requestBody)

		if requestBody.Content == "" {
			http.Error(w, "Field 'content' was not provided", http.StatusBadRequest)
			return
		}
		if requestBody.LangSlug == "" {
			http.Error(w, "Field 'lang_slug' was not provided", http.StatusBadRequest)
			return
		}

		switch requestBody.LangSlug {
		case "cpp":
			{
				_, err := h.execCPP.ExecuteFromSource(requestBody.Content)
				if err != nil {
					http.Error(w, fmt.Sprintf("Error: %s", err.Error()), http.StatusInternalServerError)
					return
				}
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(map[string]string{"Result": ""})
				return
			}
		default:
			{
				http.Error(w, fmt.Sprintf("Language '%s' is not supported", requestBody.LangSlug), http.StatusBadRequest)
			}
		}
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
