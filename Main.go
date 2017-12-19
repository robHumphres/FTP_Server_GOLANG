package main

// https://thenewstack.io/make-a-restful-json-api-go/
// https://github.com/MedBridge/sample-app/blob/master/main.go
import (
	"archive/zip"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gorilla/mux"
	"github.com/mholt/archiver"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

//Attempt of comment
func UploadFile(w http.ResponseWriter, r *http.Request) {
	//Post uploads a single file
	if r.Method == "POST" {
		file, handler, err := r.FormFile("file")
		fmt.Printf("Post came through\n")
		defer file.Close()

		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Fprintf(w, "%v", handler.Header)
		f, err := os.OpenFile(handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			fmt.Println(err)
			return
		}

		defer f.Close()
		io.Copy(f, file)

		//Delete older ones if past 10
		UnzipNClean(handler.Filename)

	} else {
		fmt.Fprintf(w, "This is just a POST Method, see documentation")
	}

	return
}

func UnzipNClean(fileToUnzip string) {

	pathString, err := os.Getwd()
	if err != nil {
		check(err)
	}

	files, _ := ioutil.ReadDir(pathString)
	fmt.Println(len(files))

	fmt.Println("File to Unzip's name is.... : " + fileToUnzip)

	//Unzip the folder
	errr := archiver.Zip.Open(fileToUnzip, "")

	if errr != nil {
		check(errr)
	}

	//Delete the folder
	os.Remove(fileToUnzip)

	fmt.Println("Deleted the old file")

	if len(files) > 12 {
		deleteOldFolder()
	}
}

func deleteOldFolder() {
	/*
		todo, need to grab a list of the files in the directory sorted by file mod date, delete the last added
	*/
	fmt.Println("Made it to delete a file")
	// files, error := ioutil.ReadDir(".")
	// if error != nil {
	// 	panic(error)
	// }

	// for _, file := range files {
	// 	fmt.Println(file.Name())
	// }

}

func Unzip(src, dest string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer func() {
		if err := r.Close(); err != nil {
			panic(err)
		}
	}()

	os.MkdirAll(dest, 0755)

	// Closure to address file descriptors issue with all the deferred .Close() methods
	extractAndWriteFile := func(f *zip.File) error {
		rc, err := f.Open()
		if err != nil {
			return err
		}
		defer func() {
			if err := rc.Close(); err != nil {
				panic(err)
			}
		}()

		path := filepath.Join(dest, f.Name)

		if f.FileInfo().IsDir() {
			os.MkdirAll(path, f.Mode())
		} else {
			os.MkdirAll(filepath.Dir(path), f.Mode())
			f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return err
			}
			defer func() {
				if err := f.Close(); err != nil {
					panic(err)
				}
			}()

			_, err = io.Copy(f, rc)
			if err != nil {
				return err
			}
		}
		return nil
	}

	for _, f := range r.File {
		err := extractAndWriteFile(f)
		if err != nil {
			return err
		}
	}

	return nil
}

func main() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/upload", UploadFile)
	log.Fatal(http.ListenAndServe(":9000", router))
}
