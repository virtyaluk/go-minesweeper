package main

import (
	"fmt"
	"github.com/visualfc/atk/tk"
)

// List of available image assets
var assetsList = []string{"1", "2", "3", "4", "5", "6", "7", "8", "bomb", "default", "empty", "flag"}

// AssetsManager is used to manipulate game assets
type AssetsManager struct {
	assets       map[string]*tk.Image
	assetsLoaded bool
}

// NewAssetsManager returns an instance of asset manager object.
func NewAssetsManager() *AssetsManager {
	return &AssetsManager{
		make(map[string]*tk.Image),
		false,
	}
}

// LoadAssets performs asset loading.
func (am *AssetsManager) LoadAssets() {
	if am.assetsLoaded {
		return
	}

	for _, assetName := range assetsList {
		loadedImg, err := tk.LoadImage(fmt.Sprintf("assets/%s.png", assetName))

		if err != nil {
			fmt.Println(err.Error())
			continue
		}

		//loadedImg.SetSizeN(32, 32)

		am.assets[assetName] = loadedImg
	}

	am.assetsLoaded = true
}

// GetAsset gets an asset given its name.
func (am *AssetsManager) GetAsset(assetName string) *tk.Image {
	if asset, ok := am.assets[assetName]; ok {
		return asset
	}

	return nil
}
