package cmds

import (
	"MIA_P1_202004796/objs"
	"MIA_P1_202004796/utilities"
	"fmt"
	"regexp"
	"strings"
)

func ParseUnmount(entrada string, id *string) {
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
		}
	}

	Unmount(*id)
}

func Unmount(id string) {
	driveletter := id[0]
	driveletter = byte(strings.ToUpper(string(driveletter))[0])
	fmt.Println("driveletter:", string(driveletter))
	//abrimos dsk
	ruta := "./MIA/P1/" + string(driveletter) + ".dsk"
	//fmt.Println("ruta:", ruta)
	file, err := utilities.OpenFile(ruta)
	if err != nil {
		return
	}

	//leemos mbr
	var tmpMbr objs.MBR
	if err := utilities.ReadObject(file, &tmpMbr, 0); err != nil {
		return
	}

	//buscamos particion
	for i := 0; i < 4; i++ {

		if string(tmpMbr.Mbr_partitions[i].Part_id[:]) == strings.ToUpper(id) {
			//particion extendida
			if string(tmpMbr.Mbr_partitions[i].Part_type[:]) == "e" {
				fmt.Println("-err no se puede desmontar una particion extendida")
				return
			}
			//desmontamos
			copy(tmpMbr.Mbr_partitions[i].Part_status[:], "0")
			tmpMbr.Mbr_partitions[i].Part_id = [4]byte{}
			tmpMbr.Mbr_partitions[i].Part_correlative = int32(0)
			utilities.WriteObject(file, tmpMbr, 0)
			fmt.Println("-particion desmontada")
			return
		}
	}

}
