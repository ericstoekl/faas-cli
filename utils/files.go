package utils

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
)

// CopyFiles copies files from src to destination, optionally recursively.
func CopyFiles(src string, destination string, recursive bool) {

	files, err := ioutil.ReadDir(src)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {

		if file.IsDir() == false {

			cp(src+"/"+file.Name(), destination+file.Name())

		} else {
			//make new destination dir
			newDir := destination + file.Name() + "/"

			if !pathExists(newDir) {

				debugPrint(fmt.Sprintf("Creating directory: %s at %s", file.Name(), newDir))
				newDirErr := os.Mkdir(newDir, 0700)

				if newDirErr != nil {
					fmt.Printf("Error creating path %s - %s.\n", newDir, newDirErr.Error())
				}
			}

			//did the call ask to recurse into sub directories?
			if recursive == true {
				//call CopyFiles to copy the contents
				CopyFiles(src+"/"+file.Name(), newDir, true)
			}
		}
	}
}

func pathExists(path string) bool {
	exists := true

	if _, err := os.Stat(path); os.IsNotExist(err) {
		exists = false
	}

	return exists
}

func cp(src string, destination string) error {

	debugPrint(fmt.Sprintf("cp - %s %s", src, destination))

	memoryBuffer, readErr := ioutil.ReadFile(src)
	if readErr != nil {
		return fmt.Errorf("Error reading source file: %s\n" + readErr.Error())
	}
	writeErr := ioutil.WriteFile(destination, memoryBuffer, 0660)
	if writeErr != nil {
		return fmt.Errorf("Error writing file: %s\n" + writeErr.Error())
	}

	return nil
}

func debugPrint(message string) {

	if val, exists := os.LookupEnv("debug"); exists && (val == "1" || val == "true") {
		fmt.Println(message)
	}
}

// Copy the src file to dst. Any existing file will be overwritten and will not
// copy file attributes.
func Copy(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}
	return out.Close()
}
