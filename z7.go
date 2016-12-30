// Package z7 provides methods for working with 7z archives (p7zip wrapper)
package z7

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                     Copyright (c) 2009-2016 Essential Kaos                         //
//      Essential Kaos Open Source License <http://essentialkaos.com/ekol?en>         //
//                                                                                    //
// ////////////////////////////////////////////////////////////////////////////////// //

import (
	"errors"
	"fmt"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"pkg.re/essentialkaos/ek.v6/fsutil"
	"pkg.re/essentialkaos/ek.v6/mathutil"
)

// ////////////////////////////////////////////////////////////////////////////////// //

// List of supported formats
const (
	TYPE_7Z   = "7z"
	TYPE_ZIP  = "zip"
	TYPE_GZIP = "gzip"
	TYPE_XZ   = "xz"
	TYPE_BZIP = "bzip2"
)

// ////////////////////////////////////////////////////////////////////////////////// //

const _BINARY = "7za"

const (
	_COMPRESSION_MIN     = 0
	_COMPRESSION_MAX     = 9
	_COMPRESSION_DEFAULT = 4
)

const (
	_COMMAND_ADD       = "a"
	_COMMAND_BENCHMARK = "b"
	_COMMAND_DELETE    = "d"
	_COMMAND_LIST      = "l"
	_COMMAND_TEST      = "t"
	_COMMAND_UPDATE    = "u"
	_COMMAND_EXTRACT   = "x"
)

const _TEST_OK_VALUE = "Everything is Ok"
const _TEST_ERROR_VALUE = "ERRORS:"

// ////////////////////////////////////////////////////////////////////////////////// //

// Props contains properties for packing/unpacking data
type Props struct {
	Dir         string // Directory with files (for relative paths)
	File        string // Output file name
	IncludeFile string // File with include filenames
	Exclude     string // Exclude filenames
	ExcludeFile string // File with exclude filenames
	Compression int    // Compression level (0-9)
	OutputDir   string // Output dir (for extract command)
	Password    string // Password
	Threads     int    // Number of CPU threads
	Recursive   bool   // Recurse subdirectories
	WorkingDir  string // Working dir
	Delete      bool   // Delete files after compression
}

// Info contains info about archive
type Info struct {
	Path         string
	Type         string
	Method       []string
	Solid        bool
	Blocks       int
	PhysicalSize int
	HeadersSize  int
	Files        []*FileInfo
}

// FileInfo contains info about file inside archive
type FileInfo struct {
	Path       string
	Folder     string
	Size       int
	PackedSize int
	Modified   time.Time
	Created    time.Time
	Accessed   time.Time
	Attributes string
	CRC        int
	Encrypted  bool
	Method     []string
	Block      int
	Comment    string
	HostOS     string
	Version    int
}

// ////////////////////////////////////////////////////////////////////////////////// //

// Add add files to archive
func Add(properties interface{}, files ...string) (string, error) {
	return AddList(properties, files)
}

// AddList add files as string slice
func AddList(properties interface{}, files ...[]string) (string, error) {
	props, err := processProps(properties, false)

	if err != nil {
		return "", err
	}

	file, err := filepath.Abs(props.File)

	if err != nil {
		return "", err
	}

	if props.Dir != "" {
		fsutil.Push(props.Dir)
	}

	args := propsToArgs(props, _COMMAND_ADD)

	if len(files) != 0 {
		args = append(args, files[0]...)
	}

	out, err := execBinary(file, _COMMAND_ADD, args)

	if props.Dir != "" {
		fsutil.Pop()
	}

	return out, err
}

// Extract extract arhive
func Extract(properties interface{}) (string, error) {
	props, err := processProps(properties, true)

	if err != nil {
		return "", err
	}

	file, err := filepath.Abs(props.File)

	if err != nil {
		return "", err
	}

	args := propsToArgs(props, _COMMAND_EXTRACT)
	out, err := execBinary(file, _COMMAND_EXTRACT, args)

	return out, err
}

