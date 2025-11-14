# Diagrama de Secuencia - Todo App

```mermaid
sequenceDiagram
    autonumber
    participant U as Usuario
    participant F as Formulario
    participant JS as app.js
    participant D as Array todos
    participant L as Lista DOM

    rect rgb(200, 230, 255)
        Note over U,L: Carga Inicial
        U->>JS: Carga la página
        JS->>D: Inicializa array vacío
        JS->>L: render() - Lista vacía
    end

    rect rgb(200, 255, 200)
        Note over U,L: Agregar Tarea
        U->>F: Escribe "Nueva tarea"
        U->>F: Click en "Agregar"
        F->>JS: submit event
        JS->>JS: e.preventDefault()
        JS->>JS: trim() del input
        JS->>D: push({ id, text, done: false })
        JS->>F: Limpia input
        JS->>L: render()
        L->>L: Crea elementos li
        L->>U: Muestra tarea nueva
    end

    rect rgb(255, 240, 200)
        Note over U,L: Marcar como Completada
        U->>L: Click en checkbox
        L->>JS: change event
        JS->>JS: toggle(id)
        JS->>D: map() - invierte done
        JS->>L: render()
        L->>L: Agrega clase 'completed'
        L->>U: Muestra tarea tachada
    end

    rect rgb(255, 200, 200)
        Note over U,L: Eliminar Tarea
        U->>L: Click en "Eliminar"
        L->>JS: click event
        JS->>JS: removeItem(id)
        JS->>D: filter() - remueve item
        JS->>L: render()
        L->>L: Reconstruye lista sin item
        L->>U: Muestra lista actualizada
    end

    rect rgb(230, 200, 255)
        Note over U,L: Desmarcar Tarea
        U->>L: Click en checkbox marcado
        L->>JS: change event
        JS->>JS: toggle(id)
        JS->>D: map() - invierte done a false
        JS->>L: render()
        L->>L: Remueve clase 'completed'
        L->>U: Muestra tarea normal
    end
```

## Descripción de las Etapas

- **Azul claro**: Carga inicial de la aplicación
- **Verde**: Proceso de agregar una nueva tarea
- **Amarillo**: Marcar una tarea como completada
- **Rojo claro**: Eliminar una tarea
- **Morado claro**: Desmarcar una tarea completada

## Componentes Clave

- **Usuario**: Interactúa con la interfaz
- **Formulario**: Captura nuevas tareas
- **app.js**: Lógica de la aplicación
- **Array todos**: Almacenamiento en memoria
- **Lista DOM**: Renderizado visual de las tareas
