package test

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"server/internal/entity"
	"server/internal/executor"
	"server/internal/handler"
	"server/pkg/logger"
	"strings"
	"testing"

	"github.com/joho/godotenv"
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
			name:       "Field 'lang_slug' not provided",
			method:     http.MethodPost,
			input:      &entity.ExecuteRequest{Content: "#include <iostream>\n\nint main() {\n\treturn 0;\n}"},
			want:       "Field 'lang_slug' was not provided",
			statusCode: http.StatusBadRequest,
		},
		{
			name:       "Field 'content' not provided",
			method:     http.MethodPost,
			input:      &entity.ExecuteRequest{LangSlug: "cpp"},
			want:       "Field 'content' was not provided",
			statusCode: http.StatusBadRequest,
		},
		{
			name:       "Unsupported lang_slug provided",
			method:     http.MethodPost,
			input:      &entity.ExecuteRequest{Content: "#include <iostream>\n\nint main() {\n\treturn 0;\n}", LangSlug: "rs"},
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
			name:       "Uncompilable code",
			method:     http.MethodPost,
			input:      &entity.ExecuteRequest{Content: "#include <iostream>\nint main(){", LangSlug: "cpp"},
			want:       ``,
			statusCode: http.StatusUnprocessableEntity,
		},
		{
			name:       "Hello world program",
			method:     http.MethodPost,
			input:      &entity.ExecuteRequest{Content: "#include <iostream>\nint main(){std::cout<<\"Hello world!\"<<'\\n';return 0;}", LangSlug: "cpp"},
			want:       `{"Result":"Hello world!\n"}`,
			statusCode: http.StatusOK,
		},
		{
			name:       "Print numbers using a for loop",
			method:     http.MethodPost,
			input:      &entity.ExecuteRequest{Content: "#include <iostream>\nint main(){for(int i=0;i<5;i++)std::cout<<i<<'\\n';return 0;}", LangSlug: "cpp"},
			want:       `{"Result":"0\n1\n2\n3\n4\n"}`,
			statusCode: http.StatusOK,
		},
		{
			name:       "Call a function",
			method:     http.MethodPost,
			input:      &entity.ExecuteRequest{Content: "#include <iostream>\nusing namespace std;\n\nint num_sum(int a, int b) {\n\treturn a + b;\n}\n\nint main() {\n\tint res = num_sum(2, 3);\n\tstd::cout << res << '\\n';\n\treturn 0;\n}", LangSlug: "cpp"},
			want:       `{"Result":"5\n"}`,
			statusCode: http.StatusOK,
		},
	}

	if err := godotenv.Load("../.env"); err != nil {
		log.Fatal("Could not load the .env file")
	}

	execcpp := executor.NewExecutorCPP()
	lgr := logger.NewLogger()

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {

			requestByte, _ := json.Marshal(tc.input)
			body := bytes.NewReader(requestByte)

			request := httptest.NewRequest(tc.method, "http://localhost:8080/api/submit/", body)
			responseRecorder := httptest.NewRecorder()

			h := handler.NewHandler(execcpp, lgr)
			h.HandleSubmit(responseRecorder, request)

			if responseRecorder.Code != tc.statusCode {
				t.Errorf("Want status %d, got %d", tc.statusCode, responseRecorder.Code)
			}

			if tc.statusCode != http.StatusUnprocessableEntity {
				if strings.TrimSpace(responseRecorder.Body.String()) != tc.want {
					t.Errorf("Want response '%s', got '%s'", tc.want, responseRecorder.Body)
				}
			}
		})
	}
}
