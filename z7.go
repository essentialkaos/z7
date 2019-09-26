// Package z7 provides methods for working with 7z archives (p7zip wrapper)
package z7

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                     Copyright (c) 2009-2019 ESSENTIAL KAOS                         //
//        Essential Kaos Open Source License <https://essentialkaos.com/ekol>         //
//                                                                                    //
// ////////////////////////////////////////////////////////////////////////////////// //

import (
	"errors"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"pkg.re/essentialkaos/ek.v11/fsutil"
	"pkg.re/essentialkaos/ek.v11/mathutil"
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
func Add(props Props, files ...string) (string, error) {
	return AddList(props, files)
}

// AddList add files as string slice
func AddList(props Props, files []string) (string, error) {
	if len(files) == 0 {
		return "", errors.New("You should define files to compress")
	}

	err := props.Validate(false)

	if err != nil {
		return "", err
	}

	if props.Dir != "" {
		fsutil.Push(props.Dir)
	}

	out, err := execBinary(_COMMAND_ADD, props, files)

	if props.Dir != "" {
		fsutil.Pop()
	}

	return out, err
}

// Extract extract arhive
func Extract(props Props) (string, error) {
	err := props.Validate(true)

	if err != nil {
		return "", err
	}

	return execBinary(_COMMAND_EXTRACT, props, nil)
}

// List return info about archive
func List(props Props) (*Info, error) {
	err := props.Validate(true)

	if err != nil {
		return nil, err
	}

	out, err := execBinary(_COMMAND_LIST, props, nil)

	if err != nil {
		return nil, err
	}

	return parseInfoString(out), nil
}

// Check test archive
func Check(props Props) (bool, error) {
	err := props.Validate(true)

	if err != nil {
		return false, err
	}

	out, err := execBinary(_COMMAND_TEST, props, nil)

	if err != nil {
		return false, err
	}

	outData := strings.Split(out, "\n")

	for index, line := range outData {
		if line == _TEST_OK_VALUE {
			return true, nil
		} else if line == _TEST_ERROR_VALUE {
			return false, fmt.Errorf(outData[index+1])
		}
	}

	return false, errors.New("Can't parse 7zip output")
}

// Delete remove files from archive
func Delete(props Props, files ...string) (string, error) {
	if len(files) == 0 {
		return "", errors.New("You should define files to delete")
	}

	err := props.Validate(true)

	if err != nil {
		return "", err
	}

	return execBinary(_COMMAND_DELETE, props, files)
}

// ////////////////////////////////////////////////////////////////////////////////// //

// Validate validate properties values
func (p Props) Validate(checkFile bool) error {
	if checkFile {
		if !fsutil.IsExist(p.File) {
			return fmt.Errorf("File %s does not exist", p.File)
		}

		if !fsutil.IsReadable(p.File) {
			return fmt.Errorf("File %s is not readable", p.File)
		}
	}

	if p.IncludeFile != "" {
		if !fsutil.IsExist(p.IncludeFile) {
			return fmt.Errorf("Included file %s does not exist", p.IncludeFile)
		}

		if !fsutil.IsReadable(p.IncludeFile) {
			return fmt.Errorf("Included file %s is not readable", p.IncludeFile)
		}
	}

	if p.ExcludeFile != "" {
		if !fsutil.IsExist(p.ExcludeFile) {
			return fmt.Errorf("Included file %s does not exist", p.ExcludeFile)
		}

		if !fsutil.IsReadable(p.ExcludeFile) {
			return fmt.Errorf("Included file %s is not readable", p.ExcludeFile)
		}
	}

	if p.OutputDir != "" {
		if !fsutil.IsWritable(p.OutputDir) {
			return fmt.Errorf("Directory %s is not writable", p.OutputDir)
		}
	}

	return nil
}

// ToArgs convert properties to p7zip arguments
func (p Props) ToArgs(command string) []string {
	var args = []string{p.File, "", "-y", "-bd"}

	if command == _COMMAND_ADD {
		var compression int

		if p.Compression == 0 {
			compression = _COMPRESSION_DEFAULT
		} else {
			compression = mathutil.Between(p.Compression, _COMPRESSION_MIN, _COMPRESSION_MAX)
		}

		args = append(args, "-mx="+strconv.Itoa(compression))

		switch {
		case p.Threads < 1:
			args = append(args, "-mmt=1")
		case p.Threads >= 1:
			args = append(args, "-mmt="+strconv.Itoa(mathutil.Between(p.Threads, 1, 128)))
		}

		if p.Exclude != "" {
			args = append(args, "-x"+p.Exclude)
		} else if p.ExcludeFile != "" {
			args = append(args, "-xr@"+p.ExcludeFile)
		}

		if p.IncludeFile != "" {
			args = append(args, "-ir@"+p.IncludeFile)
		}

	} else if command == _COMMAND_EXTRACT {
		if p.OutputDir != "" {
			args = append(args, "-o"+p.OutputDir)
		}
	} else if command == _COMMAND_LIST {
		args = append(args, "-slt")
	}

	if p.Password != "" {
		args = append(args, "-p"+p.Password)
	}

	if p.Recursive {
		args = append(args, "-r")
	}

	if p.WorkingDir != "" {
		args = append(args, "-w"+p.WorkingDir)
	}

	return args
}

// ////////////////////////////////////////////////////////////////////////////////// //

// execBinary exec 7zip binary
func execBinary(command string, props Props, files []string) (string, error) {
	args := props.ToArgs(command)

	if len(files) != 0 {
		args = append(args, files...)
	}

	cmd := exec.Command(_BINARY)

	cmd.Args = append(cmd.Args, command)
	cmd.Args = append(cmd.Args, args...)

	out, err := cmd.Output()

	if err != nil {
		return string(out[:]), errors.New(string(out[:]))
	}

	return string(out[:]), nil
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
