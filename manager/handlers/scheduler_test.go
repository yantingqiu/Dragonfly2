/*
 *     Copyright 2024 The Dragonfly Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"d7y.io/dragonfly/v2/manager/models"
	"d7y.io/dragonfly/v2/manager/service/mocks"
	"d7y.io/dragonfly/v2/manager/types"
)

var (
	mockSchedulerReqBody = `
		{
		   "host_name": "foo",
		   "ip": "127.0.0.1",
		   "port": 8003,
		   "scheduler_cluster_id": 2
		}`
	mockCreateSchedulerRequest = types.CreateSchedulerRequest{
		Hostname:           "foo",
		IP:                 "127.0.0.1",
		Port:               8003,
		SchedulerClusterID: 2,
	}
	mockUpdateSchedulerRequest = types.UpdateSchedulerRequest{
		IP:                 "127.0.0.1",
		Port:               8003,
		SchedulerClusterID: 2,
	}
	mockSchedulerModel = &models.Scheduler{
		BaseModel: mockBaseModel,
		Hostname:  "foo",
		IP:        "127.0.0.1",
		Port:      8003,
	}
)

func mockSchedulerRouter(h *Handlers) *gin.Engine {
	r := gin.Default()
	apiv1 := r.Group("/api/v1")
	s := apiv1.Group("/schedulers")
	s.POST("", h.CreateScheduler)
	s.DELETE(":id", h.DestroyScheduler)
	s.PATCH(":id", h.UpdateScheduler)
	s.GET(":id", h.GetScheduler)
	s.GET("", h.GetSchedulers)
	return r
}

func TestHandlers_CreateScheduler(t *testing.T) {
	tests := []struct {
		name   string
		req    *http.Request
		mock   func(ms *mocks.MockServiceMockRecorder)
		expect func(t *testing.T, w *httptest.ResponseRecorder)
	}{
		{
			name: "unprocessable entity",
			req:  httptest.NewRequest(http.MethodPost, "/api/v1/schedulers", nil),
			mock: func(ms *mocks.MockServiceMockRecorder) {},
			expect: func(t *testing.T, w *httptest.ResponseRecorder) {
				assert := assert.New(t)
				assert.Equal(http.StatusUnprocessableEntity, w.Code)
			},
		},
		{
			name: "success",
			req:  httptest.NewRequest(http.MethodPost, "/api/v1/schedulers", strings.NewReader(mockSchedulerReqBody)),
			mock: func(ms *mocks.MockServiceMockRecorder) {
				ms.CreateScheduler(gomock.Any(), gomock.Eq(mockCreateSchedulerRequest)).Return(mockSchedulerModel, nil).Times(1)
			},
			expect: func(t *testing.T, w *httptest.ResponseRecorder) {
				assert := assert.New(t)
				assert.Equal(http.StatusOK, w.Code)
				scheduler := models.Scheduler{}
				err := json.Unmarshal(w.Body.Bytes(), &scheduler)
				assert.NoError(err)
				assert.Equal(mockSchedulerModel, &scheduler)
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctl := gomock.NewController(t)
			defer ctl.Finish()
			svc := mocks.NewMockService(ctl)
			w := httptest.NewRecorder()
			h := New(svc)
			mockRouter := mockSchedulerRouter(h)

			tc.mock(svc.EXPECT())
			mockRouter.ServeHTTP(w, tc.req)
			tc.expect(t, w)
		})
	}
}

func TestHandlers_DestroyScheduler(t *testing.T) {
	tests := []struct {
		name   string
		req    *http.Request
		mock   func(ms *mocks.MockServiceMockRecorder)
		expect func(t *testing.T, w *httptest.ResponseRecorder)
	}{
		{
			name: "unprocessable entity",
			req:  httptest.NewRequest(http.MethodDelete, "/api/v1/schedulers/test", nil),
			mock: func(ms *mocks.MockServiceMockRecorder) {},
			expect: func(t *testing.T, w *httptest.ResponseRecorder) {
				assert := assert.New(t)
				assert.Equal(http.StatusUnprocessableEntity, w.Code)
			},
		},
		{
			name: "success",
			req:  httptest.NewRequest(http.MethodDelete, "/api/v1/schedulers/2", nil),
			mock: func(ms *mocks.MockServiceMockRecorder) {
				ms.DestroyScheduler(gomock.Any(), gomock.Eq(uint(2))).Return(nil).Times(1)
			},
			expect: func(t *testing.T, w *httptest.ResponseRecorder) {
				assert := assert.New(t)
				assert.Equal(http.StatusOK, w.Code)
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctl := gomock.NewController(t)
			defer ctl.Finish()
			svc := mocks.NewMockService(ctl)
			w := httptest.NewRecorder()
			h := New(svc)
			mockRouter := mockSchedulerRouter(h)

			tc.mock(svc.EXPECT())
			mockRouter.ServeHTTP(w, tc.req)
			tc.expect(t, w)
		})
	}
}

func TestHandlers_UpdateScheduler(t *testing.T) {
	tests := []struct {
		name   string
		req    *http.Request
		mock   func(ms *mocks.MockServiceMockRecorder)
		expect func(t *testing.T, w *httptest.ResponseRecorder)
	}{
		{
			name: "unprocessable entity caused by uri",
			req:  httptest.NewRequest(http.MethodPatch, "/api/v1/schedulers/test", nil),
			mock: func(ms *mocks.MockServiceMockRecorder) {},
			expect: func(t *testing.T, w *httptest.ResponseRecorder) {
				assert := assert.New(t)
				assert.Equal(http.StatusUnprocessableEntity, w.Code)
			},
		},
		{
			name: "unprocessable entity caused by body",
			req:  httptest.NewRequest(http.MethodPatch, "/api/v1/schedulers/2", nil),
			mock: func(ms *mocks.MockServiceMockRecorder) {},
			expect: func(t *testing.T, w *httptest.ResponseRecorder) {
				assert := assert.New(t)
				assert.Equal(http.StatusUnprocessableEntity, w.Code)
			},
		},
		{
			name: "success",
			req:  httptest.NewRequest(http.MethodPatch, "/api/v1/schedulers/2", strings.NewReader(mockSchedulerReqBody)),
			mock: func(ms *mocks.MockServiceMockRecorder) {
				ms.UpdateScheduler(gomock.Any(), gomock.Eq(uint(2)), gomock.Eq(mockUpdateSchedulerRequest)).Return(mockSchedulerModel, nil).Times(1)
			},
			expect: func(t *testing.T, w *httptest.ResponseRecorder) {
				assert := assert.New(t)
				assert.Equal(http.StatusOK, w.Code)
				scheduler := models.Scheduler{}
				err := json.Unmarshal(w.Body.Bytes(), &scheduler)
				assert.NoError(err)
				assert.Equal(mockSchedulerModel, &scheduler)
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctl := gomock.NewController(t)
			defer ctl.Finish()
			svc := mocks.NewMockService(ctl)
			w := httptest.NewRecorder()
			h := New(svc)
			mockRouter := mockSchedulerRouter(h)

			tc.mock(svc.EXPECT())
			mockRouter.ServeHTTP(w, tc.req)
			tc.expect(t, w)
		})
	}
}

func TestHandlers_GetScheduler(t *testing.T) {
	tests := []struct {
		name   string
		req    *http.Request
		mock   func(ms *mocks.MockServiceMockRecorder)
		expect func(t *testing.T, w *httptest.ResponseRecorder)
	}{
		{
			name: "unprocessable entity",
			req:  httptest.NewRequest(http.MethodGet, "/api/v1/schedulers/test", nil),
			mock: func(ms *mocks.MockServiceMockRecorder) {},
			expect: func(t *testing.T, w *httptest.ResponseRecorder) {
				assert := assert.New(t)
				assert.Equal(http.StatusUnprocessableEntity, w.Code)
			},
		},
		{
			name: "success",
			req:  httptest.NewRequest(http.MethodGet, "/api/v1/schedulers/2", nil),
			mock: func(ms *mocks.MockServiceMockRecorder) {
				ms.GetScheduler(gomock.Any(), gomock.Eq(uint(2))).Return(mockSchedulerModel, nil).Times(1)
			},
			expect: func(t *testing.T, w *httptest.ResponseRecorder) {
				assert := assert.New(t)
				assert.Equal(http.StatusOK, w.Code)
				scheduler := models.Scheduler{}
				err := json.Unmarshal(w.Body.Bytes(), &scheduler)
				assert.NoError(err)
				assert.Equal(mockSchedulerModel, &scheduler)
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctl := gomock.NewController(t)
			defer ctl.Finish()
			svc := mocks.NewMockService(ctl)
			w := httptest.NewRecorder()
			h := New(svc)
			mockRouter := mockSchedulerRouter(h)

			tc.mock(svc.EXPECT())
			mockRouter.ServeHTTP(w, tc.req)
			tc.expect(t, w)
		})
	}
}

func TestHandlers_GetSchedulers(t *testing.T) {
	tests := []struct {
		name   string
		req    *http.Request
		mock   func(ms *mocks.MockServiceMockRecorder)
		expect func(t *testing.T, w *httptest.ResponseRecorder)
	}{
		{
			name: "unprocessable entity",
			req:  httptest.NewRequest(http.MethodGet, "/api/v1/schedulers?page=-1", nil),
			mock: func(ms *mocks.MockServiceMockRecorder) {},
			expect: func(t *testing.T, w *httptest.ResponseRecorder) {
				assert := assert.New(t)
				assert.Equal(http.StatusUnprocessableEntity, w.Code)
			},
		},
		{
			name: "success",
			req:  httptest.NewRequest(http.MethodGet, "/api/v1/schedulers", nil),
			mock: func(ms *mocks.MockServiceMockRecorder) {
				ms.GetSchedulers(gomock.Any(), gomock.Eq(types.GetSchedulersQuery{
					Page:    1,
					PerPage: 10,
				})).Return([]models.Scheduler{*mockSchedulerModel}, int64(1), nil).Times(1)
			},
			expect: func(t *testing.T, w *httptest.ResponseRecorder) {
				assert := assert.New(t)
				assert.Equal(http.StatusOK, w.Code)
				scheduler := models.Scheduler{}
				err := json.Unmarshal(w.Body.Bytes()[1:w.Body.Len()-1], &scheduler)
				assert.NoError(err)
				assert.Equal(mockSchedulerModel, &scheduler)
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctl := gomock.NewController(t)
			defer ctl.Finish()
			svc := mocks.NewMockService(ctl)
			w := httptest.NewRecorder()
			h := New(svc)
			mockRouter := mockSchedulerRouter(h)

			tc.mock(svc.EXPECT())
			mockRouter.ServeHTTP(w, tc.req)
			tc.expect(t, w)
		})
	}
}
