/*
 * go4api - an api testing tool written in Go
 * Created by: Ping Zhu 2018
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
 *
 */

package utils

import (
    "fmt"
    "io/ioutil"                                                                                                                                              
    "os"
    "os/user"
    "io"
    "strings"
    "strconv"
    "path/filepath"
    "encoding/csv"
    "encoding/base64"
)

func GetCurrentDir() string{
    // get current dir, 
    // Note: here may be a bug if run the main.go on other path, need to use abs path
    currentDir, err := os.Getwd()
    if err != nil {
        panic(err)
    }

    return currentDir
}

func GetCsvFromFile(filePath string) [][]string {
    fi, err := ioutil.ReadFile(filePath)
    if err != nil {
        fmt.Println("!! Error: ", filePath)
        panic(err)
    }

    reader := csv.NewReader(strings.NewReader(string(fi)))
    // the csv may contain row for comments, which is leading with #, to ignore
    reader.Comment = '#'

    csvRows, err := reader.ReadAll()
    if err != nil {
        fmt.Println("!! Error: ", filePath)
        panic(err)
    }

    return csvRows
}

func GetContentFromFile(filePath string) []byte {
    fi,err := ioutil.ReadFile(filePath)
    if err != nil {
        panic(err)
    }
    // contents := strings.NewReader(string(fi))

    return fi
}

func GetJsonFromFile(filePath string) string {
    fi, err := os.Open(filePath)
    if err != nil {
        panic(err)
    }
    defer fi.Close()
    
    fd, err := ioutil.ReadAll(fi)
    if err != nil {
        panic(err)
    }

    return string(fd)
}

// for path contains ~
func GetAbsPath (path string) string {
    var absPath string

    u, _ := user.Current()

    if path == "~" {
        absPath = u.HomeDir
    } else if strings.HasPrefix(path, "~/") {
        absPath = filepath.Join(u.HomeDir, path[2:])
    } else {
        absPath = path
    }

    return absPath
}

func GetAbsPaths (paths []string) []string {
    var absPaths []string

    for _, path := range paths {
        u, _ := user.Current()

        if path == "~" {
            p := u.HomeDir
            absPaths = append(absPaths, p)
        } else if strings.HasPrefix(path, "~/") {
            p := filepath.Join(u.HomeDir, path[2:])
            absPaths = append(absPaths, p)
        } else {
            absPaths = append(absPaths, path)
        }
    }

    return absPaths
}

//
func GetOsEnviron () map[string]string {
    envMap := map[string]string{}
    var envArray []string

    envArray = os.Environ()
    for _, env := range envArray {
        // find out the first = position, to get the key
        env_k := strings.Split(env, "=")[0]
        envMap[env_k] = os.Getenv(env_k)
    }
    // incase the env variables have the prefix "go4_*"
    // for _, env := range envArray {
    //     // find out the first = position, to get the key
    //     env_k := strings.Split(env, "=")[0]
    //     if strings.HasPrefix(env_k, "go4_") {
    //         if strings.TrimLeft(env_k, "go4_") != "" {
    //             envMap[strings.TrimLeft(env_k, "go4_")] = os.Getenv(env_k)
    //         }
    //     } 
    // }

    return envMap
}

// for the dir and sub-dir
func WalkPath(searchDir string, extension string) ([]string, error) {
    fileList := make([]string, 0)

    e := filepath.Walk(searchDir, func(subPath string, f os.FileInfo, err error) error {
        if filepath.Ext(subPath) == extension {
            fileList = append(fileList, subPath)
        }
        return err
    })
    
    if e != nil {
        fmt.Println("Err, search files under ", searchDir, " failed")
        panic(e)
    }

    // for _, file := range fileList {
    //     fmt.Println(file)
    // }
    return fileList, nil
}

func FileCopy(src string, dest string, info os.FileInfo) error {
    f, err := os.Create(dest)
    if err != nil {
      return err
    }
    defer f.Close()

    if err = os.Chmod(f.Name(), info.Mode())
    err != nil {
      return err
    }

    s, err := os.Open(src)
    if err != nil {
      return err
    }
    defer s.Close()

    _, err = io.Copy(f, s)
    return err
  }


func DirCopy(src string, dest string, info os.FileInfo) error {
    if err := os.MkdirAll(dest, info.Mode())
    err != nil {
      return err
    }

    infos, err := ioutil.ReadDir(src)
    if err != nil {
      return err
    }

    for _, info := range infos {
      if err := FileCopy(filepath.Join(src, info.Name()), filepath.Join(dest, info.Name()), info) 
      err != nil {
        return err
      }
    }

    return nil
}

func ConvertIntArrayToStringArray(intArray []int) []string {
    var stringArray []string
    for _, k := range intArray{
        ii := strconv.Itoa(k)
        stringArray = append(stringArray, ii)
    }

    return stringArray
}

func ConvertStringArrayToIntArray(stringArray []string) []int {
    var intArray []int
    for _, k := range stringArray{
        ii, _ := strconv.Atoi(k)
        intArray = append(intArray, ii)
    }

    return intArray
}

func GenerateFileBasedOnVarAppend(strVar string, filePath string) {
    outFile, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
    if err != nil {
       panic(err) 
    }

    defer outFile.Close()

    outFile.WriteString(strVar)
}

func GenerateFileBasedOnVarOverride(strVar string, filePath string) {
    outFile, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
    if err != nil {
       panic(err) 
    }
    defer outFile.Close()

    outFile.WriteString(strVar)
}


func GenerateCsvFileBasedOnVarAppend(strVarSlice []string, filePath string) {
    outFile, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
    if err != nil {
       panic(err) 
    }

    defer outFile.Close()
    // UTF-8 BOM
    // outFile.WriteString("\xEF\xBB\xBF")
    w := csv.NewWriter(outFile)
    w.Write(strVarSlice)
    // 
    w.Flush()
}

func GenerateCsvFileBasedOnVarOverride(strVarSlice []string, filePath string) {
    outFile, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
    if err != nil {
       panic(err) 
    }
    defer outFile.Close()
    // UTF-8 BOM
    // outFile.WriteString("\xEF\xBB\xBF")
    w := csv.NewWriter(outFile)
    w.Write(strVarSlice)
    // 
    w.Flush()
}

func GeneratePicture(bytesVar []byte, filePath string) {
    outFile, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
    if err != nil {
       panic(err) 
    }
    defer outFile.Close()

    outFile.Write(bytesVar)
}

func DecodeBase64(b64 string) []byte {
    sDec, err := base64.StdEncoding.DecodeString(b64)
    if err != nil {
       panic(err) 
    }

    return sDec
}

func CreateTempDir (filePath string) string {
    err := os.MkdirAll(filepath.Dir(filePath) + "/temp", 0777)
    if err != nil {
      panic(err) 
    }

    return filepath.Dir(filePath) + "/temp"
}

func CheckFilesExistence(fileList []string) bool {
    ifExist := true

    for _, filePath := range fileList {
        _, err := os.Stat(filePath)
        if err != nil {
            ifExist = false
            break
        }
    }

    return ifExist
}

func CheckFileExistence(filePath string) bool {
    ifExist := true

    _, err := os.Stat(filePath)
    if err != nil {
        ifExist = false
    }

    return ifExist
}



