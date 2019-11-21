package lib

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"image"
	"io/ioutil"
	"strings"

	guuid "github.com/google/uuid"
)

func UploadteImage(file string) (string, string, error) {
	idx := strings.Index(file, ";base64,")
	if idx < 0 {
		fmt.Println("=======1")
		return "", "", errors.New("ErrInvalidImage")
	}
	reader := base64.NewDecoder(base64.StdEncoding, strings.NewReader(file[idx+8:]))
	buff := bytes.Buffer{}
	_, err := buff.ReadFrom(reader)
	if err != nil {
		fmt.Println("=======2", err)
		// return "", err
	}
	imgCfg, fm, err := image.DecodeConfig(bytes.NewReader(buff.Bytes()))
	if err != nil {
		fmt.Println("=======3", err)
		// return "", err
	}

	if imgCfg.Width != 750 || imgCfg.Height != 685 {
		fmt.Println("=======4", err)
		// return "", errors.New("ErrSize 750")
	}
	if fm == "" {
		fm = ".jpg"
	}

	u := guuid.New()
	fileNameBase := "public/assets/images/%s"
	fileNameBase = fmt.Sprintf(fileNameBase, u)
	fileName := fileNameBase + fm
	fileNameBase = fmt.Sprintf("/images/%s", u) + fm

	err = ioutil.WriteFile(fileName, buff.Bytes(), 0644)
	return fileNameBase, fileName, nil
}
