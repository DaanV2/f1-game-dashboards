package jwt

import (
	"encoding/json"
	"errors"

	"github.com/DaanV2/f1-game-dashboards/server/pkg/data"
	"github.com/charmbracelet/log"
)

func GetOrCreate(database data.Database, generateNew bool) ([]*SigningInfo, error) {
	sigs, err := loadSigningInfo(database)
	if errors.Is(err, data.ErrNotFound) {
		generateNew = true
	}
	if err != nil {
		return sigs, err
	}
	if generateNew {
		sig, err := GenerateSigningInfo()
		if err != nil {
			log.Error("Failed to generate signing info", "error", err)
			return nil, err
		}
		sigs = append(sigs, sig)

		go func() {
			if err = storeSigningInfo(database, sigs); err != nil {
				log.Warn("Failed to store signing info", "error", err)
			}
		}()
	}

	return sigs, err
}

func loadSigningInfo(database data.Database) ([]*SigningInfo, error) {
	data, err := database.Config().Get("jwks")
	if err != nil {
		return nil, err
	}

	result := make([]*SigningInfo, 0)
	err = json.Unmarshal(data, &result)

	return result, err
}

func storeSigningInfo(database data.Database, sigs []*SigningInfo) error {
	data, err := json.Marshal(sigs)
	if err != nil {
		return err
	}

	return database.Config().Set("jwks", data)
}