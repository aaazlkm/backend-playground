// main package
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

func main() {
	h := newHandler()
	mux := http.NewServeMux()
	mux.HandleFunc("/create", h.Create)
	mux.HandleFunc("/update", h.Update)
	mux.HandleFunc("/delete", h.Delete)
	mux.HandleFunc("/get", h.Get)

	sev := &http.Server{
		Addr:                         ":8080",
		Handler:                      mux,
		ReadTimeout:                  30 * time.Second,  //nolint:gomnd
		ReadHeaderTimeout:            30 * time.Second,  //nolint:gomnd
		WriteTimeout:                 30 * time.Second,  //nolint:gomnd
		IdleTimeout:                  120 * time.Second, //nolint:gomnd
		MaxHeaderBytes:               1 << 20,           //nolint:gomnd
		TLSConfig:                    nil,
		TLSNextProto:                 nil,
		ConnState:                    nil,
		ErrorLog:                     nil,
		BaseContext:                  nil,
		ConnContext:                  nil,
		DisableGeneralOptionsHandler: false,
	}

	if err := sev.ListenAndServe(); err != nil {
		log.Print(err)
	}
}

/* -------------------------------------------------------------------------- */
/*                                   Handler                                  */
/* -------------------------------------------------------------------------- */

type handler struct {
	data map[string]string
}

func newHandler() *handler {
	return &handler{
		data: map[string]string{},
	}
}

func (h *handler) Create(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		if _, err := w.Write(methodNotAllowedResponse); err != nil {
			log.Print(err)
		}

		return
	}

	req := &apiRequest{
		Key:   "",
		Value: "",
	}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		if _, err := w.Write(
			createResponseInJSON(
				apiResponse{
					Code:    http.StatusBadRequest,
					Message: "request body is invalid",
				},
			),
		); err != nil {
			log.Print(err)
		}

		return
	}

	h.data[req.Key] = req.Value

	_, err := w.Write(
		createResponseInJSON(
			apiResponse{
				Code:    http.StatusOK,
				Message: fmt.Sprintf("Key %s created", req.Key),
			},
		),
	)
	if err != nil {
		log.Print(err)
	}
}

func (h *handler) Update(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		if _, err := w.Write(methodNotAllowedResponse); err != nil {
			log.Print(err)
		}

		return
	}

	req := &apiRequest{
		Key:   "",
		Value: "",
	}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		_, err := w.Write(
			createResponseInJSON(
				apiResponse{
					Code:    http.StatusBadRequest,
					Message: "request body is invalid",
				},
			),
		)
		if err != nil {
			log.Print(err)
		}

		return
	}

	if _, ok := h.data[req.Key]; !ok {
		if _, err := w.Write(methodNotAllowedResponse); err != nil {
			log.Print(err)
		}

		return
	}

	h.data[req.Key] = req.Value

	if _, err := w.Write(
		createResponseInJSON(
			apiResponse{
				Code:    http.StatusOK,
				Message: fmt.Sprintf("Key %s created", req.Key),
			},
		),
	); err != nil {
		log.Print(err)
	}
}

func (h *handler) Delete(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		if _, err := w.Write(methodNotAllowedResponse); err != nil {
			log.Println(err)
		}

		return
	}

	paramKey := r.URL.Query().Get("key")
	if len(paramKey) == 0 {
		if _, err := w.Write(notFoundResponse); err != nil {
			log.Print(err)
		}

		return
	}

	h.data[paramKey] = ""

	if _, err := w.Write(
		createResponseInJSON(
			apiResponse{
				Code:    http.StatusOK,
				Message: fmt.Sprintf("Key %s deleted", paramKey),
			},
		),
	); err != nil {
		log.Print(err)
	}
}

func (h *handler) Get(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		if _, err := w.Write(methodNotAllowedResponse); err != nil {
			log.Print(err)
		}

		return
	}

	paramKey := r.URL.Query().Get("key")
	if len(paramKey) == 0 {
		if _, err := w.Write(notFoundResponse); err != nil {
			log.Print(err)
		}

		return
	}

	_, err := w.Write(
		createResponseInJSON(
			apiResponse{
				Code:    http.StatusOK,
				Message: fmt.Sprintf("Key %s value is %s", paramKey, h.data[paramKey]),
			},
		),
	)
	if err != nil {
		log.Print(err)
	}
}

/* -------------------------------------------------------------------------- */
/*                               error response                               */
/* -------------------------------------------------------------------------- */

// 共通のレスポンスを定義するためにグローバル変数を使用
//
//nolint:gochecknoglobals
var (
	methodNotAllowedResponse = createResponseInJSON(
		apiResponse{
			Code:    http.StatusMethodNotAllowed,
			Message: "Method not supported",
		},
	)
	notFoundResponse = createResponseInJSON(
		apiResponse{
			Code:    http.StatusNotFound,
			Message: "Key not found",
		},
	)
)

/* -------------------------------------------------------------------------- */
/*                                    model                                   */
/* -------------------------------------------------------------------------- */

type apiRequest struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type apiResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

/* -------------------------------------------------------------------------- */
/*                                    util                                    */
/* -------------------------------------------------------------------------- */

func createResponseInJSON(apiResponse apiResponse) []byte {
	response, err := json.Marshal(apiResponse)
	if err != nil {
		// 通常エラーにならないはずなので、ログ出力だけしておく
		log.Print(err)
	}

	return response
}
