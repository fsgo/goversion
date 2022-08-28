# goversion
Prints the build information for Go executables.

## 1. Install
```bash
go install github.com/fsgo/goversion@latest
```

## 2. Usage

### 2.1 help
```bash
# goversion -help
prints the build information for Go executables
Usage of goversion:
  goversion [-m] [-o=json] [file ...]
  -m	need module info
  -o string
    	output format, support: txt、json、json_p
    	json: JSON, all in one line
    	json_p: JSON pretty, many lines
    	 (default "txt")
```

### 2.2 basic usage
Basic use: print the executable file's Go Version
```bash
# goversion go_fmt
go_fmt : go1.19
```

```bash
# goversion ~/bin/hello.sh
~/bin/hello.sh : [Error] could not read Go build info from ~/bin/hello.sh: unrecognized file format
```

### 2.3 with more module info
```bash
# goversion -m go_fmt
go_fmt : go1.19
	module  github.com/fsgo/go_fmt (devel)
	dep github.com/google/go-cmp v0.5.8
	dep golang.org/x/mod v0.6.0-dev.0.20220419223038-86c51ed26bb4
	dep golang.org/x/tools v0.1.12
	build -compiler=gc
	build CGO_ENABLED=1
	build CGO_CFLAGS=
	build CGO_CPPFLAGS=
	build CGO_CXXFLAGS=
	build CGO_LDFLAGS=
	build GOARCH=amd64
	build GOOS=darwin
	build GOAMD64=v1
	build vcs=git
	build vcs.revision=80b433ca1cdce663890ce8bd6d094dbd0739dcdc
	build vcs.time=2022-08-19T15:36:20Z
	build vcs.modified=true
```


### 2.4 JSON output
```bash
# goversion -m -o json_p go_fmt
{
   "Deps": [
     {
       "Module": "github.com/google/go-cmp",
       "Version": "v0.5.8"
     },
     {
       "Module": "golang.org/x/mod",
       "Version": "v0.6.0-dev.0.20220419223038-86c51ed26bb4"
     },
     {
       "Module": "golang.org/x/tools",
       "Version": "v0.1.12"
     }
   ],
   "GoVersion": "go1.19",
   "Main": "github.com/fsgo/go_fmt",
   "MainVersion": "(devel)",
   "Name": "go_fmt",
   "Settings": {
     "-compiler": "gc",
     "CGO_CFLAGS": "",
     "CGO_CPPFLAGS": "",
     "CGO_CXXFLAGS": "",
     "CGO_ENABLED": "1",
     "CGO_LDFLAGS": "",
     "GOAMD64": "v1",
     "GOARCH": "amd64",
     "GOOS": "darwin",
     "vcs": "git",
     "vcs.modified": "true",
     "vcs.revision": "80b433ca1cdce663890ce8bd6d094dbd0739dcdc",
     "vcs.time": "2022-08-19T15:36:20Z"
   }
 }
```

### 2.5 many input files
```bash
# goversion go_fmt gops
go_fmt : go1.19
gops : go1.17
```