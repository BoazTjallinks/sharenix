/*
   Copyright 2014 Franc[e]sco (lolisamurai@tfwno.gf)
   This file is part of sharenix.
   sharenix is free software: you can redistribute it and/or modify
   it under the terms of the GNU General Public License as published by
   the Free Software Foundation, either version 3 of the License, or
   (at your option) any later version.
   sharenix is distributed in the hope that it will be useful,
   but WITHOUT ANY WARRANTY; without even the implied warranty of
   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
   GNU General Public License for more details.
   You should have received a copy of the GNU General Public License
   along with sharenix. If not, see <http://www.gnu.org/licenses/>.
*/

package sharenixlib

import (
	"errors"
	"fmt"
	"os"
	"os/user"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

// taken from https://github.com/kardianos
func executable() (string, error) {
	switch runtime.GOOS {
	case "linux", "android":
		const deletedTag = " (deleted)"
		execpath, err := os.Readlink("/proc/self/exe")
		if err != nil {
			return execpath, err
		}
		execpath = strings.TrimSuffix(execpath, deletedTag)
		execpath = strings.TrimPrefix(execpath, deletedTag)
		return execpath, nil
	case "netbsd":
		return os.Readlink("/proc/curproc/exe")
	case "dragonfly":
		return os.Readlink("/proc/curproc/file")
	case "solaris":
		return os.Readlink(fmt.Sprintf("/proc/%d/path/a.out", os.Getpid()))
	}
	return "", errors.New("ExecPath not implemented for " + runtime.GOOS)
}

func GetExeDir() (execpath string, err error) {
	execpath, err = executable()
	if err != nil {
		return
	}
	return filepath.Dir(execpath), nil
}

func GetHome() (res string) {
	if res = os.Getenv("HOME"); res != "" {
		return
	}

	usr, err := user.Current()
	if err != nil {
		panic(err)
	}

	res = usr.HomeDir
	return
}

func GetStorageDir() (res string, err error) {
	res = path.Join(GetHome(), "sharenix")
	err = MkDirIfNotExists(res)
	return
}

// GetArchiveDir returns the absolute path to the archive directory.
func GetArchiveDir() (res string, err error) {
	storage, err := GetStorageDir()
	if err != nil {
		return
	}
	res = path.Join(storage, "archive")
	err = MkDirIfNotExists(res)
	return
}

// GenerateArchivedFilename returns an unique file path inside
// archive/ that contains the current date, time and nanotime.
func GenerateArchivedFilename(extension string) (string, error) {
	t := time.Now()
	ye, mo, da := t.Date()
	hour, min, sec := t.Clock()

	archiveDir, err := GetArchiveDir()
	if err != nil {
		return "", err
	}

	return path.Join(archiveDir, fmt.Sprintf("%v-%v-%v_%v-%v-%v_%v.%s",
		ye, int(mo), da, hour, min, sec, t.UnixNano(), extension)), nil
}

// GetHistoryCSV returns the absolute path to the history csv.
func GetHistoryCSV() (res string, err error) {
	storage, err := GetStorageDir()
	if err != nil {
		return
	}
	res = path.Join(storage, "sharenix.csv")
	return
}

// GetPluginsDir returns the absolute path to the plugins directory.
func GetPluginsDir() (res string, err error) {
	storage, err := GetStorageDir()
	if err != nil {
		return
	}
	res = path.Join(storage, "plugins")
	err = MkDirIfNotExists(res)
	return
}
