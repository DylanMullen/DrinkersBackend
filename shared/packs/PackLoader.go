package packs

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type PackLoader struct {
	Packs map[string]Pack

	Location string
}

func (loader *PackLoader) LoadAllPacks() {
	if loader.Location == "" {
		return
	}

	loader.Packs = make(map[string]Pack)
	packFiles, err := ioutil.ReadDir(loader.Location)

	if err != nil {
		return
	}

	for _, v := range packFiles {
		loader.loadPack(v.Name())
	}

}

func (load *PackLoader) loadPack(file string) {
	jsonFile, err := os.Open(load.Location + "/" + file)

	if err != nil {
		fmt.Println("error")
		return
	}

	defer jsonFile.Close()

	contents, err := ioutil.ReadAll(jsonFile)

	if err != nil {
		return
	}

	pack, err := unmarshalPack(contents)

	if err != nil {
		return
	}

	load.Packs[pack.Settings.UUID] = pack
}

func unmarshalPack(contents []byte) (pack Pack, err error) {
	err = json.Unmarshal(contents, &pack)
	return
}
