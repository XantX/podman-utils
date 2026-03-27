# Podutil

[![Go version](https://img.shields.io/github/go-mod/go-version/XantX/podman-utils)](https://github.com/XantX/podman-utils)
[![Latest release](https://img.shields.io/github/v/release/XantX/podman-utils)](https://github.com/XantX/podman-utils/releases/latest)
[![License](https://img.shields.io/github/license/XantX/podman-utils)](LICENSE)
[![Platform](https://img.shields.io/badge/platform-linux%20%7C%20windows%20%7C%20macOS-blue)]()

Herramienta CLI en Go que actúa como un superset de los comandos de Podman, diseñada para mejorar la experiencia de usuario en la gestión de contenedores, imágenes, volúmenes y demás recursos de Podman.

## Descripción

Podutil ofrece una interfaz de línea de comandos enriquecida sobre Podman, proporcionando:

- **Gestión simplificada** de contenedores, imágenes y volúmenes
- **Mejores flujos de trabajo** para las operaciones más comunes
- **Experiencia de usuario mejorada** con salida más legible y comandos intuitivos
- **Interfaz interactiva** TUI para seleccionar contenedores

## Características

- Comandos superset de Podman con opciones adicionales
- Interfaz TUI interactiva Bubble Tea
- Aliases para comandos frecuentes
- Operaciones simplificadas para tareas complejas

## Requisitos

- Go 1.21+
- Podman instalado y configurado

## Instalación

### Opción 1: go install (recomendado)

```bash
go install github.com/podutil/podman-utils@latest
```

### Opción 2: Scripts de instalación

**Linux/macOS:**
```bash
curl -sL https://raw.githubusercontent.com/XantX/podman-utils/master/scripts/install.sh | bash
```

**Windows (PowerShell):**
```powershell
irm -Uri https://raw.githubusercontent.com/XantX/podman-utils/master/scripts/install.ps1 | iex
```

O descarga los scripts desde la sección [Releases](https://github.com/XantX/podman-utils/releases).

### Opción 3: Descarga manual

1. Ir a [Releases](https://github.com/XantX/podman-utils/releases)
2. Descargar el binario para tu SO/arquitectura
3. Extraer y agregar al PATH

### Agregar al PATH

**Linux/macOS:**
```bash
export PATH=$PATH:$HOME/.local/bin
# Para hacer permanente:
echo 'export PATH=$PATH:$HOME/.local/bin' >> ~/.bashrc
```

**Windows (PowerShell):**
```powershell
# Temporal (solo esta sesión):
$env:PATH += ";C:\Users\TU_USUARIO\go\bin"

# Permanente:
[Environment]::SetEnvironmentVariable("PATH", $env:PATH + ";C:\Users\TU_USUARIO\go\bin", "User")
```

## Uso

### Iniciar contenedor
```bash
podutil start <container_id>  # Inicia directamente
podutil start                 # Muestra lista interactiva
```

### Detener contenedor
```bash
podutil stop <container_id>   # Detiene directamente
podutil stop                  # Muestra lista interactiva
```

### Ver contenedores corriendo
```bash
podutil ps                    # Muestra lista + seleccionar para detalles
```

## Comandos

| Comando | Descripción |
|---------|-------------|
| `podutil start [id]` | Inicia un contenedor. Sin ID muestra lista interactiva |
| `podutil stop [id]` | Detiene un contenedor. Sin ID muestra lista interactiva |
| `podutil ps` | Lista contenedores corriendo con detalles |
| `podutil help` | Muestra ayuda |

## Desarrollo

### Requisitos

- Go 1.21+
- Podman

### Compilar

```bash
CGO_ENABLED=0 go build -o podutil ./cmd
```

### Release

```bash
git tag v1.0.0
git push origin v1.0.0
```

## Licencia

MIT