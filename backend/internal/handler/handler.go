package handler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"server/internal/entity"
	"server/internal/executor"
	"server/pkg/logger"
)

type Handler struct {
	execCPP executor.Executor
	lgr     *logger.Logger
}

func NewHandler(e executor.Executor, lgr *logger.Logger) *Handler {
	return &Handler{
		execCPP: e,
		lgr:     lgr,
	}
}

func (h *Handler) HandleSubmit(w http.ResponseWriter, r *http.Request) {

	if r.Method == "POST" {
		w.Header().Set("Content-Type", "application/json")

		reqBody, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error: %s", err.Error()), http.StatusInternalServerError)
			h.lgr.SubmitError(err.Error())
			return
		}

		var requestBody entity.ExecuteRequest
		json.Unmarshal(reqBody, &requestBody)

		if requestBody.Content == "" {
			http.Error(w, "Field 'content' was not provided", http.StatusBadRequest)
			h.lgr.SubmitError("Field 'content' was not provided")
			return
		}
		if requestBody.LangSlug == "" {
			http.Error(w, "Field 'lang_slug' was not provided", http.StatusBadRequest)
			h.lgr.SubmitError("Field 'lang_slug' was not provided")
			return
		}

		switch requestBody.LangSlug {
		case "cpp":
			{
				result, err, errcode := h.execCPP.ExecuteFromSource(requestBody.Content)
				if errcode != entity.PROCESS_OK {
					switch errcode {
					case entity.PROCESS_COMPILE_ERROR:
						{
							http.Error(w, fmt.Sprintf("Compile error: %s", err.Error()), http.StatusUnprocessableEntity)
							h.lgr.SubmitError(fmt.Sprintf("Compile error: %s\n", err.Error()))
							return
						}
					case entity.PROCESS_RUNTIME_ERROR:
						{
							http.Error(w, fmt.Sprintf("Runtime error: %s", err.Error()), http.StatusUnprocessableEntity)
							h.lgr.SubmitError(fmt.Sprintf("Runtime error: %s\n", err.Error()))
							return
						}
					default:
						{
							http.Error(w, "Server error", http.StatusInternalServerError)
							h.lgr.SubmitError(fmt.Sprintf("Server error: %s\n", err.Error()))
						}
					}
				}
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(map[string]string{"Result": result})
				h.lgr.SubmitSuccess()
				return
			}
		default:
			{
				http.Error(w, fmt.Sprintf("Language '%s' is not supported", requestBody.LangSlug), http.StatusBadRequest)
				h.lgr.SubmitError(fmt.Sprintf("Language '%s' is not supported", requestBody.LangSlug))
			}
		}
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		h.lgr.SubmitError("Method is not allowed")
	}
}
