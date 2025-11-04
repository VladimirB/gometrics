package handler_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/VladimirB/gometrics/internal/handler"
	"github.com/stretchr/testify/assert"
)

func TestUpdateMetricHandler(t *testing.T) {
	type want struct {
		statusCode int
	}

	tests := []struct {
		name string
		request string
		want want
	}{
		{
			name: "incorrect path",
			request: "/update/name/10.45",
			want: want{
				statusCode: http.StatusNotFound,
			},
		},
		{
			name: "success counter metric",
			request: "/update/counter/name/100",
			want: want{
				statusCode: http.StatusOK,
			},
		},
		{
			name: "success gaige metric",
			request: "/update/gauge/name/100.33",
			want: want{
				statusCode: http.StatusOK,
			},
		},
		{
			name: "bad request on incorrect gauge metric value",
			request: "/update/gauge/name/str",
			want: want{
				statusCode: http.StatusBadRequest,
			},
		},
		{
			name: "bad request on incorrect counter metric value",
			request: "/update/counter/name/str",
			want: want{
				statusCode: http.StatusBadRequest,
			},
		},
		{
			name: "bad request on unknown metric type",
			request: "/update/unknown/name/100.33",
			want: want{
				statusCode: http.StatusBadRequest,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodPost, tt.request, nil)
			recorder := httptest.NewRecorder()
			handler := handler.UpdateMetricHandler()
			handler(recorder, request)
			result := recorder.Result()
			defer result.Body.Close()

			assert.Equal(t, tt.want.statusCode, result.StatusCode)
		})
	}
}
