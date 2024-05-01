package cmds

import (
	"MIA_P1_202004796/objs"
	"MIA_P1_202004796/utilities"
	"fmt"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func ParseMkdisk(entrada string, size *int, fit *string, unit *string) {
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
		case "fit":
			flagValue = flagValue[:1]
			if flagValue != "f" && flagValue != "b" && flagValue != "w" {
				fmt.Println("-err fit no reconocida")
				return
			}
			*fit = flagValue
		case "unit":
			if flagValue != "k" && flagValue != "m" {
				fmt.Println("-err unit no reconocida")
				return
			}
			*unit = flagValue
		default:
			fmt.Println("-err flag no reconocida")
		}
	}
	fmt.Println("size:", *size)
	fmt.Println("fit:", *fit)
	fmt.Println("unit:", *unit)
}

func Mkdisk(size int, fit string, unit string) {
	//crear archivo
	name, err := utilities.CreateFile("./MIA/P1/", 65)
	if err != nil {
		fmt.Println("-err ", err)
	}
	path := "./MIA/P1/" + string(name) + ".dsk"

	//definir size
	if unit == "k" {
		size = size * 1024
	} else {
		size = size * 1024 * 1024
	}

	//abrir archivo
	file, _ := utilities.OpenFile(path)

	//escribimos 0 binarios
	arreglo := make([]byte, 1024)
	for i := 0; i <= size/1024; i++ {
		err := utilities.WriteObject(file, arreglo, int64(i*1024))
		if err != nil {
			fmt.Println("-err ", err)
		}
	}

	//escribimos mbr
	var newMBR objs.MBR
	newMBR.Mbr_tamano = int32(size)
	copy(newMBR.Mbr_fecha_creacion[:], time.Now().Format("2006-01-02"))
	newMBR.Mbr_dsk_signature = rand.Int31()
	copy(newMBR.Mbr_dsk_fit[:], fit)

	if err := utilities.WriteObject(file, newMBR, 0); err != nil {
		return
	}

	// Close bin file
	defer file.Close()
	println("-Disco creado con exito")
}
