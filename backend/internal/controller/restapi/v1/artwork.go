package v1

import (
	"errors"
	"mime"
	"mime/multipart"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/maksimovyuriy/artfolio/backend/internal/controller/restapi/apierror"
	"github.com/maksimovyuriy/artfolio/backend/internal/controller/restapi/jsonutil"
	"github.com/maksimovyuriy/artfolio/backend/internal/controller/restapi/v1/request"
	"github.com/maksimovyuriy/artfolio/backend/internal/usecase"
)

const multipartMemoryLimit = 1 << 20

func (c *Controller) listPublishedArtworks(w http.ResponseWriter, r *http.Request) {
	artworks, err := c.artwork.ListPublished(r.Context())
	if err != nil {
		apierror.Write(w, r, err)
		return
	}
	_ = jsonutil.Encode(w, http.StatusOK, c.artworkMapper.FromEntities(artworks))
}

func (c *Controller) listAllArtworks(w http.ResponseWriter, r *http.Request) {
	artworks, err := c.artwork.ListAll(r.Context())
	if err != nil {
		apierror.Write(w, r, err)
		return
	}
	_ = jsonutil.Encode(w, http.StatusOK, c.artworkMapper.AdminFromEntities(artworks))
}

func (c *Controller) createArtwork(w http.ResponseWriter, r *http.Request) {
	form, status, err := c.parseMultipart(w, r)
	if err != nil {
		apierror.Write(w, r, apierror.WithStatus(status, err))
		return
	}
	defer form.RemoveAll()

	body, err := request.ArtworkFromValues(form.Value)
	if err != nil {
		apierror.Write(w, r, apierror.InvalidRequest(err))
		return
	}
	image, status, err := openImage(form, true)
	if err != nil {
		apierror.Write(w, r, apierror.WithStatus(status, err))
		return
	}
	defer image.Close()

	created, err := c.artwork.Create(r.Context(), body.Artwork(0), image)
	if err != nil {
		apierror.Write(w, r, err)
		return
	}
	_ = jsonutil.Encode(w, http.StatusCreated, c.artworkMapper.AdminFromEntity(created))
}

func (c *Controller) updateArtwork(w http.ResponseWriter, r *http.Request) {
	id, err := artworkID(r)
	if err != nil {
		apierror.Write(w, r, apierror.InvalidRequest(err))
		return
	}
	form, status, err := c.parseMultipart(w, r)
	if err != nil {
		apierror.Write(w, r, apierror.WithStatus(status, err))
		return
	}
	defer form.RemoveAll()

	body, err := request.ArtworkFromValues(form.Value)
	if err != nil {
		apierror.Write(w, r, apierror.InvalidRequest(err))
		return
	}
	image, status, err := openImage(form, false)
	if err != nil {
		apierror.Write(w, r, apierror.WithStatus(status, err))
		return
	}
	if image != nil {
		defer image.Close()
	}

	if err := c.artwork.Update(r.Context(), body.Artwork(id), image); err != nil {
		apierror.Write(w, r, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (c *Controller) deleteArtwork(w http.ResponseWriter, r *http.Request) {
	id, err := artworkID(r)
	if err != nil {
		apierror.Write(w, r, apierror.InvalidRequest(err))
		return
	}
	if err := c.artwork.Delete(r.Context(), id); err != nil {
		apierror.Write(w, r, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (c *Controller) reorderArtworks(w http.ResponseWriter, r *http.Request) {
	var body request.ReorderArtworks
	if err := jsonutil.Decode(w, r, &body); err != nil {
		apierror.Write(w, r, apierror.InvalidRequest(err))
		return
	}
	if err := c.artwork.Reorder(r.Context(), body.ArtworkIDs); err != nil {
		apierror.Write(w, r, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (c *Controller) parseMultipart(w http.ResponseWriter, r *http.Request) (*multipart.Form, int, error) {
	mediaType, _, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
	if err != nil || mediaType != "multipart/form-data" {
		return nil, http.StatusUnsupportedMediaType, errors.New("multipart/form-data is required")
	}
	if err := r.ParseMultipartForm(multipartMemoryLimit); err != nil {
		var maxBytesError *http.MaxBytesError
		if errors.As(err, &maxBytesError) {
			return nil, http.StatusRequestEntityTooLarge, err
		}
		return nil, http.StatusBadRequest, err
	}
	return r.MultipartForm, 0, nil
}

func openImage(form *multipart.Form, required bool) (multipart.File, int, error) {
	for key := range form.File {
		if key != "image" {
			return nil, http.StatusBadRequest, errors.New("unknown file field")
		}
	}
	files := form.File["image"]
	if len(files) == 0 {
		if required {
			return nil, http.StatusBadRequest, usecase.ErrArtworkImageRequired
		}
		return nil, 0, nil
	}
	if len(files) != 1 {
		return nil, http.StatusBadRequest, errors.New("exactly one image is allowed")
	}
	file, err := files[0].Open()
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	return file, 0, nil
}

func artworkID(r *http.Request) (int64, error) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil || id <= 0 {
		return 0, errors.New("invalid artwork id")
	}
	return id, nil
}
