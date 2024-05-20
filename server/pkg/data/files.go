package data

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"strings"
	"sync"

	"github.com/DaanV2/f1-game-dashboards/server/sessions"
	"github.com/charmbracelet/log"
)

type (
	FileStorage struct {
		folder string

		chairs *TypedStorage[sessions.Chair]
		config *DirectoryStorage
	}

	DirectoryStorage struct {
		folder string
		lock   sync.Mutex
	}
)

func NewFileStorage(folder string) *FileStorage {
	if !path.IsAbs(folder) {
		folder = path.Join(".", folder)
	}
	folder = path.Clean(folder)

	log.Info("starting file storage", "folder", folder)
	checkFolder(folder)

	return &FileStorage{
		folder: folder,

		chairs: NewTypedStorage[sessions.Chair](NewDirectoryStorage(path.Join(folder, "chairs"))),
		config: NewDirectoryStorage(path.Join(folder, "config")),
	}
}

func (fs *FileStorage) Chairs() Storage[sessions.Chair] {
	return fs.chairs
}

func (fs *FileStorage) Config() RawStorage {
	return fs.config
}

func NewDirectoryStorage(folder string) *DirectoryStorage {
	checkFolder(folder)

	return &DirectoryStorage{
		folder: folder,
		lock:   sync.Mutex{},
	}
}

func (ds *DirectoryStorage) Get(id string) ([]byte, error) {
	ds.lock.Lock()
	defer ds.lock.Unlock()
	filepath := ds.filepath(id)
	log.Debug("loading from storage", "id", id, "filepath", filepath)

	data, err := os.ReadFile(filepath)
	if os.IsNotExist(err) {
		return data, ErrNotFound
	}
	return data, err
}

func (ds *DirectoryStorage) Set(id string, value []byte) error {
	ds.lock.Lock()
	defer ds.lock.Unlock()

	filepath := ds.filepath(id)
	log.Debug("saving to storage", "id", id, "filepath", filepath, "value", value)

	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return os.WriteFile(filepath, data, 0644)
}

func (ds *DirectoryStorage) Delete(id string) error {
	ds.lock.Lock()
	defer ds.lock.Unlock()

	filepath := ds.filepath(id)
	log.Debug("deleting from storage", "id", id, "filepath", filepath)

	return os.Remove(filepath)
}

func (ds *DirectoryStorage) Keys() []string {
	ds.lock.Lock()
	defer ds.lock.Unlock()

	files, err := os.ReadDir(ds.folder)
	if err != nil {
		panic(fmt.Errorf("could not read storage directory;%s; %w", ds.folder, err))
	}

	keys := make([]string, 0, len(files))
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		if !strings.HasSuffix(file.Name(), ".json") {
			continue
		}

		filename := file.Name()
		ext := path.Ext(filename)
		filename = filename[:len(filename)-len(ext)]

		keys = append(keys, filename)
	}

	return keys
}

func (ds *DirectoryStorage) filepath(id string) string {
	return path.Join(ds.folder, fmt.Sprintf("%s.json", id))
}

func checkFolder(folder string) {
	_, err := os.Stat(folder)
	if err == nil {
		return
	}

	err = os.MkdirAll(folder, 0755)
	if err != nil {
		log.Fatal("could not create storage folder", "folder", folder, "error", err)
	}
}
