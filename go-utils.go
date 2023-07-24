package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"github.com/edsrzf/mmap-go"
	"log"
	"os"
	"os/exec"
	"strconv"
	"syscall"
	"time"
	"unicode"
)

func MMAP_FILE(file_path string, mode string) (*os.File, mmap.MMap, error) {
	mode1 := os.O_RDWR
	var chmod os.FileMode = 0644
	mode2 := mmap.RDWR

	if mode == "ro" {
		mode1 = os.O_RDONLY
		chmod = 0444
		mode2 = mmap.RDONLY
	}
	var err error
	var mmap_handle mmap.MMap
	var file_handle *os.File

	if file_handle, err = os.OpenFile(file_path, mode1, chmod); err != nil {
		return nil, nil, err
	}
	if mmap_handle, err = mmap.Map(file_handle, mode2, 0); err != nil {
		return nil, nil, err
	}
	return file_handle, mmap_handle, nil
} // end func MMAP_FILE

func MMAP_CLOSE(file_path string, file_handle *os.File, mmap_handle mmap.MMap, mode string) (bool, error) {
	if mode == "rw" {
		if err := mmap_handle.Flush(); err != nil {
			log.Printf("ERROR MMAP_CLOSE Flush file_path='%s' err='%v'", file_path, err)
			return false, err
		}
	}

	if err := mmap_handle.Unmap(); err != nil {
		log.Printf("ERROR MMAP_CLOSE Unmap file_path='%s' err='%v'", file_path, err)
		return false, err
	}

	if err := file_handle.Close(); err != nil {
		log.Printf("ERROR MMAP_CLOSE Close file_path='%s' err='%v'", file_path, err)
		return false, err
	}

	//log.Printf("File closed OK fp='%s'", file_path)
	return true, nil
} // end func MMAP_CLOSE

func Line_isPrintable(line string) bool {
	for _, char := range line {
		if !unicode.IsPrint(char) {
			return false
		}
	}
	return true
} // end func Line_isPrintable

func SoftLink(src string, dst string) bool {
	cmd := exec.Command("/bin/ln", "-sf", src, dst)
	//cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		log.Printf("softlink failed src='%s' dst='%s' err='%v'", src, dst, err)
		return false
	}
	return true
} // end func SoftLink

func GetSoftLinkTarget(link string) string {
	if out, err := exec.Command("/bin/readlink", "-f", link).CombinedOutput(); err == nil {
		retstr := string(out)
		if retstr != "" {
			log.Printf("OK GSLT link='%s' retstr='%v'", link, retstr)
			return retstr
		}
	} else {
		log.Printf("ERROR GetSoftLinkTarget link='%s' err='%v'", link, err)
	}
	return ""
} // end func SoftLink

func HardLink(src string, dst string) bool {
	cmd := exec.Command("/bin/ln", src, dst)
	//cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		log.Printf("hardlink failed src='%s' dst='%s' err='%v'", src, dst, err)
		return false
	}
	return true
} // end func HardLink

func DirExists(dir string) bool {
	//log.Printf("?DirExists dir='%s'", dir)
	info, err := os.Stat(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return false
		}
		log.Printf("ERROR DirExists err='%v'", err)
		return false
	}
	return info.IsDir()
} // end func DirExists

func FileExists(File_path string) bool {
	info, err := os.Stat(File_path)
	if err != nil {
		if os.IsNotExist(err) {
			return false
		}
		log.Printf("ERROR FileExists err='%v'", err)
		return false
	}
	return !info.IsDir()
} // end func fileExists

func Mkdir(dir string) bool {
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		log.Printf("ERROR Mkdir='%s' err='%v'", dir, err)
		return false
	} else {
		//log.Printf("CREATED DIR %s", dir)
		return true
	}
	return false
} // end func Mkdir

func Hash256(astr string) string {
	ahash := sha256.Sum256([]byte(astr))
	return hex.EncodeToString(ahash[:])
} // end func hash256

func IsDigit(s string) bool {
	for _, r := range s {
		if !unicode.IsDigit(r) {
			return false
		}
	}
	return true
} // end func is Digit

func IsSpace(b byte) bool {
	return b < 33
}

func unprintable(b byte) bool {
	return b < 32
}

func Str2int(str string) int {
	if IsDigit(str) {
		aint, err := strconv.Atoi(str)
		if err == nil {
			return aint
		}
	}
	return 0
} // end func str2int

func Str2int64(str string) int64 {
	if IsDigit(str) {
		aint64, err := strconv.ParseInt(str, 10, 64)
		if err == nil {
			return aint64
		}
	}
	return 0
} // end func str2int64

func Str2uint64(str string) uint64 {
	if IsDigit(str) {
		auint64, err := strconv.ParseUint(str, 10, 64)
		if err == nil {
			return auint64
		}
	}
	return 0
} // end func str2int64

func Lines2Bytes(lines []string) []byte {
	var buf []byte
	for _, line := range lines {
		buf = append(buf, []byte(line+"\n")...)
	}
	return buf
} // end func Lines2Bytes

// linux syscall for Fallocate
func Fallocate(file *os.File, offset int64, length int64) error {
	if length == 0 {
		return nil
	}
	return syscall.Fallocate(int(file.Fd()), 0, offset, length)
} // end func Fallocate

func Now() int64 {
	return UnixTimeSec()
}

func Nano() int64 {
	return UnixTimeNanoSec()
}

func UnixTimeSec() int64 {
	return time.Now().UnixNano() / 1e9
} // end func Now

func UnixTimeMilliSec() int64 {
	return time.Now().UnixNano() / 1e6
} // end func Milli

func UnixTimeMicroSec() int64 {
	return time.Now().UnixNano() / 1e3
} // end func Micro

func UnixTimeNanoSec() int64 {
	return time.Now().UnixNano()
} // end func Nano

func BootSleep() {
	time.Sleep(100 * time.Microsecond)
} // end func BootSleep

func DebugSleepS(sec int) {
	time.Sleep(time.Duration(sec) * time.Second)
} // end func DebugSleepS

func DebugSleepM(microsec int) {
	time.Sleep(time.Duration(microsec) * time.Microsecond)
} // end func DebugSleepM
