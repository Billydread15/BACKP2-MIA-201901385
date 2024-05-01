package cmds

import (
	"MIA_P1_202004796/objs"
	"MIA_P1_202004796/utilities"
	"encoding/binary"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

func ParseFdisk(entrada string, size *int, driveletter *string, name *string, unit *string, tipe *string, fit *string, delete *string, add *int, path *string) {
	//obtenemos clave-valor de las flags con regex y guardamos en matches
	re := regexp.MustCompile(`-(\w+)=("[^"]+"|\S+)`)
	matches := re.FindAllStringSubmatch(entrada, -1)
	for _, match := range matches {
		flagName := match[1]
		flagValue := match[2]
		// Delete quotes if they are present in the value
		flagValue = strings.Trim(flagValue, "\"")

		switch flagName {
		case "size":
			sizeValue := 0
			sizeInt, _ := strconv.Atoi(flagValue)
			if sizeInt <= 0 {
				fmt.Println("-err size no puede ser negativo")
				return
			}
			fmt.Sscanf(flagValue, "%d", &sizeValue)
			*size = sizeValue
		case "driveletter":
			*driveletter = strings.ToUpper(flagValue[:1])
		case "name":
			*name = flagValue
		case "unit":
			if flagValue != "k" && flagValue != "m" && flagValue != "b" {
				fmt.Println("-err unit no reconocida")
				return
			}
			*unit = flagValue
		case "type":
			if flagValue != "p" && flagValue != "e" && flagValue != "l" {
				fmt.Println("-err type no reconocida")
				return
			}
			*tipe = flagValue
		case "fit":
			flagValue = flagValue[:1]
			if flagValue != "f" && flagValue != "b" && flagValue != "w" {
				fmt.Println("-err fit no reconocida")
				return
			}
			*fit = flagValue
		case "delete":
			if flagValue == "full" {
				*delete = flagValue
			} else {
				*delete = ""
			}
		case "add":
			*add, _ = strconv.Atoi(flagValue)
		}
	}
	Fdisk(*size, *driveletter, *name, *unit, *tipe, *fit, *delete, *add, *path)
}

func Fdisk(size int, driveletter string, name string, unit string, tipe string, fit string, delete string, add int, path string) {
	//definimos size
	if unit == "k" {
		size = size * 1024
	} else if unit == "m" {
		size = size * 1024 * 1024
	}

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

	//si es delete
	if delete == "full" {
		fmt.Println("-delete", name)
		byteSlice := []byte(name)
		var byteArray [16]byte
		copy(byteArray[:], byteSlice[:16])
		//buscamos particion
		for i := 0; i < 4; i++ {
			if reflect.DeepEqual(tmpMbr.Mbr_partitions[i].Part_name[:], byteArray[:]) {
				fmt.Println("-encontrado")
				//borramos
				tmpMbr.Mbr_partitions[i].Part_status = [1]byte{'e'}
				tmpMbr.Mbr_partitions[i].Part_fit = [1]byte{'0'}
				//reescribimos mbr
				if err := utilities.WriteObject(file, tmpMbr, 0); err != nil {
					return
				}
				//si es extendida eliminar las logicas
				if string(tmpMbr.Mbr_partitions[i].Part_type[:]) == "e" {
					fmt.Println("-es extendida")
					if err := utilities.ReplaceWithZeros(file, int64(tmpMbr.Mbr_partitions[i].Part_start), int(tmpMbr.Mbr_partitions[i].Part_s)); err != nil {
						return
					}
				}
				fmt.Println("-particion eliminada con exito")
				return
			}
		}
		return
	}
	//si es add
	if add != 0 {
		//definimos size
		if unit == "k" {
			add = add * 1024
		} else if unit == "m" {
			add = add * 1024 * 1024
		}

		byteSlice := []byte(name)
		var byteArray [16]byte
		var index = 0
		copy(byteArray[:], byteSlice[:16])
		//buscamos particion
		for i := 0; i < 4; i++ {
			if reflect.DeepEqual(tmpMbr.Mbr_partitions[i].Part_name[:], byteArray[:]) {
				index = i
				fmt.Println("-encontrado")
				break
			}
		}
		if add < 0 {
			//restar
			tmpMbr.Mbr_partitions[index].Part_s = tmpMbr.Mbr_partitions[index].Part_s + int32(add)
			//reescribimos mbr
			if err := utilities.WriteObject(file, tmpMbr, 0); err != nil {
				return
			}
			fmt.Println("-resta de memoria realizada con exito")
		} else if add > 0 {
			//sumar
			if tmpMbr.Mbr_partitions[index].Part_s+int32(add)+tmpMbr.Mbr_partitions[index].Part_start > tmpMbr.Mbr_tamano {
				fmt.Println("-err no hay espacio suficiente")
				return
			}
			if tmpMbr.Mbr_partitions[index].Part_s+int32(add)+tmpMbr.Mbr_partitions[index].Part_start > tmpMbr.Mbr_partitions[index+1].Part_start && tmpMbr.Mbr_partitions[index+1].Part_start != 0 {
				fmt.Println("-err no hay espacio suficiente")
				return
			}
			tmpMbr.Mbr_partitions[index].Part_s = tmpMbr.Mbr_partitions[index].Part_s + int32(add)
			//reescribimos mbr
			if err := utilities.WriteObject(file, tmpMbr, 0); err != nil {
				return
			}
			fmt.Println("-suma de memoria realizada con exito")
		}
		return
	}

	//buscamos entre las particiones una vacias
	var count = 0
	var gap = int32(0)
	//iteramos para saber su tamaño
	for i := 0; i < 4; i++ {
		if tmpMbr.Mbr_partitions[i].Part_status == [1]byte{'e'} {
			fmt.Println("-encontramos eliminada")
			break
		}
		if tmpMbr.Mbr_partitions[i].Part_s != 0 {
			count++
			gap = tmpMbr.Mbr_partitions[i].Part_start + tmpMbr.Mbr_partitions[i].Part_s
		}
	}

	for i := 0; i < 4; i++ {
		if tmpMbr.Mbr_partitions[i].Part_s == 0 || tmpMbr.Mbr_partitions[i].Part_status == [1]byte{'e'} {
			if tipe == "p" {
				if tmpMbr.Mbr_partitions[i].Part_s != 0 && tmpMbr.Mbr_partitions[i].Part_s < int32(size) {
					fmt.Println("-err no hay espacio suficiente")
					continue
				}
				//partcion primaria, actualizamos datos
				if count == 0 {
					tmpMbr.Mbr_partitions[i].Part_start = int32(binary.Size(tmpMbr))
				} else {
					tmpMbr.Mbr_partitions[i].Part_start = gap
				}
				tmpMbr.Mbr_partitions[i].Part_s = int32(size)
				tmpMbr.Mbr_partitions[i].Part_correlative = int32(0)
				var tmpName [16]byte
				copy(tmpMbr.Mbr_partitions[i].Part_name[:], tmpName[:])
				copy(tmpMbr.Mbr_partitions[i].Part_name[:], name)
				copy(tmpMbr.Mbr_partitions[i].Part_fit[:], fit)
				copy(tmpMbr.Mbr_partitions[i].Part_type[:], tipe)
				copy(tmpMbr.Mbr_partitions[i].Part_status[:], "0")
				//reescribimos mbr
				if err := utilities.WriteObject(file, tmpMbr, 0); err != nil {
					return
				}

			} else if tipe == "e" {
				if tmpMbr.Mbr_partitions[i].Part_s != 0 && tmpMbr.Mbr_partitions[i].Part_s < int32(size) {
					fmt.Println("-err no hay espacio suficiente")
					continue
				}
				//particion extendida
				//buscar si ya existe alguna extendida
				for i := 0; i < 4; i++ {
					if string(tmpMbr.Mbr_partitions[i].Part_type[:]) == "e" {
						fmt.Println("-err ya existe una particion extendida")
						return
					}
				}
				//llenamos
				if count == 0 {
					tmpMbr.Mbr_partitions[i].Part_start = int32(binary.Size(tmpMbr))
				} else {
					tmpMbr.Mbr_partitions[i].Part_start = gap
				}
				tmpMbr.Mbr_partitions[i].Part_s = int32(size)
				tmpMbr.Mbr_partitions[i].Part_correlative = int32(0)
				copy(tmpMbr.Mbr_partitions[i].Part_name[:], name)
				copy(tmpMbr.Mbr_partitions[i].Part_fit[:], fit)
				copy(tmpMbr.Mbr_partitions[i].Part_type[:], tipe)
				copy(tmpMbr.Mbr_partitions[i].Part_status[:], "0")
				//reescribimos mbr
				if err := utilities.WriteObject(file, tmpMbr, 0); err != nil {
					return
				}
				//creamos ebr
				var tmpEbr objs.EBR
				tmpEbr.Part_next = -1
				tmpEbr.Part_s = 0
				tmpEbr.Part_start = tmpMbr.Mbr_partitions[i].Part_start
				if err := utilities.WriteObject(file, tmpEbr, int64(tmpMbr.Mbr_partitions[i].Part_start)); err != nil {
					return
				}

				var tmpEbr2 objs.EBR
				if err := utilities.ReadObject(file, &tmpEbr2, int64(tmpMbr.Mbr_partitions[i].Part_start)); err != nil {
					return
				}
				objs.PrintEBR(tmpEbr2)

			} else if tipe == "l" {
				//particion logica
				//buscar si ya existe alguna extendida
				var extendida = false
				var indexExtendida = 0
				for i := 0; i < 4; i++ {
					if string(tmpMbr.Mbr_partitions[i].Part_type[:]) == "e" {
						indexExtendida = i
						extendida = true
					}
				}
				if !extendida {
					fmt.Println("-err no existe particion extendida")
					return
				}
				//obtenemos incio de extendida para leer logicas
				var tmpEbr objs.EBR
				start := int64(tmpMbr.Mbr_partitions[indexExtendida].Part_start)
				if err := utilities.ReadObject(file, &tmpEbr, start); err != nil {
					return
				}
				//objs.PrintEBR(tmpEbr)
				//iteramos
				nextlog := tmpEbr.Part_next

				for nextlog != -1 {
					start = int64(tmpEbr.Part_next)
					if err := utilities.ReadObject(file, &tmpEbr, start); err != nil {
						return
					}
					fmt.Println("_encontrado_")
					objs.PrintEBR(tmpEbr)
					nextlog = tmpEbr.Part_next
				}

				//llenamos
				copy(tmpEbr.Part_mount[:], "0")
				copy(tmpEbr.Part_fit[:], fit)
				tmpEbr.Part_s = int32(size)
				tmpEbr.Part_start = int32(start) + int32(binary.Size(tmpEbr))
				tmpEbr.Part_next = int32(start) + int32(binary.Size(tmpEbr)) + int32(size)
				copy(tmpEbr.Part_name[:], name)
				next := int32(start) + int32(binary.Size(tmpEbr)) + int32(size)
				//escribimos
				if err := utilities.WriteObject(file, tmpEbr, start); err != nil {
					return
				}

				if err := utilities.ReadObject(file, &tmpEbr, start); err != nil {
					return
				}
				fmt.Println("_modificado_")
				objs.PrintEBR(tmpEbr)

				//creamos siguiente
				var tmpEbr2 objs.EBR
				tmpEbr2.Part_next = -1
				tmpEbr2.Part_s = 0
				tmpEbr2.Part_start = next
				if err := utilities.WriteObject(file, tmpEbr2, int64(next)); err != nil {
					return
				}
				if err := utilities.ReadObject(file, &tmpEbr2, int64(next)); err != nil {
					return
				}
				fmt.Println("_creado_")
				objs.PrintEBR(tmpEbr2)

				fmt.Println("-particion logica creada con exito")

			}

			break
		}
	}

	var tmpMbr2 objs.MBR
	if err := utilities.ReadObject(file, &tmpMbr2, 0); err != nil {
		return
	}
	objs.PrintMBR(tmpMbr2)
	fmt.Println("-MBR modificado con exito")

}
