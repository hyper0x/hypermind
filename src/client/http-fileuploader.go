
package main

import (
    "bytes"
    "fmt"
    "io"
    "io/ioutil"
    "mime/multipart"
    "net/http"
    "os"
    "flag"
)

func postFile(filename string, targetUrl string) error {
    bodyBuf := &bytes.Buffer{}
    bodyWriter := multipart.NewWriter(bodyBuf)

    fileWriter, err := bodyWriter.CreateFormFile("file", filename)
    if err != nil {
        fmt.Println("error writing to buffer")
        return err
    }

    fh, err := os.Open(filename)
    defer fh.Close()
    if err != nil {
        fmt.Println("error opening file")
        return err
    }

    //iocopy
    _, err = io.Copy(fileWriter, fh)
    if err != nil {
        return err
    }

    contentType := bodyWriter.FormDataContentType()
    bodyWriter.Close()

    resp, err := http.Post(targetUrl, contentType, bodyBuf)
    if err != nil {
        return err
    }
    defer resp.Body.Close()
    resp_body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return err
    }
    fmt.Println(resp.Status)
    fmt.Println(string(resp_body))
    return nil
}

// sample usage 
func main() {
    targetUrl := flag.Arg(0)
    if targetUrl == "" {
        targetUrl = "http://localhost:9090/upload"
    }
    fileFullPath := flag.Arg(1)
    if fileFullPath == "" {
        fileFullPath = "./http-upload.go"
    }
    fmt.Println("Upload file '", fileFullPath, "' to url '", targetUrl, " ...")
    postFile(fileFullPath, targetUrl)
}
