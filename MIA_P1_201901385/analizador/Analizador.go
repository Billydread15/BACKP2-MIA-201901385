package Analizador

import (
	"MIA_P1_202004796/cmds"
	"flag"
	"fmt"
	"strings"
)

func Analizar(entrada string) {
	//declaramos todas las flags para los cmds
	unit := flag.String("unit", "m", "unit of memory")
	fit := flag.String("fit", "f", "fit of disk")
	size := flag.Int("size", 0, "size of disk")
	driveletter := flag.String("driveletter", "a", "disk to erase")
	name := flag.String("name", "name", "name of partition")
	tipe := flag.String("type", "p", "type of partition")
	delete := flag.String("delete", "", "delete of partition")
	add := flag.Int("add", 0, "add of partition")
	path := flag.String("path", "", "path of file")
	id := flag.String("id", "", "id of partition")
	fs := flag.String("fs", "2", "file system of partition")
	ruta := flag.String("ruta", "", "path of file")
	user := flag.String("user", "", "user of partition")
	pass := flag.String("pass", "", "password of partition")
	flag.Parse()

	/*var input string
	fmt.Println("Enter command: ")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	input = scanner.Text()
	_, params := getCommandAndParams(input)
	pathE := strings.Split(params, "=")

	//open and read the file
	file, err := os.Open(pathE[1])
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()
	scanner2 := bufio.NewScanner(file)*/

	for _, line := range strings.Split(entrada, "\n") {
		// Process each line here
		line = strings.ReplaceAll(line, "'", "")
		line = strings.ReplaceAll(line, "\"", "")

		if strings.HasPrefix(line, "#") || line == "" {
			continue
		}
		fmt.Println("\n[CMD]", line)

		//TODO HERE!!!!!
		//parseamos a minusculas
		line = strings.ToLower(line)
		//obtenemos el comando
		cmd := strings.Split(line, " ")
		fmt.Println(cmd[0])

		if cmd[0] == "mkdisk" {
			cmds.ParseMkdisk(line, size, fit, unit)
			cmds.Mkdisk(*size, *fit, *unit)
		} else if cmd[0] == "rmdisk" {
			cmds.ParseRmdisk(line, driveletter)
			cmds.Rmdisk(*driveletter)
		} else if cmd[0] == "fdisk" {
			flag.Set("unit", "k")
			flag.Set("fit", "w")
			cmds.ParseFdisk(line, size, driveletter, name, unit, tipe, fit, delete, add, path)
		} else if cmd[0] == "mount" {
			cmds.ParseMount(line, driveletter, name)
		} else if cmd[0] == "unmount" {
			cmds.ParseUnmount(line, id)
		} else if cmd[0] == "mkfs" {
			flag.Set("tipe", "full")
			cmds.ParseMkfs(line, id, tipe, fs)
		} else if cmd[0] == "rep" {
			cmds.ParseRep(line, name, path, id, ruta)
		} else if cmd[0] == "login" {
			cmds.ParseLogin(line, user, pass, id)
		} else if cmd[0] == "logout" {
			cmds.Logout()
		}

	}

	flag.CommandLine = flag.NewFlagSet("", flag.ExitOnError)
}

/*func getCommandAndParams(input string) (string, string) {
	parts := strings.Fields(input)
	if len(parts) > 0 {
		command := strings.ToLower(parts[0])
		params := strings.Join(parts[1:], " ")
		return command, params
	}
	return "", input
}*/
