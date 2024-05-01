package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	Analizador "MIA_P1_202004796/analizador"
	"MIA_P1_202004796/cmds"
	"MIA_P1_202004796/objs"
	"MIA_P1_202004796/utilities"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {

	//Analizador.Analizar()

	//Crear nuestra aplicación de Fiber
	app := fiber.New()

	app.Use(cors.New())

	app.Get("/", func(c *fiber.Ctx) error {
		c.Set("Access-Control-Allow-Origin", "*")
		c.Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		c.Set("Access-Control-Allow-Headers", "Origin, Content-Type, Accept")
		return c.SendString("Hello, World 👋!")
	})

	app.Post("/cmds", func(c *fiber.Ctx) error {
		c.Set("Access-Control-Allow-Origin", "*")
		c.Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		c.Set("Access-Control-Allow-Headers", "Origin, Content-Type, Accept")
		data := c.FormValue("data")

		fmt.Println(data)
		Analizador.Analizar(data)

		response := struct {
			Message string `json:"message"`
		}{Message: "ok"}
		return c.Status(fiber.StatusOK).JSON(response)
	})

	app.Get("/disks", func(c *fiber.Ctx) error {
		c.Set("Access-Control-Allow-Origin", "*")
		c.Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		c.Set("Access-Control-Allow-Headers", "Origin, Content-Type, Accept")
		disks := listDisks()

		response := struct {
			Message []string `json:"disks"`
		}{Message: disks}
		//fmt.Println(response)
		return c.Status(fiber.StatusOK).JSON(response)
	})

	app.Post("/partitions", func(c *fiber.Ctx) error {
		c.Set("Access-Control-Allow-Origin", "*")
		c.Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		c.Set("Access-Control-Allow-Headers", "Origin, Content-Type, Accept")
		driveletter := c.FormValue("driveletter")
		partitions := listPartitions(driveletter)

		response := struct {
			Message []string `json:"partitions"`
		}{Message: partitions}
		//fmt.Println(response)
		return c.Status(fiber.StatusOK).JSON(response)
	})

	app.Post("/login", func(c *fiber.Ctx) error {
		c.Set("Access-Control-Allow-Origin", "*")
		c.Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		c.Set("Access-Control-Allow-Headers", "Origin, Content-Type, Accept")
		pass := c.FormValue("pass")
		user := c.FormValue("user")
		disk := c.FormValue("disk")
		part := c.FormValue("part")

		driveletter := disk[:1]
		fmt.Println("driveletter:", driveletter, "part:", part, "user:", user, "pass:", pass)
		//modificar para buscar el ID de la part y pasarselo al login
		id := searchId(part, driveletter)
		fmt.Println("id:", id, "|")
		if id == "nil" {
			response := struct {
				Message string `json:"message"`
			}{Message: "-err particion no montada o no formateada"}
			return c.Status(fiber.StatusOK).JSON(response)
		}
		msg := cmds.Login(user, pass, id)
		response := struct {
			Message string `json:"message"`
		}{Message: msg}
		//fmt.Println(response)
		return c.Status(fiber.StatusOK).JSON(response)
	})

	app.Get("/logout", func(c *fiber.Ctx) error {
		c.Set("Access-Control-Allow-Origin", "*")
		c.Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		c.Set("Access-Control-Allow-Headers", "Origin, Content-Type, Accept")
		msg := cmds.Logout()

		response := struct {
			Message string `json:"message"`
		}{Message: msg}
		//fmt.Println(response)
		return c.Status(fiber.StatusOK).JSON(response)
	})

	app.Get("/reports", func(c *fiber.Ctx) error {
		c.Set("Access-Control-Allow-Origin", "*")
		c.Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		c.Set("Access-Control-Allow-Headers", "Origin, Content-Type, Accept")
		reports := listReports()

		response := struct {
			Message []string `json:"reports"`
		}{Message: reports}
		//fmt.Println(response)
		return c.Status(fiber.StatusOK).JSON(response)
	})

	app.Post("/report", func(c *fiber.Ctx) error {
		c.Set("Access-Control-Allow-Origin", "*")
		c.Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		c.Set("Access-Control-Allow-Headers", "Origin, Content-Type, Accept")
		name := c.FormValue("name")
		report := getReport(name)

		return c.SendString(report)
	})

	log.Fatal(app.Listen(":3000"))
}

func listDisks() []string {
	files, err := os.ReadDir("./MIA/P1/")
	if err != nil {
		log.Fatal(err)
	}
	var disks []string
	for _, file := range files {
		disks = append(disks, file.Name())
	}
	if len(disks) == 0 {
		disks = append(disks, "")
	}
	return disks
}

func listPartitions(driveletter string) []string {
	driveletter = strings.ToUpper(driveletter)
	//abrimos dsk
	ruta := "./MIA/P1/" + string(driveletter) + ".dsk"
	fmt.Println("ruta:", ruta)
	file, err := utilities.OpenFile(ruta)
	if err != nil {
		return nil
	}

	//leemos mbr
	var tmpMbr objs.MBR
	if err := utilities.ReadObject(file, &tmpMbr, 0); err != nil {
		return nil
	}
	return objs.ListPartitions(tmpMbr)
}

func searchId(name string, driveletter string) string {
	driveletter = strings.ToUpper(driveletter)
	//abrimos dsk
	ruta := "./MIA/P1/" + string(driveletter) + ".dsk"
	fmt.Println("ruta:", ruta)
	file, err := utilities.OpenFile(ruta)
	if err != nil {
		return "nil"
	}

	//leemos mbr
	var tmpMbr objs.MBR
	if err := utilities.ReadObject(file, &tmpMbr, 0); err != nil {
		return "nil"
	}
	for i := 0; i < 4; i++ {
		if name == objs.ReturnPartitionName(tmpMbr.Mbr_partitions[i]) && string(tmpMbr.Mbr_partitions[i].Part_status[:]) == "1" {
			return string(tmpMbr.Mbr_partitions[i].Part_id[:])
		}
	}
	return "nil"
}

func listReports() []string {
	files, err := os.ReadDir("./reports/")
	if err != nil {
		log.Fatal(err)
	}
	var reps []string
	for _, file := range files {
		reps = append(reps, file.Name())
	}
	if len(reps) == 0 {
		reps = append(reps, "")
	}
	return reps
}

func getReport(name string) string {
	contenido, err := os.ReadFile("./reports/" + name)
	if err != nil {
		fmt.Println("Error al leer el archivo:", err)
		return ""
	}

	// Convierte el contenido a un string
	contenidoString := string(contenido)
	return contenidoString
}
