# Antas - (from the french word "Entasse")

antas is a cli tool that takes a path to a pdf file as input and render each page into a .png file

the pdf pages are always saved in the path:

**{temp_dir}/{input_file_hash}/{when}/page_{id}.png**

- `temp_dir` refers to the the temporary directory on your machine
- `input_file_hash` is the hash of the given file computed via the sha256 algorithm. so multiple files are grouped.
-  `when` is meant for uniqueness of runs (so threading, go routines, doesn't create races conditions). 
    `when` is computed with two value, Now() represented via [the Unix representation number](https://en.wikipedia.org/wiki/Unix_time) and [the current PID (Process Id of antas)](https://en.wikipedia.org/wiki/Process_identifier)

    TL.DR
    `{when} = {Now().UnixTime}-{pid}`

- `id` is simply the number of the page


## Compilation

antas, use go-pdfium to render the pdf page into actual images and other three build favor (2 direct/ 1 indirect)

- classic (webassembly, direct)
- native (CGO required and pdfium setup, direct)
- turbo (use libjpeg-turbo for the rendering system, indirect)

### 1. classic

this one will use the webassembly version of go-pdfium

command: `go build -o antas .`


in order to run the webassembly version, go-pdfium will use wazero under the hood, but you won't have to deal with that.
this also grants sandboxing out of the box, which may be useful if you have to deal with untrusted pdf files (i.e: users uploading files).
and also make cross-compilation easier since it doesn't use CGO.

however it is known to be much slower(2x) than the native version

> [!WARNING]
> The following may and will be a pain in the ass on non-linux based OS to compile and run, especially windows.

### 2. native

command: `go build -tags natif antas-natif .`


This one is the fastest but requires a few stuffs

- CGO enabled
- a compiled pdfium version for your system

    some precompiled releases are available at [https://github.com/bblanchon/pdfium-binaries/releases](https://github.com/bblanchon/pdfium-binaries/releases) unless you feel like compiling it yourself.

- pkg-config

the instructions relative setup for it is available at [https://github.com/klippa-app/go-pdfium#configure-pkg-config](https://github.com/klippa-app/go-pdfium#configure-pkg-config)

Make sure you extend your library path when running:
export LD_LIBRARY_PATH={path}/lib

