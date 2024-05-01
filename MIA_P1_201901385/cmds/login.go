package cmds

import (
	"MIA_P1_202004796/global"
	"MIA_P1_202004796/objs"
	"MIA_P1_202004796/utilities"
	"encoding/binary"
	"fmt"
	"regexp"
	"strings"
)

func ParseLogin(entrada string, user *string, password *string, id *string) {

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
		case "pass":
			if flagValue == "" {
				fmt.Println("-err password no puede ser vacio")
				return
			}
			*password = flagValue
		case "user":
			if flagValue == "" {
				fmt.Println("-err user no puede ser vacio")
				return
			}
			*user = flagValue
		}
	}
	/*
		if *id == "" {
			fmt.Println("-err id no puede ser vacio")
			return
		} else if *user == "" {
			fmt.Println("-err user no puede ser vacio")
			return
		} else if *password == "" {
			fmt.Println("-err password no puede ser vacio")
			return
		}*/

	Login(*user, *password, *id)
}

func Login(user string, password string, id string) string {
	fmt.Println("user: ", user, " password: ", password, " id: ", id)

	if global.Usuario.Status {
		fmt.Println("-err ya hay un usuario logueado")
		return "-err ya hay un usuario logueado"
	}

	//obtenemos disco y particion del user
	var login bool = false
	driveletter := string(id[0])
	driveletter = strings.ToUpper(driveletter)

	//abrimos dsk
	ruta := "./MIA/P1/" + string(driveletter) + ".dsk"
	fmt.Println("ruta:", ruta)
	file, err := utilities.OpenFile(ruta)
	if err != nil {
		return ""
	}

	//leemos mbr
	var tmpMbr objs.MBR
	if err := utilities.ReadObject(file, &tmpMbr, 0); err != nil {
		return ""
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
		return "-err no se encontro la particion o no está montada"
	} else {
		objs.PrintPartition(tmpMbr.Mbr_partitions[index])
	}

	var tmpSuperBlock objs.Superblock
	if err := utilities.ReadObject(file, &tmpSuperBlock, int64(tmpMbr.Mbr_partitions[index].Part_start)); err != nil {
		return ""
	}
	indexInode := utilities.InitSearch("/users.txt", file, tmpSuperBlock)
	// indexInode := int32(1)

	var crrInode objs.Inode
	// Read object from bin file
	if err := utilities.ReadObject(file, &crrInode, int64(tmpSuperBlock.S_inode_start+indexInode*int32(binary.Size(objs.Inode{})))); err != nil {
		return ""
	}

	// read file data
	data := utilities.GetInodeFileData(crrInode, file, tmpSuperBlock)

	//fmt.Println("Fileblock------------")
	// Dividir la cadena en líneas
	lines := strings.Split(data, "\n")

	// login -user=root -pass=123 -id=A119

	// Iterar a través de las líneas
	for _, line := range lines {
		// Imprimir cada línea
		//fmt.Println("CONTENT----", line)
		words := strings.Split(line, ",")

		if len(words) == 5 {
			if (strings.Contains(words[3], user)) && (strings.Contains(words[4], password)) {
				login = true
				break
			}
		}
	}

	// Print object
	//fmt.Println("Inode", crrInode.I_block)

	// Close bin file
	defer file.Close()

	if login {
		fmt.Println("-user logged in")
		global.Usuario.Id = id
		global.Usuario.Status = true
		global.Usuario.Username = user
		return "-user logged in"
	}
	return "-err credenciales incorrectas"

}

func Logout() string {
	if global.Usuario.Status {
		global.Usuario.Id = ""
		global.Usuario.Status = false
		global.Usuario.Username = ""
		fmt.Println("-user logged out")
		return "-user logged out"
	} else {
		fmt.Println("-err no user logged in")
		return "-err no user logged in"
	}
}
