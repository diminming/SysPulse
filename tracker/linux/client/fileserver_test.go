package client

import (
	"testing"

	"github.com/google/uuid"
)

func TestUpload2FileServer(t *testing.T) {
	Upload2FileServer("syspulse", uuid.NewString(), "/tmp/f7f42faa-147b-4a20-8958-ea254f51c288", "application/octet-stream")
}
