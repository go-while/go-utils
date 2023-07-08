package utils

import (
    "crypto/sha256"
    "encoding/hex"
    "unicode"
    "strconv"
    "time"
    "os"
    "github.com/edsrzf/mmap-go"
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
    if file_handle, err = os.OpenFile(file_path, mode1, chmod); err == nil {
        if mmap_handle, err = mmap.Map(file_handle, mode2, 0); err == nil {
            return file_handle, mmap_handle, nil
        }
    }
    return nil, nil, err
} // end func MMAP_FILE

func MMAP_CLOSE(file_path string, file_handle *os.File, mmap_handle mmap.MMap, mode string) (retval bool, err error) {
    if mode == "rw" {
        if err = mmap_handle.Flush(); err != nil {
            return
        }
    }
    if err = mmap_handle.Unmap(); err == nil {
        if err = file_handle.Close(); err == nil {
            //log.Printf("File closed OK fp='%s'", file_path)
            retval = true
        }
    }
    return
} // end func MMAP_CLOSE



func Line_isPrintable(line string) bool {
    for _, char := range line {
        if !unicode.IsPrint(char) {
            return false
        }
    }
    return true
} // end func Line_isPrintable

func FileExists(File_path string) bool {
    info, err := os.Stat(File_path)
    if os.IsNotExist(err) {
        return false
    }
    return !info.IsDir()
} // end func fileExists


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

