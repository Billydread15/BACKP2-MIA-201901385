package utilities

import (
	"MIA_P1_202004796/objs"
	"encoding/binary"
	"fmt"
	"os"
	"strings"
)

// login -user=root -pass=123 -id=A119
func InitSearch(path string, file *os.File, tempSuperblock objs.Superblock) int32 {
	fmt.Println("Start INITSEARCH")
	fmt.Println("path:", path)
	// path = "/ruta/nueva"

	// split the path by /
	TempStepsPath := strings.Split(path, "/")
	StepsPath := TempStepsPath[1:]

	fmt.Println("StepsPath:", StepsPath, "len(StepsPath):", len(StepsPath))
	for _, step := range StepsPath {
		fmt.Println("step:", step)
	}

	var Inode0 objs.Inode
	// Read object from bin file
	if err := ReadObject(file, &Inode0, int64(tempSuperblock.S_inode_start)); err != nil {
		return -1
	}

	return SarchInodeByPath(StepsPath, Inode0, file, tempSuperblock)
}

func pop(s *[]string) string {
	lastIndex := len(*s) - 1
	last := (*s)[lastIndex]
	*s = (*s)[:lastIndex]
	return last
}

// login -user=root -pass=123 -id=A119
func SarchInodeByPath(StepsPath []string, Inode objs.Inode, file *os.File, tempSuperblock objs.Superblock) int32 {
	fmt.Println("Start SARCHINODEBYPATH")
	index := int32(0)
	SearchedName := strings.Replace(pop(&StepsPath), " ", "", -1)

	fmt.Println("========== SearchedName:", SearchedName)

	// Iterate over i_blocks from Inode
	for _, block := range Inode.I_block {
		if block != -1 {
			if index < 13 {
				//CASO DIRECTO

				var crrFolderBlock objs.Folderblock
				// Read object from bin file
				if err := ReadObject(file, &crrFolderBlock, int64(tempSuperblock.S_block_start+block*int32(binary.Size(objs.Folderblock{})))); err != nil {
					return -1
				}

				for _, folder := range crrFolderBlock.B_content {
					// fmt.Println("Folder found======")
					fmt.Println("Folder === Name:", string(folder.B_name[:]), "B_inodo", folder.B_inodo)

					if strings.Contains(string(folder.B_name[:]), SearchedName) {

						fmt.Println("len(StepsPath)", len(StepsPath), "StepsPath", StepsPath)
						if len(StepsPath) == 0 {
							fmt.Println("Folder found======")
							return folder.B_inodo
						} else {
							fmt.Println("NextInode======")
							var NextInode objs.Inode
							// Read object from bin file
							if err := ReadObject(file, &NextInode, int64(tempSuperblock.S_inode_start+folder.B_inodo*int32(binary.Size(objs.Inode{})))); err != nil {
								return -1
							}
							return SarchInodeByPath(StepsPath, NextInode, file, tempSuperblock)
						}
					}
				}

			}
		}
		index++
	}

	fmt.Println("End SARCHINODEBYPATH")
	return 0
}

func GetInodeFileData(Inode objs.Inode, file *os.File, tempSuperblock objs.Superblock) string {
	fmt.Println("Start GETINODEFILEDATA")
	index := int32(0)
	// define content as a string
	var content string

	// Iterate over i_blocks from Inode
	for _, block := range Inode.I_block {
		if block != -1 {
			if index < 13 {
				//CASO DIRECTO

				var crrFileBlock objs.Fileblock
				// Read object from bin file
				if err := ReadObject(file, &crrFileBlock, int64(tempSuperblock.S_block_start+block*int32(binary.Size(objs.Fileblock{})))); err != nil {
					return ""
				}

				content += string(crrFileBlock.B_content[:])

			}
		}
		index++
	}

	fmt.Println("End GETINODEFILEDATA")
	return content
}
