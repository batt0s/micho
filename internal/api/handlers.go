package api

import (
	"errors"
	"net/http"

	"github.com/batt0s/micho/internal/helm"
	"github.com/batt0s/micho/internal/logging"
	"github.com/go-chi/chi/v5"
)

func DeployPyMenuHandler(w http.ResponseWriter, r *http.Request) {
	body, err := decodeJSON[DeployRequest](w, r)
	if err != nil {
		var mr *malformedRequest
		if errors.As(err, *mr) {
			http.Error(w, mr.msg, mr.status)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	if err := body.validate(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	cfg := helm.ReleaseConfig{
		Namespace:   "tenant-" + body.Slug,
		ReleaseName: "pymenu-" + body.Slug,
		ChartPath:   "./charts/pymenu",
		Values:      body.ToHelmValues(),
	}

	if err := helm.InstallRelease(cfg); err != nil {
		http.Error(w, "Deployment Failed: "+err.Error(), http.StatusInternalServerError)
		logging.Record(body.Slug, "deploy", "failed", err)
		return
	}

	logging.Record(body.Slug, "deploy", "successfull", nil)

	respondJSON(w, http.StatusOK, map[string]string{
		"url": "https://" + body.Domain,
	})
}

func UpgradePyMenuHandler(w http.ResponseWriter, r *http.Request) {
	body, err := decodeJSON[DeployRequest](w, r)
	if err != nil {
		var mr *malformedRequest
		if errors.As(err, *mr) {
			http.Error(w, mr.msg, mr.status)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	if err := body.validate(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	cfg := helm.ReleaseConfig{
		Namespace:   "tenant-" + body.Slug,
		ReleaseName: "pymenu-" + body.Slug,
		ChartPath:   "./charts/pymenu",
		Values:      body.ToHelmValues(),
	}

	if err := helm.UpgradeRelease(cfg); err != nil {
		http.Error(w, "Upgrade Failed: "+err.Error(), http.StatusInternalServerError)
		logging.Record(body.Slug, "upgrade", "failed", err)
		return
	}

	logging.Record(body.Slug, "upgrade", "successfull", nil)

	respondJSON(w, http.StatusOK, map[string]string{
		"url": "https://" + body.Domain,
	})
}

func UninstallPyMenuHandler(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	if slug == "" {
		http.Error(w, "Slug parameter missing", http.StatusBadRequest)
		return
	}

	namespace := "tenant-" + slug
	releaseName := "pymenu-" + slug

	if err := helm.UninstallRelease(namespace, releaseName); err != nil {
		http.Error(w, "Uninstall Failed: "+err.Error(), http.StatusInternalServerError)
		logging.Record(slug, "uninstall", "failed", err)
		return
	}

	logging.Record(slug, "uninstall", "successfull", nil)

	respondJSON(w, http.StatusOK, map[string]string{})
}

func StatusPyMenuHandler(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	if slug == "" {
		http.Error(w, "Slug parameter missing", http.StatusBadRequest)
		return
	}

	namespace := "tenant-" + slug
	releaseName := "pymenu-" + slug

	status, err := helm.GetReleaseStatus(namespace, releaseName)
	if err != nil {
		http.Error(w, "Release Not Found", http.StatusNotFound)
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{
		"release": releaseName,
		"status":  status,
	})
}
