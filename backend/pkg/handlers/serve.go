package handlers

import (
	"net/http"

	"social-network/pkg/models"
	"social-network/pkg/tools"
)

func (rt *Root) NewServeFilesHandler() (serveFilesMux *http.ServeMux) {
	serveFilesMux = http.NewServeMux()

	serveFilesMux.HandleFunc("GET /img", rt.ServeImage)

	return serveFilesMux
}

func (rt *Root) ServeImage(w http.ResponseWriter, r *http.Request) {
	imgUUID := r.URL.Query().Get("img_uuid")
	table := r.URL.Query().Get("table")

	if imgUUID == "" || table == "" {
		rt.DL.Logger.Log(models.LogEntry{
			Level:   "ERROR",
			Message: "Missing image UUID or table query params",
			Metadata: map[string]any{
				"ip":   r.RemoteAddr,
				"path": r.URL.Path,
			},
		})
		tools.RespondError(w, "Missing query parameters", http.StatusBadRequest)
		return
	}
	rt.DL.Logger.Log(models.LogEntry{Level: "DEBUG", Message: "Image GET request received"})

	requesterID := rt.DL.GetRequesterID(w, r)
	if requesterID == 0 {
		rt.DL.Logger.Log(models.LogEntry{
			Level:   "ERROR",
			Message: "Unauthorized request (no requester ID)",
			Metadata: map[string]any{
				"ip":   r.RemoteAddr,
				"path": r.URL.Path,
			},
		})
		tools.RespondError(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	err := rt.DL.Images.IsAuthorized(imgUUID, table, requesterID)
	if err != nil {
		rt.DL.Logger.Log(models.LogEntry{
			Level:   "ERROR",
			Message: "Image not found or access denied",
			Metadata: map[string]any{
				"uuid":   imgUUID,
				"table":  table,
				"userID": requesterID,
				"err":    err.Error(),
			},
		})
		tools.RespondError(w, "Access denied or image not found", http.StatusForbidden)
		return
	}
	rt.DL.Logger.Log(models.LogEntry{Level: "DEBUG", Message: "Authorized to access image"})

	if err := models.ImageServe(w, imgUUID); err != nil {
		rt.DL.Logger.Log(models.LogEntry{
			Level:   "ERROR",
			Message: "Error opening image file",
			Metadata: map[string]any{
				"image_uuid": imgUUID,
				"error":      err.Error(),
			},
		})
		tools.RespondError(w, "Failed to open image", http.StatusInternalServerError)
	}

	rt.DL.Logger.Log(models.LogEntry{
		Level:   "INFO",
		Message: "Image served successfully",
		Metadata: map[string]any{
			"uuid":   imgUUID,
			"table":  table,
			"userID": requesterID,
			"path":   r.URL.Path,
		},
	})
}
