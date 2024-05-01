package cmds

import (
	"MIA_P1_202004796/utilities"
	"fmt"
	"regexp"
	"strings"
)

func ParseRmdisk(entrada string, driveletter *string) {
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
		default:
			fmt.Println("-err flag no reconocida")
		}
	}
}

func Rmdisk(driveletter string) {
	//eliminar archivo
	path := "./MIA/P1/" + driveletter + ".dsk"
	utilities.DeleteFile(path)
}
