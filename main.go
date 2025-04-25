package main

import (
	"crypto/sha1"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: myvcs <command>")
		return
	}

	command := os.Args[1]

	switch command {
	case "init":
		err := initRepo()
		if err != nil {
			fmt.Println("Error initializing the repo : ", err)
		} else {
			fmt.Println("Initialized the repo successfully!")
		}
	case "hash-object":
		if len(os.Args) < 3 {
			fmt.Println("Command usage : myvcs hash-object <filename>")
			return
		}
		filename := os.Args[2]
		hash, err := hashObject(filename)
		if err != nil {
			fmt.Println("Error creating the hash", err)
			return
		}
		fmt.Println(hash)
	case "cat-file":
		if len(os.Args) < 3 {
			fmt.Println("Command usage : myvcs cat-file <hash>")
			return
		}
		hash := os.Args[2]
		content, err := catFile(hash)
		if err != nil {
			fmt.Println("Error: ", err)
			return
		}
		fmt.Print(string(content))
	default:
		fmt.Println("Unknown command: ", command)
	}
}

func initRepo() error {
	repoDir := ".myvcs"

	subDirs := []string{"objects", "refs"}

	err := os.Mkdir(repoDir, 0755)
	if err != nil {
		fmt.Println(err)
	}

	for _, dir := range subDirs {
		err = os.MkdirAll(repoDir+"/"+dir, 0755) //it creates all necessary parent directories if they don't exist.
		if err != nil {
			fmt.Println(err)
		}
	}

	return nil
}

func hashObject(filename string) (string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hasher := sha1.New()

	_, err = io.Copy(hasher, file)
	if err != nil {
		return "", err
	}

	hashSum := hasher.Sum(nil)

	//convert hash to hex value
	hashString := fmt.Sprintf("%x", hashSum)

	_, err = file.Seek(0, 0)
	if err != nil {
		return "", err
	}

	content, err := io.ReadAll(file)
	if err != nil {
		return "", err
	}

	objectPath := filepath.Join(".myvcs", "objects", hashString)
	err = os.WriteFile(objectPath, content, 0644)
	if err != nil {
		return "", err
	}

	return hashString, nil

}

func catFile(hash string) ([]byte, error) {
	objectPath := filepath.Join(".myvcs", "objects", hash)

	if len(hash) != 40 {
		return nil, fmt.Errorf("invalid hash provided")
	}

	content, err := os.ReadFile(objectPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("path does not exist for given hash")
		}
		return nil, err
	}
	return content, nil
}
