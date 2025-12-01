package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

// Códigos ANSI para colores y estilos
const (
	ColorReset  = "\033[0m"
	ColorRed    = "\033[31m"
	ColorGreen  = "\033[32m"
	ColorYellow = "\033[33m"
	ColorBlue   = "\033[34m"
	ColorPurple = "\033[35m"
	ColorCyan   = "\033[36m"
	ColorWhite  = "\033[37m"
	Bold        = "\033[1m"
)

// Task representa una tarea en la lista
type Task struct {
	ID          int        `json:"id"`
	Description string     `json:"description"`
	Completed   bool       `json:"completed"`
	CreatedAt   time.Time  `json:"created_at"`
	CompletedAt *time.Time `json:"completed_at,omitempty"`
}

// TodoList maneja la lista de tareas
type TodoList struct {
	Tasks    []Task `json:"tasks"`
	NextID   int    `json:"next_id"`
	filename string
}

// NewTodoList crea una nueva lista de tareas
func NewTodoList(filename string) *TodoList {
	return &TodoList{
		Tasks:    []Task{},
		NextID:   1,
		filename: filename,
	}
}

// Load carga las tareas desde un archivo JSON
func (tl *TodoList) Load() error {
	file, err := os.ReadFile(tl.filename)
	if err != nil {
		if os.IsNotExist(err) {
			return nil // Archivo no existe, comenzar con lista vacía
		}
		return err
	}

	return json.Unmarshal(file, tl)
}

