package initdatabase

import (
	"github.com/kawasoki/gowork/configs"
	"testing"
)

func TestNewEsClient(t *testing.T) {
}

func GetEs() {
	NewEsClient(&configs.Config{})
}
