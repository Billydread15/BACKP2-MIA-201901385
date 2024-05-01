package cmds

import (
	"MIA_P1_202004796/objs"
	"MIA_P1_202004796/utilities"
	"encoding/binary"
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"
)

func ParseMkfs(entrada string, id *string, tipe *string, fs *string) {
	//obtenemos clave-valor de las flags con regex y guardamos en matches
	re := regexp.MustCompile(`-(\w+)=("[^"]+"|\S+)`)
	matches := re.FindAllStringSubmatch(entrada, -1)
	for _, match := range matches {
		flagName := match[1]
		flagValue := match[2]
		// Delete quotes if they are present in the value
		flagValue = strings.Trim(flagValue, "\"")

		switch flagName {
		case "id":
			if flagValue == "" {
				fmt.Println("-err id no puede ser vacio")
				return
			}
			*id = flagValue
		case "type":
			if flagValue != "full" {
				fmt.Println("-err type no reconocida")
				return
			}
			*tipe = flagValue
		case "fs":
			flagValue = flagValue[:1]
			if flagValue != "2" && flagValue != "3" {
				fmt.Println("-err fs no reconocida")
				return
			}
			*fs = flagValue
		}
	}
	if *id == "" {
		fmt.Println("-err id no puede ser vacio")
		return
	}
	Mkfs(*id, *tipe, *fs)
}

func Mkfs(id string, tipe string, fs string) {
	driveletter := strings.ToUpper(string(id[0]))

	//abrimos dsk
	ruta := "./MIA/P1/" + string(driveletter) + ".dsk"
	fmt.Println("ruta:", ruta)
	file, err := utilities.OpenFile(ruta)
	if err != nil {
		return
	}

	//leemos mbr
	var tmpMbr objs.MBR
	if err := utilities.ReadObject(file, &tmpMbr, 0); err != nil {
		return
	}
	//objs.PrintMBR(tmpMbr)

	//buscamos particion
	var index = -1
	for i := 0; i < 4; i++ {
		if tmpMbr.Mbr_partitions[i].Part_s != 0 {
			if string(tmpMbr.Mbr_partitions[i].Part_id[:]) == strings.ToUpper(id) && string(tmpMbr.Mbr_partitions[i].Part_status[:]) == "1" {
				index = i
				fmt.Println("-encontrado")
				break
			}
		}
	}
	if index == -1 {
		fmt.Println("-err no se encontro la particion")
		return
	} else {
		objs.PrintPartition(tmpMbr.Mbr_partitions[index])
	}

	//calculamos sizes
	numerador := int32(tmpMbr.Mbr_partitions[index].Part_s - int32(binary.Size(objs.Superblock{})))
	denominador_base := int32(4 + int32(binary.Size(objs.Inode{})) + 3*int32(binary.Size(objs.Fileblock{})))
	temp := int32(0)
	if fs == "2" {
		temp = 0
	} else {
		temp = int32(binary.Size(objs.Journaling{}))
	}
	denominador := denominador_base + temp
	n := int32(numerador / denominador)
	fmt.Println("n:", n)

	if fs == "2" {
		Ext2(file, n, tmpMbr.Mbr_partitions[index])
	} else {
		fmt.Println("- fs3")
	}

}

