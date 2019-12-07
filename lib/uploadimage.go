package lib

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"image"
	"io/ioutil"
	"os"
	"strings"

	guuid "github.com/google/uuid"
)

const PathLocal = "public/assets"

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
		// return "", err
	}
	imgCfg, fm, err := image.DecodeConfig(bytes.NewReader(buff.Bytes()))
	if err != nil {
		// return "", err
	}

	if imgCfg.Width != 750 || imgCfg.Height != 685 {
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

func UploadFile(file string, typeFile string) (string, string, error) {
	u := guuid.New()
	fileNameBase := "public/assets/mp3/%s"
	fileNameBase = fmt.Sprintf(fileNameBase, u)
	fileName := fileNameBase + typeFile
	fileNameBase = fmt.Sprintf("/mp3/%s", u) + typeFile
	if err := Decode(file, fileName); err != nil {
		fmt.Println(err)
		return "", "", err
	}
	return fileNameBase, fileName, nil

}

func Decode(code string, dest string) error {
	idx := strings.Index(code, ";base64,")
	if idx < 0 {
		fmt.Println("=======1")
		return errors.New("ErrInvalidImage")
	}
	reader := base64.NewDecoder(base64.StdEncoding, strings.NewReader(code[idx+8:]))

	buff := bytes.Buffer{}
	buff.ReadFrom(reader)
	if err := ioutil.WriteFile(dest, buff.Bytes(), 0644); err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func DeleteFile(path string) error {
	// delete file
	path = PathLocal + path
	var err = os.Remove(path)
	if isError(err) {
		return err
	}

	fmt.Println("==> done deleting file")
	return nil
}

func isError(err error) bool {
	if err != nil {
		fmt.Println(err.Error())
	}

	return (err != nil)
}
