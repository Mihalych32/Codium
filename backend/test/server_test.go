package test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"server/internal/entity"
	"server/internal/executor"
	"server/internal/handler"
	"strings"
	"testing"
)

func TestHandleSubmit(t *testing.T) {
	tt := []struct {
		name       string
		method     string
		input      *entity.ExecuteRequest
		want       string
		statusCode int
	}{
		{
			name:       "Empty request",
			method:     http.MethodPost,
			input:      &entity.ExecuteRequest{},
			want:       "Field 'content' was not provided",
			statusCode: http.StatusBadRequest,
		},
		{
			name:       "Only content provided",
			method:     http.MethodPost,
			input:      &entity.ExecuteRequest{Content: "#include<iostream>\n\nint main() {\n\treturn0;\n}"},
			want:       "Field 'lang_slug' was not provided",
			statusCode: http.StatusBadRequest,
		},
		{
			name:       "Only lang_slug provided",
			method:     http.MethodPost,
			input:      &entity.ExecuteRequest{LangSlug: "cpp"},
			want:       "Field 'content' was not provided",
			statusCode: http.StatusBadRequest,
		},
		{
			name:       "Unsupported lang_slug provided",
			method:     http.MethodPost,
			input:      &entity.ExecuteRequest{Content: "#include<iostream>\n\nint main() {\n\treturn0;\n}", LangSlug: "rs"},
			want:       "Language 'rs' is not supported",
			statusCode: http.StatusBadRequest,
		},
		{
			name:       "Unsupported method",
			method:     http.MethodGet,
			input:      &entity.ExecuteRequest{},
			want:       "Method not allowed",
			statusCode: http.StatusMethodNotAllowed,
		},
		{
			name:       "Normal request",
			method:     http.MethodPost,
			input:      &entity.ExecuteRequest{Content: "#include<iostream>\n\nint main() {\n\treturn0;\n}", LangSlug: "cpp"},
			want:       `{"Result":""}`,
			statusCode: http.StatusOK,
		},
	}

	execcpp := executor.NewExecutorCPP()

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {

			requestByte, _ := json.Marshal(tc.input)
			body := bytes.NewReader(requestByte)

			request := httptest.NewRequest(tc.method, "http://localhost:8080/api/submit/", body)
			responseRecorder := httptest.NewRecorder()

			h := handler.NewHandler(execcpp)
			h.HandleSubmit(responseRecorder, request)

			if responseRecorder.Code != tc.statusCode {
				t.Errorf("Want status %d, got %d", tc.statusCode, responseRecorder.Code)
			}
			if strings.TrimSpace(responseRecorder.Body.String()) != tc.want {
				t.Errorf("Want response '%s', got '%s'", tc.want, responseRecorder.Body)
			}
		})
	}
}
