package logger

import (
	"github.com/Cyborg0X/ASRS/Agent/internal/pkg/handler"
)

func DetectionMarker() bool {
	var detector handler.Config
	if detector.Detectionmarker.Markerisdetected {
		return detector.Detectionmarker.Markerisdetected
	}
	return false
}