// List return info about archive
func List(properties interface{}) (*Info, error) {
	props, err := processProps(properties, true)

	if err != nil {
		return &Info{}, err
	}

	args := propsToArgs(props, _COMMAND_LIST)
	out, err := execBinary(props.File, _COMMAND_LIST, args)

	if err != nil {
		return nil, errors.New(out)
	}

	return parseInfoString(out), nil
}

// Check test archive
func Check(properties interface{}) (bool, error) {
	props, err := processProps(properties, true)

	if err != nil {
		return false, err
	}

	args := propsToArgs(props, _COMMAND_TEST)
	out, _ := execBinary(props.File, _COMMAND_TEST, args)

	outData := strings.Split(out, "\n")

	for i, line := range outData {
		if line == _TEST_OK_VALUE {
			return true, nil
		} else if line == _TEST_ERROR_VALUE {
			return false, fmt.Errorf(outData[i+1])
		}
	}

	return false, fmt.Errorf("Unknown error")
}

// Delete remove files from archive
func Delete(properties interface{}, files ...string) (string, error) {
	props, err := processProps(properties, true)

	if err != nil {
		return "", err
	}

	args := propsToArgs(props, _COMMAND_DELETE)

	if len(files) != 0 {
		args = append(args, files...)
	}

	out, err := execBinary(props.File, _COMMAND_DELETE, args)

	return out, err
}

// ////////////////////////////////////////////////////////////////////////////////// //

// execBinary exec 7zip binary
func execBinary(target string, command string, args []string) (string, error) {
	var cmd = exec.Command(_BINARY)

	cmd.Args = append(cmd.Args, command)
	cmd.Args = append(cmd.Args, target)
	cmd.Args = append(cmd.Args, args...)

	out, err := cmd.Output()

	return string(out[:]), err
}

// processProps parse properties and return props struct
func processProps(properties interface{}, checkFile bool) (*Props, error) {
	var props *Props

	switch properties.(type) {
	case *Props:
		props = properties.(*Props)
	case string:
		props = &Props{File: properties.(string), Compression: _COMPRESSION_DEFAULT}
	default:
		return nil, fmt.Errorf("Unknown properties type")
	}

	if checkFile {
		if !fsutil.IsExist(props.File) {
			return nil, fmt.Errorf("File %s is not exist", props.File)
		}

		if !fsutil.IsReadable(props.File) {
			return nil, fmt.Errorf("File %s is not readable", props.File)
		}
	}

	if props.IncludeFile != "" {
		if !fsutil.IsExist(props.IncludeFile) {
			return nil, fmt.Errorf("Include file %s is not exist", props.IncludeFile)
		}

		if !fsutil.IsReadable(props.IncludeFile) {
			return nil, fmt.Errorf("Include file %s is not readable", props.IncludeFile)
		}
	}

	if props.ExcludeFile != "" {
		if !fsutil.IsExist(props.ExcludeFile) {
			return nil, fmt.Errorf("Include file %s is not exist", props.ExcludeFile)
		}

		if !fsutil.IsReadable(props.ExcludeFile) {
			return nil, fmt.Errorf("Include file %s is not readable", props.ExcludeFile)
		}
	}

	if props.OutputDir != "" {
		if !fsutil.IsWritable(props.OutputDir) {
			return nil, fmt.Errorf("Directory %s is not writable", props.OutputDir)
		}
	}

	return props, nil
}

// propsToArgs convert props struct to p7zip arguments
func propsToArgs(props *Props, command string) []string {
	var args = []string{"", "-y", "-bd"}

	if command == _COMMAND_ADD {
		compLvl := strconv.Itoa(mathutil.Between(props.Compression, _COMPRESSION_MIN, _COMPRESSION_MAX))

		args = append(args, "-mx="+compLvl)

		switch {
		case props.Threads < 1:
			args = append(args, "-mmt=1")
		case props.Threads >= 1:
			args = append(args, "-mmt="+strconv.Itoa(mathutil.Between(props.Threads, 1, 128)))
		}

		if props.Exclude != "" {
			args = append(args, "-x"+props.Exclude)
		} else if props.ExcludeFile != "" {
			args = append(args, "-xr@"+props.ExcludeFile)
		}

		if props.IncludeFile != "" {
			args = append(args, "-ir@"+props.IncludeFile)
		}

	} else if command == _COMMAND_EXTRACT {
		if props.OutputDir != "" {
			args = append(args, "-o"+props.OutputDir)
		}
	} else if command == _COMMAND_LIST {
		args = append(args, "-slt")
	}

	if props.Password != "" {
		args = append(args, "-p"+props.Password)
	}

	if props.Recursive {
		args = append(args, "-r")
	}

	if props.WorkingDir != "" {
		args = append(args, "-w"+props.WorkingDir)
	}

	return args
}

