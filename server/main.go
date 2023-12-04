// main package
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"server/db"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	db, err := db.NewDB()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	slog.Info("start server")

	h := newHandler(db)
	mux := http.NewServeMux()
	mux.HandleFunc("/create", h.Create)
	mux.HandleFunc("/update", h.Update)
	mux.HandleFunc("/delete", h.Delete)
	mux.HandleFunc("/get", h.Get)
	mux.HandleFunc("/getAll", h.GetAll)

	sev := &http.Server{
		Addr:              ":8080",
		Handler:           mux,
		ReadTimeout:       30 * time.Second,
		ReadHeaderTimeout: 30 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       120 * time.Second,
		MaxHeaderBytes:    1 << 20,
	}

	slog.Info("listenAndServe is started")

	if err := sev.ListenAndServe(); err != nil {
		log.Print(err)
	}

	slog.Info("listenAndServe is completed")
}

/* -------------------------------------------------------------------------- */
/*                                   Handler                                  */
/* -------------------------------------------------------------------------- */

type handler struct {
	db   *db.DB
	data map[string]string
}

func newHandler(
	db *db.DB,
) *handler {
	return &handler{
		db:   db,
		data: map[string]string{},
	}
}

func (h *handler) Create(w http.ResponseWriter, r *http.Request) {
	if r.Method !=
		"POST" {
		if _, err := w.Write(methodNotAllowedResponse); err != nil {
			log.Print(err)
		}

		return
	}

	alb := &db.Album{
		Title:  "",
		Artist: "",
		Price:  0,
	}
	if err := json.NewDecoder(r.Body).Decode(alb); err != nil {
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

	id, err := h.db.Create(*alb)
	if err != nil {
		_, err = w.Write(
			createResponseInJSON(
				apiResponse{
					Code:    http.StatusInternalServerError,
					Message: fmt.Sprintf("create is error %v", err),
				},
			),
		)
		if err != nil {
			log.Print(err)
		}
	}

	_, err = w.Write(
		createResponseInJSON(
			apiResponse{
				Code:    http.StatusOK,
				Message: fmt.Sprintf("id %d created", id),
			},
		),
	)
	if err != nil {
		log.Print(err)
	}
}

func (h *handler) Update(w http.ResponseWriter, r *http.Request) {
	if r.Method !=
		"POST" {
		if _, err := w.Write(methodNotAllowedResponse); err != nil {
			log.Print(err)
		}

		return
	}

	alb := &db.Album{}
	if err := json.NewDecoder(r.Body).Decode(alb); err != nil {
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

	_, err := h.db.Update(*alb)
	if err != nil {
		_, err := w.Write(
			createResponseInJSON(
				apiResponse{
					Code:    http.StatusInternalServerError,
					Message: fmt.Sprintf("update is error %v", err),
				},
			),
		)
		if err != nil {
			log.Print(err)
		}

		return
	}

	if _, err := w.Write(
		createResponseInJSON(
			apiResponse{
				Code:    http.StatusOK,
				Message: fmt.Sprintf("%v updated", alb),
			},
		),
	); err != nil {
		log.Print(err)

		return
	}
}

func (h *handler) Delete(w http.ResponseWriter, r *http.Request) {
	if r.Method !=
		"Delete" {
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
	slog.Info("handler is %v", h)

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

	id, err := strconv.Atoi(paramKey)
	if err != nil {
		log.Print(err)

		return
	}

	slog.Info("id is %v", id)
	album, err := h.db.Read(id)
	if err != nil {
		w.Write(
			createResponseInJSON(
				apiResponse{
					Code:    http.StatusInternalServerError,
					Message: fmt.Sprintf("get is error %v", err),
				},
			),
		)

		return
	}

	_, err = w.Write(
		createResponseInJSON(
			apiResponse{
				Code:    http.StatusOK,
				Message: fmt.Sprintf("album is %v", album),
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

func (h *handler) GetAll(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		if _, err := w.Write(methodNotAllowedResponse); err != nil {
			log.Print(err)
		}

		return
	}

	albums, err := h.db.ReadAll()
	if err != nil {
		w.Write(
			createResponseInJSON(
				apiResponse{
					Code:    http.StatusInternalServerError,
					Message: fmt.Sprintf("get all is error %v", err),
				},
			),
		)

		return
	}

	albumsJSON, err := json.Marshal(albums)
	if err != nil {
		w.Write(
			createResponseInJSON(
				apiResponse{
					Code:    http.StatusInternalServerError,
					Message: fmt.Sprintf("get all is error %v", err),
				},
			),
		)

		return
	}

	_, err = w.Write(
		createResponseInJSON(
			apiResponse{
				Code:    http.StatusOK,
				Message: fmt.Sprintf("album is %s", string(albumsJSON)),
			},
		),
	)
	if err != nil {
		log.Print(err)
	}
}

/* -------------------------------------------------------------------------- */
/*                                    model                                   */
/* -------------------------------------------------------------------------- */

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
