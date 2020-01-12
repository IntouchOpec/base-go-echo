package lib

import (
	"bytes"
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"image"
	"io"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	"cloud.google.com/go/storage"
	guuid "github.com/google/uuid"
	"google.golang.org/api/option"
)

const PathLocal = "public/assets"

func UploadteImage(file string) (string, string, error) {
	idx := strings.Index(file, ";base64,")
	if idx < 0 {
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

func UploadGoolgeStorage(ctx context.Context, code, imagePath string) (string, error) {
	client, err := storage.NewClient(ctx, option.WithCredentialsFile("/lineconnect-99ca66b2bd16.json"))
	if err != nil {
		return "", err
	}
	idx := strings.Index(code, ";base64,")
	if idx < 0 {
		return "", errors.New("ErrInvalidImage")
	}

	reader := base64.NewDecoder(base64.StdEncoding, strings.NewReader(code[idx+8:]))
	imagePath = fmt.Sprintf("%s%s", imagePath, guuid.New())
	wc := client.Bucket("triple-t").Object(imagePath).NewWriter(ctx)
	if err != nil {
		return "", err
	}
	fmt.Println(code)
	if _, err := io.Copy(wc, reader); err != nil {
		fmt.Println(err, "===3")
		return "", err
	}
	if err := wc.Close(); err != nil {
		fmt.Println(err)
		return "", err
	}
	return imagePath, nil
}

func GetGoolgeStorage(ctx context.Context, bucket, folder string) ([]byte, string, error) {
	client, err := storage.NewClient(ctx, option.WithCredentialsFile("/lineconnect-99ca66b2bd16.json"))
	if err != nil {
		return nil, "", err
	}
	rc, err := client.Bucket(bucket).Object(folder).NewReader(ctx)
	if err != nil {
		return nil, "", err
	}
	defer rc.Close()
	data, err := ioutil.ReadAll(rc)
	if err != nil {
		return nil, "", err
	}
	length := len(data)
	sz := strconv.Itoa(length)

	return data, sz, nil
}
