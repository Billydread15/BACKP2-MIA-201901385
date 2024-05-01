package cmds

import (
	"MIA_P1_202004796/objs"
	"MIA_P1_202004796/utilities"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

func ParseMount(entrada string, driveletter *string, name *string) {
	//obtenemos clave-valor de las flags con regex y guardamos en matches
	re := regexp.MustCompile(`-(\w+)=("[^"]+"|\S+)`)
	matches := re.FindAllStringSubmatch(entrada, -1)
	for _, match := range matches {
		flagName := match[1]
		flagValue := match[2]
		// Delete quotes if they are present in the value
		flagValue = strings.Trim(flagValue, "\"")

		switch flagName {
		case "driveletter":
			*driveletter = strings.ToUpper(flagValue[:1])
		case "name":
			*name = flagValue
		}
	}

	Mount(*driveletter, *name)
}

func Mount(driveletter string, name string) {
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

	//buscamos el correlativo
	var count = 1
	//iteramos para saber su tamaño
	for i := 0; i < 4; i++ {
		if tmpMbr.Mbr_partitions[i].Part_correlative != 0 {
			count++
		}
	}

	//buscamos particion
	byteSlice := []byte(name)
	var byteArray [16]byte
	copy(byteArray[:], byteSlice[:16])
	for i := 0; i < 4; i++ {

		if reflect.DeepEqual(tmpMbr.Mbr_partitions[i].Part_name[:], byteArray[:]) {
			//particion extendida
			if string(tmpMbr.Mbr_partitions[i].Part_type[:]) == "e" {
				fmt.Println("-err no se puede montar una particion extendida")
				return
			}

			//modificamos valores
			tmpMbr.Mbr_partitions[i].Part_correlative = int32(count)
			id := driveletter + strconv.Itoa(int(tmpMbr.Mbr_partitions[i].Part_correlative)) + "96"
			copy(tmpMbr.Mbr_partitions[i].Part_status[:], []byte("1"))
			copy(tmpMbr.Mbr_partitions[i].Part_id[:], []byte(id)[:4])

			fmt.Println("id:", id)

			//escribimos mbr
			if err := utilities.WriteObject(file, tmpMbr, 0); err != nil {
				return
			}

			//leemos mbr
			var tmpMbr2 objs.MBR
			if err := utilities.ReadObject(file, &tmpMbr2, 0); err != nil {
				return
			}
			objs.PrintMBR(tmpMbr2)
			fmt.Println("-particion ", name, "montada")
			break
		}

		//buscamos en la particion extendida
		if string(tmpMbr.Mbr_partitions[i].Part_type[:]) == "e" {
			//leemos ebr
			var tmpEbr objs.EBR
			start := int64(tmpMbr.Mbr_partitions[i].Part_start)
			if err := utilities.ReadObject(file, &tmpEbr, start); err != nil {
				return
			}
			//buscamos particion
			nextlog := tmpEbr.Part_next
			for nextlog != -1 {
				if reflect.DeepEqual(tmpEbr.Part_name[:], byteArray[:]) {
					fmt.Println("-encontrado")
					//modificamos valores
					tmpEbr.Part_mount = [1]byte{'1'}
					//reescribimos ebr
					if err := utilities.WriteObject(file, tmpEbr, start); err != nil {
						return
					}
					fmt.Println("-particion ", name, "montada")
				}
				start = int64(tmpEbr.Part_next)
				if err := utilities.ReadObject(file, &tmpEbr, start); err != nil {
					return
				}
				//fmt.Println("_encontrado_")
				//objs.PrintEBR(tmpEbr)
				nextlog = tmpEbr.Part_next
			}
		}
	}
}
