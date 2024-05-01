package utilities

import (
	"encoding/binary"
	"fmt"
	"os"
	"path/filepath"
)

// Crear file bin
func CreateFile(path string, name int32) (int32, error) {
	//Ensure the directory exists
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		fmt.Println("-err CreateFile dir==", err)
		return 0, err
	}

	name = rune(name)
	file := filepath.Join(dir, string(name)+".dsk")

	// Create file
	if _, err := os.Stat(file); os.IsNotExist(err) {
		file, err := os.Create(file)
		if err != nil {
			fmt.Println("-err CreateFile create==", err)
			return 0, err
		}
		defer file.Close()
		return name, nil

	} else {
		name, _ = CreateFile(path, name+1)
	}
	return name, nil
}

// Funtion to open bin file in read/write mode
func OpenFile(name string) (*os.File, error) {
	file, err := os.OpenFile(name, os.O_RDWR, 0644)
	if err != nil {
		fmt.Println("-err OpenFile==", err)
		return nil, err
	}
	return file, nil
}

// Function to Write an object in a bin file
func WriteObject(file *os.File, data interface{}, position int64) error {
	file.Seek(position, 0)
	err := binary.Write(file, binary.LittleEndian, data)
	if err != nil {
		fmt.Println("-err WriteObject==", err)
		return err
	}
	return nil
}

// Function to replace data with 0s
func ReplaceWithZeros(file *os.File, position int64, size int) error {
	file.Seek(position, 0)
	for i := 0; i < size; i++ {
		err := binary.Write(file, binary.LittleEndian, byte(0))
		if err != nil {
			fmt.Println("-err ReplaceWithZeros", err)
			return err
		}
	}
	return nil
}

// Function to Read an object from a bin file
func ReadObject(file *os.File, data interface{}, position int64) error {
	file.Seek(position, 0)
	err := binary.Read(file, binary.LittleEndian, data)
	if err != nil {
		fmt.Println("Err ReadObject==", err)
		return err
	}
	return nil
}

// Function to delete a file
func DeleteFile(name string) error {
	println("Seguro que quiere remover ", name, "? (y/n)")
	var input string
	fmt.Scanln(&input)
	if input == "y" {
		err := os.Remove(name)
		if err != nil {
			fmt.Println("-err DeleteFile==", err)
			return err
		}
		println("-Archivo eliminado")
	}
	return nil
}