func Ext2(file *os.File, n int32, partition objs.Partition) {
	//creamos superblock
	var newSuper objs.Superblock

	newSuper.S_filesystem_type = 2
	newSuper.S_inodes_count = n
	newSuper.S_blocks_count = 3 * n
	newSuper.S_free_blocks_count = 3 * n
	newSuper.S_free_inodes_count = n

	copy(newSuper.S_mtime[:], time.Now().Format("2006-01-02"))
	copy(newSuper.S_umtime[:], time.Now().Format("2006-01-02"))
	newSuper.S_mnt_count = 1

	newSuper.S_inode_s = int32(binary.Size(objs.Inode{}))
	newSuper.S_block_s = int32(binary.Size(objs.Fileblock{}))

	newSuper.S_bm_inode_start = partition.Part_start + int32(binary.Size(objs.Superblock{}))
	newSuper.S_bm_block_start = newSuper.S_bm_inode_start + n
	newSuper.S_inode_start = newSuper.S_bm_block_start + (3 * n)
	newSuper.S_block_start = newSuper.S_inode_start + (int32(binary.Size(objs.Inode{})) * n)

	//se crearan 2 bloques y 2 inodos por defecto
	newSuper.S_free_blocks_count -= 2
	newSuper.S_free_inodes_count -= 2

	//escribimos superblock
	if err := utilities.WriteObject(file, newSuper, int64(partition.Part_start)); err != nil {
		return
	}

	//escribimos bitmap de inodos
	for i := 0; i < int(n); i++ {
		if err := utilities.WriteObject(file, byte(0), int64(newSuper.S_bm_inode_start)+int64(i)); err != nil {
			return
		}
	}

	//escribimos bitmap de bloques
	for i := 0; i < 3*int(n); i++ {
		if err := utilities.WriteObject(file, byte(0), int64(newSuper.S_bm_block_start)+int64(i)); err != nil {
			return
		}
	}

	//creamos inodos vacios
	var newInode objs.Inode
	for i := 0; i < 15; i++ {
		newInode.I_block[i] = -1
	}
	for i := 0; i < int(n); i++ {
		if err := utilities.WriteObject(file, newInode, int64(newSuper.S_inode_start)+int64(i)*int64(binary.Size(objs.Inode{}))); err != nil {
			return
		}
	}

	//creamos bloques vacios
	var newFblock objs.Fileblock
	for i := 0; i < 3*int(n); i++ {
		if err := utilities.WriteObject(file, newFblock, int64(newSuper.S_block_start)+int64(i)*int64(binary.Size(objs.Fileblock{}))); err != nil {
			return
		}
	}

	//creamos root folder (inode)
	var inode0 objs.Inode

	inode0.I_uid = 1
	inode0.I_gid = 1
	inode0.I_s = 0
	copy(inode0.I_atime[:], time.Now().Format("2006-01-02"))
	copy(inode0.I_ctime[:], time.Now().Format("2006-01-02"))
	copy(inode0.I_mtime[:], time.Now().Format("2006-01-02"))
	copy(inode0.I_perm[:], "0")
	copy(inode0.I_perm[:], "664")
	for i := int32(0); i < 15; i++ {
		inode0.I_block[i] = -1
	}
	inode0.I_block[0] = 0

	//creamos root folder (block)
	var fBlock0 objs.Folderblock

	fBlock0.B_content[0].B_inodo = 0
	copy(fBlock0.B_content[0].B_name[:], ".")

	fBlock0.B_content[1].B_inodo = 0
	copy(fBlock0.B_content[1].B_name[:], "..")

	fBlock0.B_content[2].B_inodo = 1
	copy(fBlock0.B_content[2].B_name[:], "users.txt")

	var inode1 objs.Inode
	inode1.I_uid = 1
	inode1.I_gid = 1
	inode1.I_s = int32(binary.Size(objs.Folderblock{}))
	copy(inode1.I_atime[:], time.Now().Format("2006-01-02"))
	copy(inode1.I_ctime[:], time.Now().Format("2006-01-02"))
	copy(inode1.I_mtime[:], time.Now().Format("2006-01-02"))
	copy(inode1.I_perm[:], "0")
	copy(inode1.I_perm[:], "664")
	for i := int32(0); i < 15; i++ {
		inode1.I_block[i] = -1
	}
	inode1.I_block[0] = 1

	//creamos users.txt en fileblock
	data := "1,G,root\n1,U,root,root,123\n"
	var fBlock1 objs.Fileblock
	copy(fBlock1.B_content[:], data)

	//reescribimos bitmaps
	err := utilities.WriteObject(file, byte(1), int64(newSuper.S_bm_inode_start))
	if err != nil {
		fmt.Println("-err ", err)
		return
	}
	err = utilities.WriteObject(file, byte(1), int64(newSuper.S_bm_inode_start)+1)
	if err != nil {
		fmt.Println("-err ", err)
		return
	}
	err = utilities.WriteObject(file, byte(1), int64(newSuper.S_bm_block_start))
	if err != nil {
		fmt.Println("-err ", err)
		return
	}
	err = utilities.WriteObject(file, byte(1), int64(newSuper.S_bm_block_start)+1)
	if err != nil {
		fmt.Println("-err ", err)
		return
	}

	//reescribimos inodos
	err = utilities.WriteObject(file, inode0, int64(newSuper.S_inode_start))
	if err != nil {
		fmt.Println("-err ", err)
		return
	}
	err = utilities.WriteObject(file, inode1, int64(newSuper.S_inode_start)+int64(binary.Size(objs.Inode{})))
	if err != nil {
		fmt.Println("-err ", err)
		return
	}

	//reescribimos bloques
	err = utilities.WriteObject(file, fBlock0, int64(newSuper.S_block_start))
	if err != nil {
		fmt.Println("-err ", err)
		return
	}

	err = utilities.WriteObject(file, fBlock1, int64(newSuper.S_block_start)+int64(binary.Size(objs.Fileblock{})))
	if err != nil {
		fmt.Println("-err ", err)
		return
	}
	fmt.Println("-formateado")
}
