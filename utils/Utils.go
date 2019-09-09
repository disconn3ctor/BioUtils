package utils

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/dchest/uniuri"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net"
	"net/http"
	"os"
	"time"
)

const (
	_                     = iota
	KB                int = 1 << (10 * iota)
	MB
	GB
	RandomChar = "abcdefghijklmnopqrstuvwxyz0123456789_"
)

type Sizer interface {
	Size() int64
}

func GetAllFormRequestValue(r *http.Request) map[string]interface{} {
	clearMapData := make(map[string]interface{})


	// chon r.Form tamame maghadir ro to array mirikht on maghidir ro az array dar avrodam
	for i, value := range r.Form {

		clearMapData[i] = value[0]

	}

	return clearMapData
}

func GetLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, address := range addrs {
		// check the address type and if it is not a loopback the display it
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return ""
}



func WriteAllPostImageFromRequest(r *http.Request, keyFileValue string, path string , maxWith int , maxHeight int , maxSize int) (chan string, error) {

	err := r.ParseMultipartForm(32 << 20) //32 MB

	if err != nil{
		return nil, err
	}

	allFile := r.MultipartForm.File[keyFileValue]

	fileCount := len(allFile)

	if fileCount > 0 {

		if err := FolderMaker(path); err != nil {
			fmt.Println(err.Error())
			return nil, err
		}

		fileNameChan := make(chan string, fileCount)

		for _, value := range allFile {

			if fileName, errImageWriter := ImageWriterByFileHeader(value , path , maxWith , maxHeight , maxSize ); errImageWriter != nil{
				return nil, errImageWriter
			}else
			{
				fileNameChan <- fileName
			}

		}

		return fileNameChan, nil

	}


	return nil, nil

}


func ImageWriterByFileHeader(fileHeader *multipart.FileHeader, path string , maxWith int , maxHeight int , maxSize int) (string, error) {

	file, err := fileHeader.Open()

	if err != nil{
		return "", err
	}

	buffer := make([]byte, fileHeader.Size)

	for {

		value, err := file.Read(buffer)

		if err != nil && err != io.EOF {
			return "", err
		}
		if value == 0 {
			break
		}
	}

	if err := file.Close(); err != nil {
		return "", err
	}

	conf, format, err := image.DecodeConfig(bytes.NewReader(buffer))
	if err != nil {
		return "", err
	}

	if format != "jpeg" && format != "png" && format != "jpg" {
		return "", errors.New(" format haye jpeg , png , jpg pazirofte mishavad")
	}

	if conf.Height > maxHeight || conf.Width > maxWith {

		return "", errors.New(fmt.Sprintf("tasvir bayad kochaktar az %d x %d pixel bashand", maxWith , maxHeight))
	}

	if len(buffer) >= maxSize*MB {
		return "", errors.New(fmt.Sprintf("tasvir bayad kochaktar az %d hajm dashte bashand", MB))
	}

	fileName := uniuri.NewLenChars(10, []byte(RandomChar)) + fmt.Sprint(time.Now().Unix())

	if errIo := ioutil.WriteFile(path+fileName+"."+format, buffer, 0700); errIo != nil {
		return "", errIo
	}

	return fileName + "." + format, nil
}


func WriteAllPostAudioFromRequest(r *http.Request, keyFileValue string, path string , maxSize int) (chan string, error) {

	err := r.ParseMultipartForm(32 << 20) //32 MB

	if err != nil {
		return nil, err
	}

	allFile := r.MultipartForm.File[keyFileValue]

	fileCount := len(allFile)

	if fileCount > 0 {

		if err := FolderMaker(path); err != nil {
			fmt.Println(err.Error())
			return nil, err
		}

		fileNameChan := make(chan string, fileCount)

		for _, value := range allFile {


			if fileName, errAudioWriter := AudioWriterByFileHeader(value , path , maxSize ); errAudioWriter != nil{
				return nil, errAudioWriter
			}else
			{
				fileNameChan <- fileName
			}

		}

		return fileNameChan, nil

	}

	return nil, nil
}

func AudioWriterByFileHeader(fileHeader *multipart.FileHeader,  path string, maxSize int)(string, error) {

	format := fileHeader.Filename[len(fileHeader.Filename)-3:]

	if format != "mp3" && format != "wav" && format != "aac" {
		return "", errors.New(" format haye mp3 , wav , aac pazirofte mishavad")
	}

	file, err := fileHeader.Open()

	if err != nil{
		return "", err
	}

	buffer := make([]byte, fileHeader.Size)

	for {

		value, err := file.Read(buffer)

		if err != nil && err != io.EOF {
			return "", err
		}
		if value == 0 {
			break
		}
	}

	if err := file.Close(); err != nil {
		return "", err
	}


	if len(buffer) >= maxSize*MB {
		return "", errors.New(fmt.Sprintf("file bayad kochaktar az %d MB hajm dashte bashand", maxSize))
	}

	fileName := uniuri.NewLenChars(10, []byte(RandomChar)) + fmt.Sprint(time.Now().Unix())


	if errIo := ioutil.WriteFile(path+fileName+"."+format, buffer, 0700); errIo != nil {
		return "", errIo
	}

	return fileName + "." + format, nil
}

func FolderMaker(path string) error {

	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			if errMkDir := os.MkdirAll(path, 0700); err != nil {
				return errMkDir
			}
		}
	} else {
		return err
	}

	return nil
}