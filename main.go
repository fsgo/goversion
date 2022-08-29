// Copyright(C) 2022 github.com/fsgo  All Rights Reserved.
// Author: hidu <duv123@gmail.com>
// Date: 2022/8/28

package main

import (
	"debug/buildinfo"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

var helpMessage = `
Self-Update :
          go install github.com/fsgo/goversion@latest

Site    : https://github.com/fsgo/goversion
Version : 0.0.2
Date    : 2022-08-29
`

func init() {
	flag.Usage = func() {
		fo := flag.CommandLine.Output()
		fmt.Fprintf(fo, "prints the build information for Go executables\n")
		fmt.Fprintf(fo, "Usage of %s:\n", os.Args[0])
		fmt.Fprintf(fo, "  %s [-m] [-o=json] [file ...]\n", os.Args[0])
		flag.PrintDefaults()
		fmt.Fprintf(fo, "\n"+strings.TrimSpace(helpMessage)+"\n")
	}
}

var module = flag.Bool("m", false, "need module info")
var out = flag.String("o", "txt", `output format, support: txt、json、json_p
json: JSON, all in one line
json_p: JSON pretty, many lines
`)

func main() {
	flag.Parse()
	files := flag.Args()
	if len(files) == 0 {
		log.Fatalln("no executable file, should:", os.Args[0], "[-m] [file ...]")
	}
	var hasFail bool
	for i := 0; i < len(files); i++ {
		if !do(files[i]) {
			hasFail = true
		}
	}
	if hasFail {
		os.Exit(2)
	}
}

func do(f string) bool {
	ret := &result{
		Name: f,
	}
	ret.BuildInfo, ret.Error = buildinfo.ReadFile(f)
	fmt.Println(ret.String())
	return ret.Error == nil
}

type result struct {
	Name      string
	BuildInfo *buildinfo.BuildInfo
	Error     error
}

func (ret *result) String() string {
	switch *out {
	case "txt":
		return ret.txt()
	case "json":
		return ret.json(false)
	case "json_p":
		return ret.json(true)
	default:
		return fmt.Sprintf("not support format: %q", *out)
	}
}

func (ret *result) txt() string {
	var b strings.Builder
	b.WriteString(ret.Name)
	b.WriteString(" : ")
	if ret.Error != nil {
		b.WriteString("[Error] ")
		b.WriteString(ret.Error.Error())
		return b.String()
	}
	b.WriteString(ret.BuildInfo.GoVersion)

	if !*module {
		return b.String()
	}

	b.WriteString("\n")
	b.WriteString(fmt.Sprintf("\tpath  %s\n", ret.BuildInfo.Path))
	b.WriteString(fmt.Sprintf("\tmodule  %s %s\n", ret.BuildInfo.Main.Path, ret.BuildInfo.Main.Version))
	for i := 0; i < len(ret.BuildInfo.Deps); i++ {
		dep := ret.BuildInfo.Deps[i]
		b.WriteString(fmt.Sprintf("\tdep %s %s\n", dep.Path, dep.Version))
	}

	for i := 0; i < len(ret.BuildInfo.Settings); i++ {
		s := ret.BuildInfo.Settings[i]
		b.WriteString(fmt.Sprintf("\tbuild %s=%s\n", s.Key, s.Value))
	}

	return b.String()
}

func jsonEncode(val any, pretty bool) string {
	var bf []byte
	var err error
	if !pretty {
		bf, err = json.Marshal(val)
	} else {
		bf, err = json.MarshalIndent(val, " ", "  ")
	}
	if err != nil {
		return err.Error()
	}
	return string(bf)

}

func (ret *result) json(pretty bool) string {
	data := map[string]any{
		"Name": ret.Name,
	}
	if ret.Error != nil {
		data["Error"] = ret.Error.Error()
		return jsonEncode(data, pretty)
	}
	data["GoVersion"] = ret.BuildInfo.GoVersion

	if !*module {
		return jsonEncode(data, pretty)
	}
	data["Path"] = ret.BuildInfo.Path
	data["Main"] = ret.BuildInfo.Main.Path
	data["MainVersion"] = ret.BuildInfo.Main.Version

	deps := make([]any, 0)
	for i := 0; i < len(ret.BuildInfo.Deps); i++ {
		dep := ret.BuildInfo.Deps[i]
		item := map[string]any{
			"Module":  dep.Path,
			"Version": dep.Version,
		}
		deps = append(deps, item)
	}

	data["Deps"] = deps

	setting := make(map[string]string)
	for i := 0; i < len(ret.BuildInfo.Settings); i++ {
		s := ret.BuildInfo.Settings[i]
		setting[s.Key] = s.Value
	}
	data["Settings"] = setting

	return jsonEncode(data, pretty)
}
