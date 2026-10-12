package httpapi

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"todogo/internal/storage/memory"
	"todogo/internal/todo"
)

var _ TodoStore = (*memory.Store)(nil)

func newTestRouter() http.Handler {
	return NewRouter(NewHandler(memory.New()))
}

func TestCreateAndGet(t *testing.T) {
	router := newTestRouter()

	req := httptest.NewRequest(http.MethodPost, "/todos",
		strings.NewReader(`{"title":"купить молоко"}`))
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("status = %d, want 201, body: %s", rec.Code, rec.Body.String())
	}

	var created todo.Todo
	if err := json.NewDecoder(rec.Body).Decode(&created); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if created.ID != 1 || created.Title != "купить молоко" {
		t.Errorf("unexpected todo: %+v", created)
	}

	req = httptest.NewRequest(http.MethodGet, "/todos/1", nil)
	rec = httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}
}

func TestErrorStatuses(t *testing.T) {
	tests := []struct {
		name       string
		method     string
		path       string
		body       string
		wantStatus int
	}{
		{"пустой title", http.MethodPost, "/todos", `{"title":""}`, http.StatusBadRequest},
		{"кривой json", http.MethodPost, "/todos", `мусор`, http.StatusBadRequest},
		{"задачи нет", http.MethodGet, "/todos/999", "", http.StatusNotFound},
		{"кривой id", http.MethodGet, "/todos/abc", "", http.StatusBadRequest},
		{"удаление несуществующей", http.MethodDelete, "/todos/999", "", http.StatusNotFound},
		{"неизвестный метод", http.MethodPatch, "/todos", "", http.StatusMethodNotAllowed},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := newTestRouter()
			req := httptest.NewRequest(tt.method, tt.path, strings.NewReader(tt.body))
			rec := httptest.NewRecorder()

			router.ServeHTTP(rec, req)

			if rec.Code != tt.wantStatus {
				t.Errorf("status = %d, want %d, body: %s",
					rec.Code, tt.wantStatus, rec.Body.String())
			}
		})
	}
}
