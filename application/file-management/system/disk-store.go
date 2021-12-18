package system

import (
	"bytes"
	"fmt"
	"os"
	"sync"

	"raspstore.github.io/file-manager/model"
)

type DiskStore interface {
	Save(file *model.File, data bytes.Buffer) error
	Delete(uri string) error
}

type diskStore struct {
	mutex      sync.RWMutex
	rootFolder string
}

func NewDiskStore(rootFolder string) DiskStore {
	return &diskStore{rootFolder: rootFolder}
}

func (d *diskStore) Save(file *model.File, data bytes.Buffer) error {

	path := fmt.Sprintf("%s/%s", d.rootFolder, file.Filename)

	touch, err := os.Create(path)

	if err != nil {
		fmt.Println("could not create file ", file.Filename, " due to error: ", err.Error())
		return err
	}

	if _, err := data.WriteTo(touch); err != nil {
		fmt.Println("an error occurred while writing file to disk: ", err.Error())
		return err
	}

	d.mutex.Lock()

	defer d.mutex.Unlock()

	file.Uri = path

	return nil
}

func (d *diskStore) Delete(uri string) error {
	return os.Remove(uri)
}
