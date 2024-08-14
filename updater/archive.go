package updater

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"os"
	
)

func GZ_extractor(url string) {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error downloading file:", err)
		return
	}
	defer resp.Body.Close()

	// Gzip ile aç
	gzReader, err := gzip.NewReader(resp.Body)
	if err != nil {
		fmt.Println("Error creating gzip reader:", err)
		return
	}
	defer gzReader.Close()

	// Tar arşivini aç
	tarReader := tar.NewReader(gzReader)

	tmpDir := ".tmp"
	if err := os.MkdirAll(tmpDir, os.ModePerm); err != nil {
		fmt.Errorf("error creating .tmp directory: %w", err)
		return
	}

	// Tar arşivindeki dosyaları çıkar
	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break // Tüm dosyalar çıkarıldı
		}
		if err != nil {
			fmt.Println("Error reading tar:", err)
			return
		}

		// Çıkarılacak dosyanın tam yolu
		//target := filepath.Join(".tmp", header.Name)
		target := ".tmp/selfupdate-wails"

		fmt.Println("TARGET1:",target)
		// Dizin oluştur
		if header.Typeflag == tar.TypeDir {
			if err := os.MkdirAll(target, os.ModePerm); err != nil {
				fmt.Println("Error creating directory:", err)
				return
			}
			continue
		}

		fmt.Println("TARGET2:",target)
		// Dosyayı çıkar
		outFile, err := os.Create(target)
		if err != nil {
			fmt.Println("Error creating file:", err)
			return
		}
		defer outFile.Close()

		// Dosyayı yaz
		if _, err := io.Copy(outFile, tarReader); err != nil {
			fmt.Println("Error writing file:", err)
			return
		}
	}

	fmt.Println("Extraction completed successfully.")
}
