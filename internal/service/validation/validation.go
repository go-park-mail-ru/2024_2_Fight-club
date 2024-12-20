package validation

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"net/http"
	"regexp"
	"strings"
)

func ValidateEmail(email string) bool {
	re := regexp.MustCompile(`^.+@.+$`)
	return re.MatchString(email)
}

func ValidateLogin(login string) bool {
	re := regexp.MustCompile(`^[A-Za-zА-Яа-яЁё0-9][A-Za-zА-Яа-яЁё0-9-_.]{3,20}[A-Za-zА-Яа-яЁё0-9]$`)
	return re.MatchString(login)
}

func ValidatePassword(password string) bool {
	re := regexp.MustCompile(`^[a-zA-ZА-Яа-яЁё0-9!@#$%^&*()_+=-]{8,16}$`)
	return re.MatchString(password)
}

func ValidateName(name string) bool {
	re := regexp.MustCompile(`^[A-Za-zА-Яа-яЁё\s.]{5,50}$`)
	return re.MatchString(name)
}

func ValidateImages(files [][]byte, maxSize int64, allowedMimeTypes []string, maxWidth, maxHeight int) error {
	for _, file := range files {
		if err := ValidateImage(file, maxSize, allowedMimeTypes, maxWidth, maxHeight); err != nil {
			return err
		}
	}
	return nil
}

func ValidateImage(file []byte, maxSize int64, allowedMimeTypes []string, maxWidth, maxHeight int) error {
	if int64(len(file)) > maxSize {
		return fmt.Errorf("file exceeds maximum size of %d bytes", maxSize)
	}

	buffer := file[:512]
	mimeType := http.DetectContentType(buffer)
	allowed := false
	for _, t := range allowedMimeTypes {
		if mimeType == t {
			allowed = true
			break
		}
	}
	if !allowed {
		return errors.New("file type is not allowed, please use (png, jpg, jpeg) types")
	}

	var img image.Image
	var err error
	reader := bytes.NewReader(file)

	switch {
	case strings.HasSuffix(mimeType, "jpeg") || strings.HasSuffix(mimeType, "jpg"):
		img, err = jpeg.Decode(reader)
	case strings.HasSuffix(mimeType, "png"):
		img, err = png.Decode(reader)
	default:
		return errors.New("unsupported image format")
	}
	if err != nil {
		return errors.New("could not decode image")
	}

	width := img.Bounds().Dx()
	height := img.Bounds().Dy()
	if width > maxWidth || height > maxHeight {
		return errors.New("image resolution exceeds maximum allowed size of 2000 x 2000")
	}

	return nil
}
