package main

import (
    "encoding/json"
    "io/ioutil"
    "log"
    "net/http"
    "os/exec"
    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/service/s3/s3manager"
    "fmt"
    "strconv"
    "os"
)

// to parse config json
type Config struct {
    ScanWorkflow string `json:"scan_workflow"`
    // include other properties as required
}

// for the response from the url
type Response struct {
    Programs []Program `json:"programs"`
}

// to parse programs array
type Program struct {
    Domains []string `json:"domains"`
}

func main() {
    // config := loadConfig("config_path_here")
    // fmt.Println(config)
    fmt.Println("FUNCTION STARTED")
    list := loadChaos()
    if list == nil {
        log.Fatal("Failed to load chaos")
        return
    }

    // making slice with capacity as the length of list
    domains := make([]string, 0, len(list))
    for _, p := range list {
        domains = append(domains, p.Domains...)
    }

    fileCount := 0
    fmt.Println("STARTING TO WRITE FILES : ")
    fmt.Println()
    for _, domain := range domains {
       fileCount++
       fileName := "domains" + strconv.Itoa(fileCount) + ".txt"
       fmt.Println()
       fmt.Println("Starting to write File : ", fileName)
       err := ioutil.WriteFile(fileName, []byte(domain+"\n"), 0644)
       if err != nil {
          panic(err)
       }
       fmt.Println("File : ", fileName, "Successfully Written")
    }
    fmt.Println()
    for i := 1; i <= fileCount; i++ {
       outputFile := "output" + strconv.Itoa(i) + ".json"
       inputFile := "domains" + strconv.Itoa(i) + ".txt"
       fmt.Println()
       fmt.Println("Starting Nuclei Execution for : ", inputFile)
       cmd := exec.Command("nuclei", "-l", inputFile, "-v", "-je", outputFile)
       output, err := cmd.Output()
       if err != nil {
                log.Fatal("Error running command: %v", err)
        }
       fmt.Println("Nuclei Output : ", output)

       // Your S3 bucket name you wish to upload to
       UploadToS3("scanningreports", outputFile)
       fmt.Println("Uploaded file  : ", outputFile, "to S3")
        }
}

func loadConfig(path string) Config {
    b, err := ioutil.ReadFile(path)
    if err != nil {
        log.Fatal("Failed to load config file")
    }

    var conf Config
    json.Unmarshal(b, &conf)

    if conf.ScanWorkflow == "*" {
        conf.ScanWorkflow = ""
    }
    return conf
}

func loadChaos() []Program {
    resp, err := http.Get("https://raw.githubusercontent.com/projectdiscovery/public-bugbounty-programs/main/chaos-bugbounty-list.json")
    if err != nil {
        log.Fatal("Failed to load url")
        return nil
    }
    defer resp.Body.Close()

    var r Response
    json.NewDecoder(resp.Body).Decode(&r)

    return r.Programs
}

func UploadToS3(bucket, filename string) error {
    // Create a session using your current aws configuration
    sess, err := session.NewSession(&aws.Config{
        Region: aws.String("us-west-2")}, // change to your desired region Eg. us-west-2
    )

    if err != nil {
        return log.Fatal("Failed to create a session, %v", err)
    }

    // Create a new uploader using session
    uploader := s3manager.NewUploader(sess)

    // Open the file for upload
    file, err := os.Open(filename)
    if err != nil {
        return log.Fatal("Failed to open file %q, %v", filename, err)
    }

    // Upload the file, returns an UploadOutput struct and any error
    _, err = uploader.Upload(&s3manager.UploadInput{
        Bucket: aws.String(bucket),
        Key:    aws.String(filename),
        Body:   file,
    })

    // Close the file when done
    defer file.Close()

    if err != nil {
        // Print the error and exit.
        return log.Fatal("Unable to upload %q to %q, %v", filename, bucket, err)
    }

    fmt.Printf("Successfully uploaded %q to %q\n", filename, bucket)
    return nil
}