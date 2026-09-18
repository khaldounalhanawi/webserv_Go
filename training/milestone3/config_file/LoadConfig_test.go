package LoadConfig

import "testing"

func TestLoadConfig(t *testing.T) {
	_, err := LoadConfig("config.config")
	if err != nil {t.Fatal("failz")}

}
