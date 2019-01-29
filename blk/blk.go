package blk

import (
	"errors"
	"io/ioutil"
	"path/filepath"
	"strings"
)

// Errors
var (
	ErrNotFound = errors.New("Not Found")
)

// Paths
var (
	MtabPath     = filepath.Join("/", "etc", "mtab")
	SysBlockPath = filepath.Join("/", "sys", "block")
	UdevDataPath = filepath.Join("/", "run", "udev", "data")
)

func mtab() (string, error) {
	b, err := ioutil.ReadFile(MtabPath)
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(b)), nil
}

// GetDevPath ...
func GetDevPath(mountpoint string) (string, error) {
	mtab, err := mtab()
	if err != nil {
		return "", err
	}

	lines := strings.Split(mtab, "\n")

	for _, line := range lines {
		fields := strings.Fields(line)
		if fields[1] == mountpoint {
			devPath, err := filepath.EvalSymlinks(fields[0])
			if err != nil {
				return "", err
			}

			return devPath, nil
		}
	}

	return "", ErrNotFound
}

// GetMajorMinor ...
func GetMajorMinor(devPath string) (string, error) {
	devName := filepath.Base(devPath)

	fileInfos, err := ioutil.ReadDir(SysBlockPath)
	if err != nil {
		return "", err
	}

	path := ""
	for _, fileInfo := range fileInfos {
		fileInfoName := fileInfo.Name()
		if strings.HasPrefix(devName, fileInfoName) {
			path = filepath.Join(SysBlockPath, fileInfoName)
			if devName != fileInfoName {
				// name is a partition.
				path = filepath.Join(path, devName)
			}

			path = filepath.Join(path, "dev")
			break
		}
	}

	if path == "" {
		return "", ErrNotFound
	}

	b, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}

	majorMinor := "b" + strings.TrimSpace(string(b))
	return majorMinor, nil
}

// GetUdevData ...
func GetUdevData(majorMinor string) (string, error) {
	b, err := ioutil.ReadFile(filepath.Join(UdevDataPath, majorMinor))
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(b)), nil
}
