package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
)

type Predio struct {
	Codigo               string
	Ciudad               string
	Departamento         string
	Barrio               string
	Direccion            string
	AreaTerreno          string
	AreaConstruida       string
	DetalleDisponibilidad string
	Estrato              string
	Precio               float64
	TipoInmueble         string
	DatosAdicionales     string
}

func parsePrecio(precioStr string) (float64, error) {
	precioStr = strings.TrimSpace(precioStr)
	precioStr = strings.ReplaceAll(precioStr, ",", "")

	precio, err := strconv.ParseFloat(precioStr, 64)
	if err != nil {
		return 0, err
	}

	precio = precio / 1000000

	return precio, nil
}

func leerPredios(nombreArchivo string) ([]Predio, error) {
	file, err := os.Open(nombreArchivo)
	if err != nil {
		return nil, fmt.Errorf("error abriendo archivo: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)

	_, err = reader.Read()
	if err != nil {
		return nil, fmt.Errorf("error leyendo encabezados: %w", err)
	}

	var predios []Predio

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("error leyendo registro: %w", err)
		}

		if len(record) < 12 {
			log.Printf("Registro incompleto, saltando: %v", record)
			continue
		}

		precio, err := parsePrecio(record[9])
		if err != nil {
			log.Printf("Error parseando precio '%s': %v, saltando registro", record[9], err)
			continue
		}

		predio := Predio{
			Codigo:               record[0],
			Ciudad:               record[1],
			Departamento:         record[2],
			Barrio:               record[3],
			Direccion:            record[4],
			AreaTerreno:          record[5],
			AreaConstruida:       record[6],
			DetalleDisponibilidad: record[7],
			Estrato:              record[8],
			Precio:               precio,
			TipoInmueble:         record[10],
			DatosAdicionales:     record[11],
		}

		predios = append(predios, predio)
	}

	return predios, nil
}

func obtenerTop5MasCostosos(predios []Predio) []Predio {
	sort.Slice(predios, func(i, j int) bool {
		return predios[i].Precio > predios[j].Precio
	})

	if len(predios) > 5 {
		return predios[:5]
	}

	return predios
}

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Uso: go run main.go <archivo.csv>")
	}

	nombreArchivo := os.Args[1]

	predios, err := leerPredios(nombreArchivo)
	if err != nil {
		log.Fatalf("Error procesando archivo: %v", err)
	}

	if len(predios) == 0 {
		log.Fatal("No se encontraron predios válidos en el archivo")
	}

	top5 := obtenerTop5MasCostosos(predios)

	fmt.Println("========================================")
	fmt.Println("TOP 5 PREDIOS MÁS COSTOSOS")
	fmt.Println("========================================\n")

	for i, predio := range top5 {
		fmt.Printf("%d. Código: %s\n", i+1, predio.Codigo)
		fmt.Printf("   Ciudad: %s\n", predio.Ciudad)
		fmt.Printf("   Departamento: %s\n", predio.Departamento)
		fmt.Printf("   Barrio: %s\n", predio.Barrio)
		fmt.Printf("   Dirección: %s\n", predio.Direccion)
		fmt.Printf("   Tipo: %s\n", predio.TipoInmueble)
		fmt.Printf("   Precio: $%.2f\n", predio.Precio)
		fmt.Printf("   Área Terreno: %s\n", predio.AreaTerreno)
		fmt.Printf("   Área Construida: %s\n", predio.AreaConstruida)
		fmt.Printf("   Estrato: %s\n", predio.Estrato)
		fmt.Println("   ----------------------------------------")
	}
}
