package data

import (
	"fmt"
	"path"

	"github.com/DaanV2/f1-game-dashboards/server/sessions"
	"github.com/charmbracelet/log"
	flag "github.com/spf13/pflag"
)

func NewStorage(flags *flag.FlagSet) (Database, error) {
	storageType := flags.Lookup("storage-type").Value.String()

	switch storageType {
	case "files":
		storageFolder := flags.Lookup("files-storage-directory").Value.String()
		if storageFolder == "" {
			storageFolder = path.Join(".", "data", "files")
		}

		return NewFileStorage(storageFolder), nil
	default:
		return nil, fmt.Errorf("unknown storage type: %s", storageType)
	}
}

// DatabaseHooks sets up hooks to store changes in the database
func DatabaseHooks(database Database, chairs *sessions.ChairManager) {
	chairs.OnChairAdded.Add(func(chair sessions.Chair) {
		err := database.Chairs().Set(chair.Id(), chair)
		if err != nil {
			log.Error("could not store chair", "error", err)
		}
	})
	chairs.OnChairUpdated.Add(func(chair sessions.Chair) {
		err := database.Chairs().Set(chair.Id(), chair)
		if err != nil {
			log.Error("could not store chair", "error", err)
		}
	})
	chairs.OnChairRemoved.Add(func(chair sessions.Chair) {
		database.Chairs()
	})
}
