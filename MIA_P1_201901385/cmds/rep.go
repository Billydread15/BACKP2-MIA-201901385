package cmds

import (
	"MIA_P1_202004796/objs"
	"MIA_P1_202004796/utilities"
	"bytes"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func ParseRep(entrada string, name *string, path *string, id *string, ruta *string) {
	//obtenemos clave-valor de las flags con regex y guardamos en matches
	re := regexp.MustCompile(`-(\w+)=("[^"]+"|\S+)`)
	matches := re.FindAllStringSubmatch(entrada, -1)
	for _, match := range matches {
		flagName := match[1]
		flagValue := match[2]
		// Delete quotes if they are present in the value
		flagValue = strings.Trim(flagValue, "\"")
		switch flagName {
		case "name":
			if flagValue == "" {
				fmt.Println("-err name no puede ser vacio")
				return
			}
			*name = flagValue
		case "path":
			if flagValue == "" {
				fmt.Println("-err path no puede ser vacio")
				return
			}
			*path = flagValue
		case "id":
			if flagValue == "" {
				fmt.Println("-err id no puede ser vacio")
				return
			}
			*id = flagValue
		case "ruta":
			if flagValue == "" {
				fmt.Println("-err ruta no puede ser vacio")
				return
			}
			*ruta = flagValue
		}
	}
	if *name == "" {
		fmt.Println("-err name no puede ser vacio")
		return
	}
	if *path == "" {
		fmt.Println("-err path no puede ser vacio")
		return
	}
	if *id == "" {
		fmt.Println("-err id no puede ser vacio")
		return
	}

	switch *name {
	case "mbr":
		RepMBR(*id, *path)
	case "disk":
		Repdsk(*id, *path)
	case "bm_inode":
		RepBMInode(*id, *path)
	case "bm_block":
		RepBMBlock(*id, *path)
	case "sb":
		RepSB(*id, *path)
	}
}

func RepMBR(id string, path string) {
	driveletter := id[0]
	driveletter = byte(strings.ToUpper(string(driveletter))[0])
	fmt.Println("driveletter:", string(driveletter))
	//abrimos dsk
	rutas := "./MIA/P1/" + string(driveletter) + ".dsk"
	//fmt.Println("ruta:", ruta)
	file, err := utilities.OpenFile(rutas)
	if err != nil {
		return
	}

	//leemos mbr
	var tmpMbr objs.MBR
	if err := utilities.ReadObject(file, &tmpMbr, 0); err != nil {
		return
	}
	objs.PrintMBR(tmpMbr)

	dotContent := `digraph G {node [shape=plaintext]node1 [label=<<TABLE BORDER="1" CELLBORDER="1" CELLSPACING="0"><TR><TD bgcolor="#5364ed" COLSPAN="2">Reporte MBR</TD></TR>`
	dotContent += `	<TR> <TD >mbr_tamano</TD> <TD>` + fmt.Sprint(tmpMbr.Mbr_tamano) + `</TD> </TR>`
	dotContent += `	<TR> <TD >mbr_fecha_creacion</TD> <TD>` + string(tmpMbr.Mbr_fecha_creacion[:]) + `</TD> </TR>`
	dotContent += `	<TR> <TD >mbr_dsk_signature</TD> <TD>` + fmt.Sprint(tmpMbr.Mbr_dsk_signature) + `</TD> </TR>`
	for i := 0; i < 4; i++ {
		if tmpMbr.Mbr_partitions[i].Part_start == 0 {
			continue
		}
		dotContent += `	<TR><TD bgcolor="#5364ed" COLSPAN="2">Particion ` + fmt.Sprint(i) + `</TD></TR>`
		dotContent += `	<TR> <TD >Part_status</TD> <TD>` + string(tmpMbr.Mbr_partitions[i].Part_status[:]) + `</TD> </TR>`
		dotContent += `	<TR> <TD >Part_type</TD> <TD>` + string(tmpMbr.Mbr_partitions[i].Part_type[:]) + `</TD> </TR>`
		dotContent += `	<TR> <TD >Part_fit</TD> <TD>` + string(tmpMbr.Mbr_partitions[i].Part_fit[:]) + `</TD> </TR>`
		dotContent += `	<TR> <TD >Part_start</TD> <TD>` + fmt.Sprint(tmpMbr.Mbr_partitions[i].Part_start) + `</TD> </TR>`
		dotContent += `	<TR> <TD >Part_s</TD> <TD>` + fmt.Sprint(tmpMbr.Mbr_partitions[i].Part_s) + `</TD> </TR>`
		partName := string(tmpMbr.Mbr_partitions[i].Part_name[:])
		nullIndex := bytes.IndexByte(tmpMbr.Mbr_partitions[i].Part_name[:], 0)
		partName = partName[:nullIndex] // Recorta la cadena en el byte nulo
		dotContent += `	<TR> <TD >Part_name</TD> <TD>` + string(partName) + `</TD> </TR>`
		if string(tmpMbr.Mbr_partitions[i].Part_type[:]) == "e" {
			var tmpEbr objs.EBR
			start := int64(tmpMbr.Mbr_partitions[i].Part_start)
			if err := utilities.ReadObject(file, &tmpEbr, start); err != nil {
				return
			}
			//buscamos particion
			nextlog := tmpEbr.Part_next
			for nextlog != -1 {
				dotContent += `	<TR><TD bgcolor="#8f53ed" COLSPAN="2">EBR</TD></TR>`
				dotContent += `	<TR> <TD >Part_status</TD> <TD>` + string(tmpEbr.Part_mount[:]) + `</TD> </TR>`
				dotContent += `	<TR> <TD >Part_fit</TD> <TD>` + string(tmpEbr.Part_fit[:]) + `</TD> </TR>`
				dotContent += `	<TR> <TD >Part_start</TD> <TD>` + fmt.Sprint(tmpEbr.Part_start) + `</TD> </TR>`
				dotContent += `	<TR> <TD >Part_s</TD> <TD>` + fmt.Sprint(tmpEbr.Part_s) + `</TD> </TR>`
				dotContent += `	<TR> <TD >Part_next</TD> <TD>` + fmt.Sprint(tmpEbr.Part_next) + `</TD> </TR>`
				partName := string(tmpMbr.Mbr_partitions[i].Part_name[:])
				nullIndex := bytes.IndexByte(tmpMbr.Mbr_partitions[i].Part_name[:], 0)
				partName = partName[:nullIndex] // Recorta la cadena en el byte nulo
				dotContent += `	<TR> <TD >Part_name</TD> <TD>` + string(partName) + `</TD> </TR>`
				start = int64(tmpEbr.Part_next)
				if err := utilities.ReadObject(file, &tmpEbr, start); err != nil {
					return
				}
				nextlog = tmpEbr.Part_next
			}
		}
	}
	dotContent += `
			</TABLE>
		>]
	}
	`
	fmt.Println(dotContent)
	generate(dotContent, path)
}

func Repdsk(id string, path string) {
	driveletter := id[0]
	driveletter = byte(strings.ToUpper(string(driveletter))[0])
	fmt.Println("driveletter:", string(driveletter))
	//abrimos dsk
	rutas := "./MIA/P1/" + string(driveletter) + ".dsk"
	//fmt.Println("ruta:", ruta)
	file, err := utilities.OpenFile(rutas)
	if err != nil {
		return
	}

	//leemos mbr
	var tmpMbr objs.MBR
	if err := utilities.ReadObject(file, &tmpMbr, 0); err != nil {
		return
	}
	sizeDsk := tmpMbr.Mbr_tamano
	used := float64(0)
	totalused := float64(0)
	dotContent := `digraph D {subgraph cluster_0 {bgcolor="#68d9e2"node [style="rounded" style=filled];node_A [shape=record    label=`
	label := "MBR|"
	for i := 0; i < 4; i++ {
		if tmpMbr.Mbr_partitions[i].Part_start == 0 {
			continue
		}

		if string(tmpMbr.Mbr_partitions[i].Part_status[:]) == "e" {
			used = float64(tmpMbr.Mbr_partitions[i].Part_s) * 100 / float64(sizeDsk)
			totalused += used
			fmt.Println("errased:", used)
			label += "{Libre|{" + strconv.FormatFloat(used, 'f', 2, 64) + "%}}|"
			continue
		}

		if string(tmpMbr.Mbr_partitions[i].Part_type[:]) == "p" {
			used = float64(tmpMbr.Mbr_partitions[i].Part_s) * 100 / float64(sizeDsk)
			totalused += used
			fmt.Println("primary:", used)

			label += "{Primaria|{" + strconv.FormatFloat(used, 'f', 2, 64) + "%}}|"
		} else if string(tmpMbr.Mbr_partitions[i].Part_type[:]) == "e" {
			label += "{Extendida|{"
			var tmpEbr objs.EBR
			start := int64(tmpMbr.Mbr_partitions[i].Part_start)
			if err := utilities.ReadObject(file, &tmpEbr, start); err != nil {
				return
			}
			//buscamos particion
			nextlog := tmpEbr.Part_next
			for nextlog != -1 {
				used = float64(tmpEbr.Part_s) * 100 / float64(sizeDsk)
				totalused += used
				fmt.Println("-logic:", used)

				label += "EBR|{LOGICA|{" + strconv.FormatFloat(used, 'f', 2, 64) + "%}}|"
				start = int64(tmpEbr.Part_next)
				if err := utilities.ReadObject(file, &tmpEbr, start); err != nil {
					return
				}
				nextlog = tmpEbr.Part_next
			}
			label += "}}|"
		}
	}
	free := 100 - totalused
	label += "{Libre|{" + strconv.FormatFloat(free, 'f', -1, 64) + "%}}"

	dotContent += "\"" + label + "\""
	dotContent += `];}}`

	fmt.Println(dotContent)
	generate(dotContent, path)
}

func RepBMInode(id string, path string) {
	name := strings.Split(path, "/")
	file_name := name[len(name)-1]
	file_name = strings.Split(file_name, ".")[0]
	fmt.Println(file_name)

	driveletter := id[0]
	driveletter = byte(strings.ToUpper(string(driveletter))[0])
	fmt.Println("driveletter:", string(driveletter))
	//abrimos dsk
	rutas := "./MIA/P1/" + string(driveletter) + ".dsk"
	//fmt.Println("ruta:", ruta)
	file, err := utilities.OpenFile(rutas)
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
			//leemos superblock
			var tmpSuper objs.Superblock
			start := int64(tmpMbr.Mbr_partitions[i].Part_start)
			if err := utilities.ReadObject(file, &tmpSuper, start); err != nil {
				return
			}
			//leemos bitmap de inodos
			inodes := tmpSuper.S_inodes_count
			bm_inodes_start := tmpSuper.S_bm_inode_start
			fmt.Println("inodes:", inodes)
			fmt.Println("bm_inodes_start:", bm_inodes_start)

			//creamos file
			file2, _ := os.Create("./reports/" + file_name + ".txt")

			buffer := make([]byte, 20)
			for i := 0; i < int(inodes); i = i + 20 {
				for j := 0; j < 20; j++ {
					if i+j >= int(inodes) {
						break
					}
					var bm byte
					if err := utilities.ReadObject(file, &bm, int64(bm_inodes_start)); err != nil {
						return
					}
					bm_inodes_start++
					buffer[j] = bm
				}
				fmt.Println("buffer:", buffer)
				file2.WriteString(fmt.Sprintf("%d", buffer))
				file2.WriteString("\n")
			}
			file2.Close()
		}
	}
}