// parseInfoString process raw info data
func parseInfoString(infoData string) *Info {
	var data = strings.Split(infoData, "\n")
	var info = &Info{}

	header, headerEnd := extractInfoHeader(data)
	headerData := parseRecordData(header)

	info.Path = headerData["Path"]
	info.Type = headerData["Type"]
	info.Method = strings.Split(headerData["Method"], " ")

	if info.Type == TYPE_7Z {
		info.Solid = headerData["Solid"] == "+"

		info.Blocks, _ = strconv.Atoi(headerData["Blocks"])
		info.PhysicalSize, _ = strconv.Atoi(headerData["Physical Size"])
		info.HeadersSize, _ = strconv.Atoi(headerData["Headers Size"])
	}

	recStart := 0
	records := data[headerEnd : len(data)-1]

	for i, v := range records {
		if v == "" {
			info.Files = append(info.Files, parseFileInfo(records[recStart:i]))
			recStart = i + 1
		}
	}

	return info
}

// parseFileInfo process raw info about file/directory
func parseFileInfo(data []string) *FileInfo {
	var info = &FileInfo{}
	var recordData = parseRecordData(data)

	crc, _ := strconv.ParseInt(recordData["CRC"], 16, 0)

	info.Path = recordData["Path"]
	info.Folder = recordData["Folder"]
	info.Size, _ = strconv.Atoi(recordData["Size"])
	info.PackedSize, _ = strconv.Atoi(recordData["Packed Size"])
	info.Modified = parseDateString(recordData["Modified"])
	info.Created = parseDateString(recordData["Created"])
	info.Accessed = parseDateString(recordData["Accessed"])
	info.Attributes = recordData["Attributes"]
	info.CRC = int(crc)
	info.Comment = recordData["Comment"]
	info.Encrypted = recordData["Encrypted"] == "+"
	info.Method = strings.Split(recordData["Method"], " ")
	info.Block, _ = strconv.Atoi(recordData["Block"])
	info.HostOS = recordData["Host OS"]
	info.Version, _ = strconv.Atoi(recordData["Version"])

	return info
}

// parseRecordData parse raw record
func parseRecordData(data []string) map[string]string {
	var result = make(map[string]string)

	for _, rec := range data {
		if rec != "" {
			name, val := parseValue(rec)
			result[name] = val
		}
	}

	return result
}

// parseDateString parse date string
func parseDateString(data string) time.Time {
	if data == "" {
		return time.Time{}
	}

	year, _ := strconv.Atoi(data[0:4])
	month, _ := strconv.Atoi(data[5:7])
	day, _ := strconv.Atoi(data[8:10])
	hour, _ := strconv.Atoi(data[11:13])
	min, _ := strconv.Atoi(data[14:16])
	sec, _ := strconv.Atoi(data[17:19])

	return time.Date(year, time.Month(month), day, hour, min, sec, 0, time.UTC)
}

// extractInfoHeader extract header from raw info data
func extractInfoHeader(data []string) ([]string, int) {
	var start int
	var end int

	for i, v := range data {
		if v == "--" {
			start = i + 1
		}

		switch v {
		case "--":
			start = i + 1
		case "----------":
			end = i - 1
			break
		}
	}

	return data[start:end], end + 2
}

// parseValue parse "name = value" string
func parseValue(s string) (string, string) {
	valSlice := strings.Split(s, " = ")

	if len(valSlice) == 2 {
		return valSlice[0], valSlice[1]
	}

	return "", ""
}