// Save guarda las tareas en un archivo JSON
func (tl *TodoList) Save() error {
	data, err := json.MarshalIndent(tl, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(tl.filename, data, 0644)
}

// AddTask agrega una nueva tarea a la lista con la descripción proporcionada.
// La función crea una tarea con un ID único autoincrementable, establece la fecha
// de creación al momento actual y la marca como no completada por defecto.
// Después de agregar la tarea, incrementa el contador NextID para la próxima tarea
// e imprime un mensaje de confirmación con el ID asignado.
//
// Parámetros:
//   - description: texto descriptivo de la tarea a agregar
func (tl *TodoList) AddTask(description string) {
	task := Task{
		ID:          tl.NextID,
		Description: description,
		Completed:   false,
		CreatedAt:   time.Now(),
	}
	tl.Tasks = append(tl.Tasks, task)
	tl.NextID++
	fmt.Printf("✓ Tarea agregada con ID: %d\n", task.ID)
}

// ListTasks muestra todas las tareas
func (tl *TodoList) ListTasks() {
	if len(tl.Tasks) == 0 {
		fmt.Println("No hay tareas en la lista.")
		return
	}

	fmt.Println("\n=== LISTA DE TAREAS ===")
	for _, task := range tl.Tasks {
		status := "⬜"
		if task.Completed {
			status = "✅"
		}
		fmt.Printf("%s ID: %d - %s\n", status, task.ID, task.Description)
		fmt.Printf("    Creada: %s\n", task.CreatedAt.Format("2006-01-02 15:04:05"))
		if task.CompletedAt != nil {
			fmt.Printf("    Completada: %s\n", task.CompletedAt.Format("2006-01-02 15:04:05"))
		}
		fmt.Println()
	}
}

// ListCompletedTasks muestra solamente las tareas completadas
func (tl *TodoList) ListCompletedTasks() {
	completedTasks := []Task{}
	for _, task := range tl.Tasks {
		if task.Completed {
			completedTasks = append(completedTasks, task)
		}
	}

	if len(completedTasks) == 0 {
		fmt.Println("No hay tareas completadas.")
		return
	}

	fmt.Println("\n=== TAREAS COMPLETADAS ===")
	for _, task := range completedTasks {
		fmt.Printf("✅ ID: %d - %s\n", task.ID, task.Description)
		fmt.Printf("    Creada: %s\n", task.CreatedAt.Format("2006-01-02 15:04:05"))
		if task.CompletedAt != nil {
			fmt.Printf("    Completada: %s\n", task.CompletedAt.Format("2006-01-02 15:04:05"))
		}
		fmt.Println()
	}
}

// CompleteTask marca una tarea como completada
func (tl *TodoList) CompleteTask(id int) error {
	for i := range tl.Tasks {
		if tl.Tasks[i].ID == id {
			if tl.Tasks[i].Completed {
				return fmt.Errorf("la tarea ya está completada")
			}
			now := time.Now()
			tl.Tasks[i].Completed = true
			tl.Tasks[i].CompletedAt = &now
			fmt.Printf("🎉 ✅ Tarea %d marcada como completada 🎊\n", id)
			return nil
		}
	}
	return fmt.Errorf("tarea con ID %d no encontrada", id)
}

// DeleteTask elimina una tarea de la lista
func (tl *TodoList) DeleteTask(id int) error {
	for i, task := range tl.Tasks {
		if task.ID == id {
			tl.Tasks = append(tl.Tasks[:i], tl.Tasks[i+1:]...)
			fmt.Printf("✓ Tarea %d eliminada\n", id)
			return nil
		}
	}
	return fmt.Errorf("tarea con ID %d no encontrada", id)
}

// ShowStats muestra estadísticas de las tareas
func (tl *TodoList) ShowStats() {
	total := len(tl.Tasks)
	completed := 0
	pending := 0

	for _, task := range tl.Tasks {
		if task.Completed {
			completed++
		} else {
			pending++
		}
	}

	fmt.Println("\n=== ESTADÍSTICAS ===")
	fmt.Printf("Total de tareas: %d\n", total)
	fmt.Printf("Completadas: %d\n", completed)
	fmt.Printf("Pendientes: %d\n", pending)
	if total > 0 {
		fmt.Printf("Progreso: %.1f%%\n", float64(completed)/float64(total)*100)
	}
	fmt.Println()
}

// Backup crea una copia de seguridad del archivo de tareas
func (tl *TodoList) Backup() error {
	// Verificar si el archivo de tareas existe
	if _, err := os.Stat(tl.filename); os.IsNotExist(err) {
		return fmt.Errorf("no hay tareas para respaldar")
	}

	// Generar nombre del archivo de backup con timestamp
	timestamp := time.Now().Format("2006-01-02_15-04-05")
	backupFilename := fmt.Sprintf("tasks_backup_%s.json", timestamp)

	// Leer el contenido del archivo original
	data, err := os.ReadFile(tl.filename)
	if err != nil {
		return fmt.Errorf("error al leer archivo de tareas: %v", err)
	}

	// Escribir el backup
	if err := os.WriteFile(backupFilename, data, 0644); err != nil {
		return fmt.Errorf("error al crear backup: %v", err)
	}

	fmt.Printf("✓ Backup creado exitosamente: %s\n", backupFilename)
	return nil
}

func main() {
	todoList := NewTodoList("tasks.json")

	// Cargar tareas existentes
	if err := todoList.Load(); err != nil {
		fmt.Printf("Error al cargar tareas: %v\n", err)
	}

	reader := bufio.NewReader(os.Stdin)

	fmt.Println("╔════════════════════════════════════╗")
	fmt.Println("║   GESTOR DE LISTA DE TAREAS GO    ║")
	fmt.Println("╚════════════════════════════════════╝")

	for {
		fmt.Println("\n--- MENÚ PRINCIPAL ---")
		fmt.Printf("%s%s1.%s %s%sAgregar tarea%s\n", ColorCyan, Bold, ColorReset, Bold, ColorWhite, ColorReset)
		fmt.Printf("%s%s2.%s %s%sListar tareas%s\n", ColorGreen, Bold, ColorReset, Bold, ColorWhite, ColorReset)
		fmt.Printf("%s%s3.%s %s%sListar tareas completadas%s\n", ColorGreen, Bold, ColorReset, Bold, ColorWhite, ColorReset)
		fmt.Printf("%s%s4.%s %s%sCompletar tarea%s\n", ColorYellow, Bold, ColorReset, Bold, ColorWhite, ColorReset)
		fmt.Printf("%s%s5.%s %s%sEliminar tarea%s\n", ColorRed, Bold, ColorReset, Bold, ColorWhite, ColorReset)
		fmt.Printf("%s%s6.%s %s%sVer estadísticas%s\n", ColorPurple, Bold, ColorReset, Bold, ColorWhite, ColorReset)
		fmt.Printf("%s%s7.%s %s%sCrear backup%s\n", ColorCyan, Bold, ColorReset, Bold, ColorWhite, ColorReset)
		fmt.Printf("%s%s8.%s %s%sSalir%s\n", ColorBlue, Bold, ColorReset, Bold, ColorWhite, ColorReset)
		fmt.Print("\nSelecciona una opción: ")

		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		// Mostrar la opción seleccionada en un cuadro
		if input != "" {
			// Determinar el color según la opción elegida
			optionColor := ColorWhite
			switch input {
			case "1":
				optionColor = ColorCyan
			case "2":
				optionColor = ColorGreen
			case "3":
				optionColor = ColorGreen
			case "4":
				optionColor = ColorYellow
			case "5":
				optionColor = ColorRed
			case "6":
				optionColor = ColorPurple
			case "7":
				optionColor = ColorCyan
			case "8":
				optionColor = ColorBlue
			}

			fmt.Printf("\n%s┌─────────────────────┐%s\n", ColorWhite, ColorReset)
			fmt.Printf("%s│%s Opción elegida: %s%s%s%s   %s│%s\n", ColorWhite, ColorReset, optionColor, Bold, input, ColorReset, ColorWhite, ColorReset)
			fmt.Printf("%s└─────────────────────┘%s\n", ColorWhite, ColorReset)
		}

		switch input {
		case "1":
			fmt.Print("Descripción de la tarea: ")
			description, _ := reader.ReadString('\n')
			description = strings.TrimSpace(description)
			if description != "" {
				todoList.AddTask(description)
				_ = todoList.Save()
			} else {
				fmt.Println("La descripción no puede estar vacía.")
			}

		case "2":
			todoList.ListTasks()

		case "3":
			todoList.ListCompletedTasks()

		case "4":
			fmt.Print("ID de la tarea a completar: ")
			idStr, _ := reader.ReadString('\n')
			idStr = strings.TrimSpace(idStr)
			id, err := strconv.Atoi(idStr)
			if err != nil {
				fmt.Println("ID inválido.")
				continue
			}
			if err := todoList.CompleteTask(id); err != nil {
				fmt.Printf("Error: %v\n", err)
			} else {
				_ = todoList.Save()
			}

		case "5":
			fmt.Print("ID de la tarea a eliminar: ")
			idStr, _ := reader.ReadString('\n')
			idStr = strings.TrimSpace(idStr)
			id, err := strconv.Atoi(idStr)
			if err != nil {
				fmt.Println("ID inválido.")
				continue
			}
			if err := todoList.DeleteTask(id); err != nil {
				fmt.Printf("Error: %v\n", err)
			} else {
				_ = todoList.Save()
			}

		case "6":
			todoList.ShowStats()

		case "7":
			if err := todoList.Backup(); err != nil {
				fmt.Printf("Error: %v\n", err)
			}

		case "8":
			fmt.Println("¡Hasta luego!")
			return

		default:
			fmt.Println("Opción inválida. Por favor, selecciona una opción del 1 al 8.")
		}
	}
}