func RepBMBlock(id string, path string) {
	name := strings.Split(path, "/")
	file_name := name[len(name)-1]
	file_name = strings.Split(file_name, ".")[0]
	fmt.Println(file_name)

	driveletter := id[0]
	driveletter = byte(strings.ToUpper(string(driveletter))[0])
	fmt.Println("driveletter:", string(driveletter))
	//abrimos dsk
	rutas := "./MIA/P1/" + string(driveletter) + ".dsk"
	//fmt.Println("ruta:", ruta)
	file, err := utilities.OpenFile(rutas)
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
			//leemos superblock
			var tmpSuper objs.Superblock
			start := int64(tmpMbr.Mbr_partitions[i].Part_start)
			if err := utilities.ReadObject(file, &tmpSuper, start); err != nil {
				return
			}
			//leemos bitmap de inodos
			inodes := tmpSuper.S_blocks_count
			bm_inodes_start := tmpSuper.S_bm_block_start
			fmt.Println("blocks:", inodes)
			fmt.Println("blocks_start:", bm_inodes_start)

			//creamos file
			file2, _ := os.Create("./reports/" + file_name + ".txt")

			buffer := make([]byte, 20)
			for i := 0; i < int(inodes); i = i + 20 {
				for j := 0; j < 20; j++ {
					if i+j >= int(inodes) {
						break
					}
					var bm byte
					if err := utilities.ReadObject(file, &bm, int64(bm_inodes_start)); err != nil {
						return
					}
					bm_inodes_start++
					buffer[j] = bm
				}
				fmt.Println("buffer:", buffer)
				file2.WriteString(fmt.Sprintf("%d", buffer))
				file2.WriteString("\n")
			}
			file2.Close()
		}
	}
}

