package handlers

import (
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/moukoublen/goboilerplate/build"
	"github.com/moukoublen/goboilerplate/internal/httpx"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type AboutHandler struct {
	info Info
}

func NewAboutHandler() AboutHandler {
	return AboutHandler{info: Info{
		Version:     build.Version,
		Branch:      build.Branch,
		Commit:      build.Commit,
		CommitShort: build.CommitShort,
		Tag:         build.Tag,
	}}
}

func (a *AboutHandler) Register(router *chi.Mux) {
	router.Get("/about", a.About)
}

func (a AboutHandler) About(w http.ResponseWriter, r *http.Request) {
	accept := r.Header.Get(`Accept`)
	switch {
	case strings.Contains(`application/vnd.google.protobuf`, accept) || strings.Contains(`application/x-protobuf`, accept) || strings.Contains(`application/protobuf`, accept):
		msg := &About{
			Version:     a.info.Version,
			Branch:      a.info.Branch,
			Commit:      a.info.Commit,
			CommitShort: a.info.CommitShort,
			Tag:         a.info.Tag,
			Dt:          timestamppb.Now(),
		}
		b, err := proto.Marshal(msg)
		if err != nil {
			return
		}

		w.Header().Add("Content-Type", `application/x-protobuf`)

		w.WriteHeader(http.StatusOK)
		w.Write(b)
	default:
		httpx.RespondJSON(r.Context(), w, http.StatusOK, a.info)
	}
}

type Info struct {
	Version     string `json:"version"`
	Branch      string `json:"branch"`
	Commit      string `json:"commit"`
	CommitShort string `json:"commit_short"`
	Tag         string `json:"tag"`
}
