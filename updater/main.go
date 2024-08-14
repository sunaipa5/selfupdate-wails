package updater

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"
)

// isUpdateAvailable - true
func (options Options) CheckUpdate() (bool, Release) {
	resp, err := http.Get("https://api.github.com/repos/" + options.Author + "/" + options.Repo + "/releases/latest")
	if err != nil {
		fmt.Println(err)
		return false, Release{}
	}
	defer resp.Body.Close()

	jsonDoc, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return false, Release{}
	}

	var release Release
	if err := json.Unmarshal(jsonDoc, &release); err != nil {
		fmt.Println(err)
		return false, Release{}
	}

	if release.Version != options.CurrentVersion {
		fmt.Println("New version available!")
		return true, release
	} else {
		fmt.Println("App up to date")
		return false, Release{}
	}

}

func (options Options) ApplyUpdate(release Release) {
	var source Source
	count := 0
	for _, asset := range release.Assets {
		if count > 0 {
			source.Name += " "
			source.Download_Url += " "
		}
		if strings.HasSuffix(asset.Name, options.TagEnd) {
			count++
			source.Name += asset.Name
			source.Download_Url += asset.Download_Url
		}
	}

	if count > 1 {
		fmt.Println("Multiple source found! please change 'TagEnd' :")
		fmt.Println(source)
		return
	} else if count < 1 {
		fmt.Println("Not found any source! plese change 'TagEnd'")
		return
	}

	GZ_extractor(source.Download_Url)
	installUpdateNew(options.AppName)

	releaseJson, err := json.Marshal(release)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(string(releaseJson))
}

func installUpdateNew(appName string) {
	var err error
	newBinary := ".tmp/"+appName
	oldBinary := "./"+appName
	/*oldBinary, err := os.Executable()
	if err != nil {
		fmt.Println("Already app location cannot get:", err)
		return
	}*/

	// Yeni binary'ye çalıştırma izni ver
	err = os.Chmod(newBinary, 0755)
	if err != nil {
		fmt.Println("Exec permission cannot get:", err)
		return
	}

	// Yeni binary'yi çalıştır
	cmd := exec.Command(newBinary, "updated")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Start()
	if err != nil {
		fmt.Println("New app cannot start:", err)
		return
	}

	// Kısa bir süre bekle
	time.Sleep(time.Second)

	// Eski binary'yi yenisiyle değiştir
	bakPath := oldBinary + ".bak"
	err = os.Rename(oldBinary, bakPath)
	if err != nil {
		fmt.Println("Old file cannot backup:", err)
		return
	}

	err = os.Rename(newBinary, oldBinary)
	if err != nil {
		// Hata durumunda eski dosyayı geri getir
		os.Rename(bakPath, oldBinary)
		fmt.Println("File cannot changed:", err)
		return
	}

	// Yedek dosyayı sil
	os.Remove(bakPath)

	fmt.Println("Update Finised")
	os.Exit(0)
}