func RepSB(id string, path string) {
	name := strings.Split(path, "/")
	file_name := name[len(name)-1]
	file_name = strings.Split(file_name, ".")[0]
	fmt.Println(file_name)

	driveletter := id[0]
	driveletter = byte(strings.ToUpper(string(driveletter))[0])
	fmt.Println("driveletter:", string(driveletter))
	//abrimos dsk
	rutas := "./MIA/P1/" + string(driveletter) + ".dsk"
	//fmt.Println("ruta:", ruta)
	file, err := utilities.OpenFile(rutas)
	if err != nil {
		return
	}
	//leemos mbr
	var tmpMbr objs.MBR
	if err := utilities.ReadObject(file, &tmpMbr, 0); err != nil {
		return
	}
	dotContent := ""
	//buscamos particion
	for i := 0; i < 4; i++ {
		if string(tmpMbr.Mbr_partitions[i].Part_id[:]) == strings.ToUpper(id) {
			//leemos superblock
			var tmpSuper objs.Superblock
			start := int64(tmpMbr.Mbr_partitions[i].Part_start)
			if err := utilities.ReadObject(file, &tmpSuper, start); err != nil {
				return
			}
			//creamos dot
			dotContent := `digraph G {
						node [shape=plaintext]
						node1 [label=<
						<TABLE BORDER="1" CELLBORDER="1" CELLSPACING="0">
						<TR><TD bgcolor="#ed5364" COLSPAN="2">Reporte Superblock</TD></TR>`
			dotContent += `	<TR> <TD >s_filesystem_type</TD> <TD>` + fmt.Sprint(tmpSuper.S_filesystem_type) + `</TD> </TR>`
			dotContent += `	<TR> <TD >s_inodes_count</TD> <TD>` + fmt.Sprint(tmpSuper.S_inodes_count) + `</TD> </TR>`
			dotContent += `	<TR> <TD >s_blocks_count</TD> <TD>` + fmt.Sprint(tmpSuper.S_blocks_count) + `</TD> </TR>`
			dotContent += `	<TR> <TD >s_free_blocks_count</TD> <TD>` + fmt.Sprint(tmpSuper.S_free_blocks_count) + `</TD> </TR>`
			dotContent += `	<TR> <TD >s_free_inodes_count</TD> <TD>` + fmt.Sprint(tmpSuper.S_free_inodes_count) + `</TD> </TR>`
			dotContent += `	<TR> <TD >s_mtime</TD> <TD>` + string(tmpSuper.S_mtime[:]) + `</TD> </TR>`
			dotContent += `	<TR> <TD >s_umtime</TD> <TD>` + string(tmpSuper.S_umtime[:]) + `</TD> </TR>`
			dotContent += `	<TR> <TD >s_mnt_count</TD> <TD>` + fmt.Sprint(tmpSuper.S_mnt_count) + `</TD> </TR>`
			dotContent += `	<TR> <TD >s_magic</TD> <TD>` + fmt.Sprint(tmpSuper.S_magic) + `</TD> </TR>`
			dotContent += `	<TR> <TD >s_inode_size</TD> <TD>` + fmt.Sprint(tmpSuper.S_inode_s) + `</TD> </TR>`
			dotContent += `	<TR> <TD >s_block_size</TD> <TD>` + fmt.Sprint(tmpSuper.S_block_s) + `</TD> </TR>`
			dotContent += `	<TR> <TD >s_first_ino</TD> <TD>` + fmt.Sprint(tmpSuper.S_fist_ino) + `</TD> </TR>`
			dotContent += `	<TR> <TD >s_first_blo</TD> <TD>` + fmt.Sprint(tmpSuper.S_first_blo) + `</TD> </TR>`
			dotContent += `	<TR> <TD >s_bm_inode_start</TD> <TD>` + fmt.Sprint(tmpSuper.S_bm_inode_start) + `</TD> </TR>`
			dotContent += `	<TR> <TD >s_bm_block_start</TD> <TD>` + fmt.Sprint(tmpSuper.S_bm_block_start) + `</TD> </TR>`
			dotContent += `	<TR> <TD >s_inode_start</TD> <TD>` + fmt.Sprint(tmpSuper.S_inode_start) + `</TD> </TR>`
			dotContent += `	<TR> <TD >s_block_start</TD> <TD>` + fmt.Sprint(tmpSuper.S_block_start) + `</TD> </TR>`
			dotContent += `</TABLE>>]}`

			fmt.Println(dotContent)
			//objs.PrintSuperblock(tmpSuper)
		}
	}
	generate(dotContent, path)
}

func generate(dotContent string, path string) {
	file := path

	// Guardar el contenido DOT en un archivo
	fileName := "./reports/" + file
	if err := os.WriteFile(fileName, []byte(dotContent), 0644); err != nil {
		fmt.Printf("Error al guardar el archivo DOT: %v\n", err)
		return
	}
	fmt.Printf("Archivo DOT guardado correctamente: %s\n", fileName)

}
