package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
)

func (s *server) uploadCatalogPhotos() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		r.ParseMultipartForm(10 << 20)

		file, handler, err := r.FormFile("photo")

		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		defer file.Close()

		dst, err := os.Create("./static/catalog/" + handler.Filename)
		defer dst.Close()

		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		if _, err := io.Copy(dst, file); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		filePath := path.Join("./static/catalog/" + handler.Filename)

		s.respond(w, r, http.StatusCreated, map[string]interface{}{"status": "successfully create", "filePath": filePath})

	}
}

func (s *server) uploadGalleryPhotos() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		r.ParseMultipartForm(10 << 20)

		file, handler, err := r.FormFile("photo")

		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		defer file.Close()

		dst, err := os.Create("./static/gallery/" + handler.Filename)
		defer dst.Close()

		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		if _, err := io.Copy(dst, file); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		filePath := path.Join("./static/gallery/" + handler.Filename)

		s.respond(w, r, http.StatusCreated, map[string]interface{}{"status": "successfully create", "filePath": filePath})

	}
}

func (s *server) uploadProfilePhotos() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		r.ParseMultipartForm(10 << 20)

		file, handler, err := r.FormFile("photo")

		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		defer file.Close()

		dst, err := os.Create("./static/profile/" + handler.Filename)
		defer dst.Close()

		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		if _, err := io.Copy(dst, file); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		filePath := path.Join("./static/profile/" + handler.Filename)

		s.respond(w, r, http.StatusCreated, map[string]interface{}{"status": "successfully create", "filePath": filePath})

	}
}

func (s *server) uploadReviewsPhotos() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		r.ParseMultipartForm(10 << 20)

		file, handler, err := r.FormFile("photo")

		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		defer file.Close()

		dst, err := os.Create("./static/reviews/" + handler.Filename)
		defer dst.Close()

		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		if _, err := io.Copy(dst, file); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		filePath := path.Join("./static/reviews/" + handler.Filename)

		s.respond(w, r, http.StatusCreated, map[string]interface{}{"status": "successfully create", "filePath": filePath})

	}
}

func (s *server) uploadApks() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		r.ParseMultipartForm(10 << 20)

		file, handler, err := r.FormFile("apk")

		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		defer file.Close()

		dst, err := os.Create("./static/apks/" + handler.Filename)
		defer dst.Close()

		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		if _, err := io.Copy(dst, file); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		filePath := path.Join("./static/apks/" + handler.Filename)

		s.respond(w, r, http.StatusCreated, map[string]interface{}{"status": "successfully create", "filePath": filePath})

	}
}

func (s *server) getApks() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		entries, err := os.ReadDir("./static/apks")
		if err != nil {
			log.Fatal(err)
		}

		for _, e := range entries {
			fmt.Println(e.Name())
			s.respond(w, r, http.StatusOK, map[string]interface{}{"files": e.Name()})
		}
	}
}

func (s *server) downloadApk() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		filename := &fileName{}
		if err := json.NewDecoder(r.Body).Decode(filename); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		file, err := os.Open("./static/apks/" + filename.Name)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}
		defer file.Close()

		contentType := "application/octet-stream"
		w.Header().Set("Content-Type", contentType)
		w.Header().Set("Content-Disposition", "attachment; filename="+filename.Name)

		_, err = io.Copy(w, file)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}
	}
}

func (s *server) deleteCatalogPhoto() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		filename := &fileName{}
		if err := json.NewDecoder(r.Body).Decode(filename); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		filePath := "./static/catalog/" + filename.Name

		err := os.Remove(filePath)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		s.respond(w, r, http.StatusOK, map[string]interface{}{"status": "successfully delete", "filePath": filePath})
	}
}

func (s *server) deleteGalleryPhoto() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		filename := &fileName{}
		if err := json.NewDecoder(r.Body).Decode(filename); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		filePath := "./static/gallery/" + filename.Name

		err := os.Remove(filePath)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		s.respond(w, r, http.StatusOK, map[string]interface{}{"status": "successfully delete", "filePath": filePath})
	}
}

func (s *server) deleteProfilePhoto() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		filename := &fileName{}
		if err := json.NewDecoder(r.Body).Decode(filename); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		filePath := "./static/profile/" + filename.Name

		err := os.Remove(filePath)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		s.respond(w, r, http.StatusOK, map[string]interface{}{"status": "successfully delete", "filePath": filePath})
	}
}

func (s *server) deleteReviewsPhoto() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		filename := &fileName{}
		if err := json.NewDecoder(r.Body).Decode(filename); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		filePath := "./static/reviews/" + filename.Name

		err := os.Remove(filePath)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		s.respond(w, r, http.StatusOK, map[string]interface{}{"status": "successfully delete", "filePath": filePath})
	}
}

func (s *server) deleteApk() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		filename := &fileName{}
		if err := json.NewDecoder(r.Body).Decode(filename); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		filePath := "./static/apks/" + filename.Name

		err := os.Remove(filePath)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		s.respond(w, r, http.StatusOK, map[string]interface{}{"status": "successfully delete", "filePath": filePath})
	}
}

func (s *server) handleGetPhoto() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fileName := r.URL.Query().Get("filename")
		if fileName == "" {
			s.error(w, r, http.StatusBadRequest, errors.New("filename parameter is required"))
			return
		}

		file, err := os.Open(fileName)
		if err != nil {
			s.error(w, r, http.StatusNotFound, err)
			return
		}
		defer file.Close()

		w.Header().Set("Content-Type", "image/jpeg")
		w.Header().Set("Content-Type", "image/png")

		if _, err := io.Copy(w, file); err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}
	}
}
