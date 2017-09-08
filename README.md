# MinoloBot
Bot de Telegram oficial del miembro de [Sugus](https://sugus.eii.us.es/) _Minolo_.

## Requisitos
- **Go**
- **Glide**

Podrás encontrar paquetes en la mayoría de distribuciones.

## Compilación
1. Descarga MinoloBot con el comando `get` de Go.  
`go get github.com/SUGUS-GNULinux/minolobot`
- Ahora la raíz del proyecto debe encontrarse en:
`$GOPATH/src/github.com/SUGUS-GNULinux/minolobot`,
descarga las dependencias mediante [Glide](https://github.com/Masterminds/glide) **estando en la raíz del proyecto** mediante `glide install`.
- Compilar con Go (respeta mantener el directorio datafiles junto con el binario), o bien ejecuta con `go run main.go`

## Asegúrate de
- Guardar el token del bot de Telegram en un archivo llamado `token` en la siguiente ruta: `[...]/minolobot/datafiles/`
