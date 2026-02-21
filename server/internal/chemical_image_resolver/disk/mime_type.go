package chemicalimageresolver_disk

import (
	"mime"
	"strings"
)

// imageExtensionToMime maps lowercase file extensions (with leading dot) to MIME types.
// Extensions not in this map fall back to mime.TypeByExtension, then application/octet-stream.
var imageExtensionToMime = map[string]string{
	".png":  "image/png",
	".jpg":  "image/jpeg",
	".jpeg": "image/jpeg",
	".gif":  "image/gif",
	".svg":  "image/svg+xml",
	".bmp":  "image/bmp",
	".webp": "image/webp",
	".ico":  "image/x-icon",
	".tiff": "image/tiff",
	".tif":  "image/tiff",
}

// mimeToExtension maps MIME types to a canonical file extension (with leading dot).
// Prefers common extensions (e.g. ".jpg" for image/jpeg).
var mimeToExtension = map[string]string{
	"image/png":     ".png",
	"image/jpeg":    ".jpg",
	"image/gif":     ".gif",
	"image/svg+xml": ".svg",
	"image/bmp":     ".bmp",
	"image/webp":    ".webp",
	"image/x-icon":  ".ico",
	"image/tiff":    ".tiff",
}

const (
	fallbackMimeType  = "application/octet-stream"
	fallbackExtension = ".bin"
)

// ExtensionToMimeType returns the MIME type for a file path or extension.
// The argument may be a full path (e.g. "/path/to/image.PNG") or just an extension (e.g. ".png" or "png").
// Lookup is case-insensitive. If the type is unknown, falls back to the standard mime package
// and then to application/octet-stream.
func ExtensionToMimeType(pathOrExt string) string {
	ext := normalizeExtension(pathOrExt)
	if ext == "" {
		return fallbackMimeType
	}
	if mimeType, ok := imageExtensionToMime[ext]; ok {
		return mimeType
	}
	if mimeType := mime.TypeByExtension(ext); mimeType != "" {
		return mimeType
	}
	return fallbackMimeType
}

// MimeTypeToExtension returns a file extension (with leading dot) for the given MIME type.
// The argument is normalized (trimmed, lowercase) before lookup. If the type is unknown,
// falls back to the standard mime.ExtensionsByType and then to ".bin".
func MimeTypeToExtension(mimeType string) string {
	norm := strings.ToLower(strings.TrimSpace(mimeType))
	if norm == "" {
		return fallbackExtension
	}
	if ext, ok := mimeToExtension[norm]; ok {
		return ext
	}
	if exts, err := mime.ExtensionsByType(norm); err == nil && len(exts) > 0 {
		ext := exts[0]
		if ext != "" && ext[0] != '.' {
			return "." + ext
		}
		return ext
	}
	return fallbackExtension
}

// normalizeExtension returns a lowercase extension with a leading dot, e.g. ".png".
// Accepts either a full path or a raw extension (with or without dot).
func normalizeExtension(pathOrExt string) string {
	s := strings.TrimSpace(pathOrExt)
	if s == "" {
		return ""
	}
	// If it looks like a path, take the last component's extension
	if strings.Contains(s, "/") || strings.Contains(s, "\\") {
		s = s[strings.LastIndexAny(s, "/\\")+1:]
	}
	ext := ""
	if i := strings.LastIndexByte(s, '.'); i >= 0 && i < len(s)-1 {
		ext = s[i:]
	} else if s != "" && s[0] != '.' {
		ext = "." + s
	} else if s != "" {
		ext = s
	}
	return strings.ToLower(ext)
}
