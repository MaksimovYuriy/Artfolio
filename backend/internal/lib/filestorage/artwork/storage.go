package artwork

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/maksimovyuriy/artfolio/backend/internal/config"
	"github.com/maksimovyuriy/artfolio/backend/internal/entity"
	"github.com/maksimovyuriy/artfolio/backend/internal/lib/filestorage"
)

type Storage struct {
	root        string
	temporary   string
	artworks    string
	maxFileSize int64
	maxPixels   int64
}

var _ filestorage.Artwork = (*Storage)(nil)

func New(cfg config.FileStorageConfig) (*Storage, error) {
	if cfg.Path == "" {
		return nil, errors.New("storage path is required")
	}
	if cfg.MaxFileSize <= 0 {
		return nil, errors.New("storage max file size must be positive")
	}
	if cfg.MaxPixels <= 0 {
		return nil, errors.New("storage max pixels must be positive")
	}

	root, err := filepath.Abs(cfg.Path)
	if err != nil {
		return nil, fmt.Errorf("resolve storage path: %w", err)
	}

	storage := &Storage{
		root:        root,
		temporary:   filepath.Join(root, ".tmp"),
		artworks:    filepath.Join(root, "artworks"),
		maxFileSize: cfg.MaxFileSize,
		maxPixels:   cfg.MaxPixels,
	}

	if err := os.MkdirAll(storage.artworks, 0o755); err != nil {
		return nil, fmt.Errorf("create artwork storage directory: %w", err)
	}
	if err := os.MkdirAll(storage.temporary, 0o700); err != nil {
		return nil, fmt.Errorf("create temporary storage directory: %w", err)
	}

	return storage, nil
}

func (s *Storage) Save(ctx context.Context, content io.Reader) (stored entity.StoredArtworkImage, resultErr error) {
	temporaryFile, err := os.CreateTemp(s.temporary, "upload-*")
	if err != nil {
		return stored, fmt.Errorf("create temporary artwork image: %w", err)
	}
	temporaryPath := temporaryFile.Name()
	defer func() {
		_ = temporaryFile.Close()
		if resultErr != nil {
			_ = os.Remove(temporaryPath)
		}
	}()

	limited := &io.LimitedReader{R: &contextReader{ctx: ctx, reader: content}, N: s.maxFileSize + 1}
	written, err := io.Copy(temporaryFile, limited)
	if err != nil {
		return stored, fmt.Errorf("write temporary artwork image: %w", err)
	}
	if written > s.maxFileSize {
		return stored, filestorage.ErrFileTooLarge
	}
	if written == 0 {
		return stored, filestorage.ErrInvalidImage
	}

	if _, err := temporaryFile.Seek(0, io.SeekStart); err != nil {
		return stored, fmt.Errorf("rewind temporary artwork image: %w", err)
	}
	imageConfig, format, err := image.DecodeConfig(temporaryFile)
	if err != nil {
		return stored, filestorage.ErrInvalidImage
	}
	extension, ok := extensionForFormat(format)
	if !ok || imageConfig.Width <= 0 || imageConfig.Height <= 0 {
		return stored, filestorage.ErrInvalidImage
	}
	if int64(imageConfig.Width) > s.maxPixels/int64(imageConfig.Height) {
		return stored, filestorage.ErrImageTooManyPixels
	}
	if _, err := temporaryFile.Seek(0, io.SeekStart); err != nil {
		return stored, fmt.Errorf("rewind temporary artwork image: %w", err)
	}
	if _, _, err := image.Decode(temporaryFile); err != nil {
		return stored, filestorage.ErrInvalidImage
	}

	name, err := randomName()
	if err != nil {
		return stored, fmt.Errorf("generate artwork image name: %w", err)
	}
	key := path.Join("artworks", name+extension)
	destination := filepath.Join(s.root, filepath.FromSlash(key))

	if err := temporaryFile.Sync(); err != nil {
		return stored, fmt.Errorf("sync temporary artwork image: %w", err)
	}
	if err := temporaryFile.Close(); err != nil {
		return stored, fmt.Errorf("close temporary artwork image: %w", err)
	}
	if err := os.Chmod(temporaryPath, 0o644); err != nil {
		return stored, fmt.Errorf("set artwork image permissions: %w", err)
	}
	if err := os.Rename(temporaryPath, destination); err != nil {
		return stored, fmt.Errorf("store artwork image: %w", err)
	}

	return entity.StoredArtworkImage{
		Key:    key,
		Width:  imageConfig.Width,
		Height: imageConfig.Height,
	}, nil
}

func (s *Storage) Delete(ctx context.Context, key string) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if !validKey(key) {
		return filestorage.ErrInvalidKey
	}

	filename := filepath.Join(s.root, filepath.FromSlash(key))
	relative, err := filepath.Rel(s.root, filename)
	if err != nil || relative == "." || strings.HasPrefix(relative, ".."+string(filepath.Separator)) {
		return filestorage.ErrInvalidKey
	}

	if err := os.Remove(filename); err != nil && !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("delete artwork image: %w", err)
	}
	return nil
}

func randomName() (string, error) {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func validKey(key string) bool {
	return key != "" &&
		key == path.Clean(key) &&
		strings.HasPrefix(key, "artworks/") &&
		!strings.Contains(key, "\\")
}

type contextReader struct {
	ctx    context.Context
	reader io.Reader
}

func (r *contextReader) Read(buffer []byte) (int, error) {
	if err := r.ctx.Err(); err != nil {
		return 0, err
	}
	return r.reader.Read(buffer)
}
