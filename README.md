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


## Usage:

```bash

# replace e with the name/path to the antas flavour

<e> <path/to/file.pdf> [options]

Options:
    -f, --format: in what format antas should make it's report [ json | human ] (default "human")
    -h, --help: print this message
```

when running antas two non-zero exit codes may be returned

1. something went wrong inside of antas (i.e: file not found)
2. you passed an unknown, incorrect flag value (i.e: -f "unknown formatter")

on exit code 2, you'll have to deal with an unstructured string

on exit code 1 (and if using -f json) the error will be returned as so in stdout

```json
{
    "ok": false,
    "error": string
}
```

if the exit code is 0 and the value `json` was passed to the `-f, --format` flag

antas will give you the given payload to parse

```json
{
    "ok": true,
    "out_dir": string,
    "page_count": int
}
```

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

command: `go build -tags="natif pdfium_use_turbojpeg" antas-natif .`

This one is the fastest but requires a few stuffs installed and configured on your system

- CGO enabled
- a compiled pdfium version for your system

    some precompiled releases are available at [https://github.com/bblanchon/pdfium-binaries/releases](https://github.com/bblanchon/pdfium-binaries/releases) unless you feel like compiling it yourself.

- pkg-config

the instructions relative setup for it is available at [https://github.com/klippa-app/go-pdfium#configure-pkg-config](https://github.com/klippa-app/go-pdfium#configure-pkg-config)

Make sure you extend your library path when running:
export LD_LIBRARY_PATH={path}/lib

I generally recommend not removing the installed libraries as they're still needed at runtime.

one may find convenient to add this at the end of their shell configuration file (i.e: .zshrc, .bashrc, etc...)

```bash
export PKG_CONFIG_PATH=/opt/pdfium/lib/pkgconfig
export LD_LIBRARY_PATH=/opt/pdfium/lib
```

### 3. turbo

command: `go build -tags pdfium_use_turbo_jpeg -o antas-turbo .`

this one is apparently supposed to use a faster encoding system via libjpeg-turbo, although most benchmarsk turned out to make him the slowest of the family.
i wouldn't recommend using it unless your the benchmark on your machine says otherwise

### Sandboxing

in most cases you'll probably want to use antas-natif, if it runs in a dedicated server
the docker file will provide a simple built docker image 

you can then sandbox antas-natif with control of the system.