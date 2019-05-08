package main
import (
    "bytes"
    "crypto/sha256"
    "encoding/hex"
    "encoding/json"
    "fmt"
    "io"
    "io/ioutil"
    "log"
    "net/http"
    "os"
    "strconv"
    "golang.org/x/crypto/ssh/terminal"
    "syscall"
)
const (
    // TokensURL is the URL for getting an ID token
    TokensURL = "https://api.pencildata.com/token"
    // VerifyURL is the URL for a document verification
    VerifyURL = "https://api.pencildata.com/verify/"
)
// UserInfo represents a struct for a user info needed for auth
type UserInfo struct {
    Name     string `json:"userId"`
    Password string `json:"password"`
}
// AuthResponse represents a struct for response of auth endpoint. Here only needed fields are parsed.
// Entire response struct is presented in API docs.
type AuthResponse struct {
    Data struct {
        AccessToken string `json:"accessToken"`
        ExpiresIn string `json:"expiresIn"`
        RefreshToken string `json:"refreshToken"`
    } `json:"data"`
}
// GetAccessToken authorizes user by passed username and password, returns accessToken needed for future requests
func GetToken(userInfo UserInfo) (string, error) {
    // marshal userInfo object to json
    userData, err := json.Marshal(userInfo)
    if err != nil {
        return "", err
    }
    // create a new http request object
    req, err := http.NewRequest("POST", TokensURL, bytes.NewBuffer(userData))
    if err != nil {
        return "", err
    }
    // set the Content-Type header for the request
    req.Header.Set("Content-Type", "application/json")
    // create a new http client object. Notice: go's http package doesn't specify request timeouts by default
    client := &http.Client{}
    // make the request
    resp, err := client.Do(req)
    if err != nil {
        return "", err
    }
    // the client must close the response body when finished with it
    defer resp.Body.Close()
    // ready the body content
    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return "", err
    }
    fmt.Println("Status Code: ", resp.Status)
    // ensure that http status equals 200
    if resp.StatusCode != http.StatusOK {
        return "", fmt.Errorf("Unable to auth user. StatusCode: %d Body: %s", resp.StatusCode, body)
    }
    authResp := AuthResponse{}
    // // parse the response body to appropriate struct
    if err := json.Unmarshal(body, &authResp); err != nil {
        return "", err
    }
    // fmt.Println(authResp.Data.AccessToken)
    return authResp.Data.AccessToken, nil
}
// PrepareFile reads the entire file by passed file path and returns the SHA256 checksum of the data
func PrepareFile(filename string) (string, error) {
    // fmt.Println(os.Getwd())
    // fmt.Println(filename)
    f, err := os.Open(filename)
    // fmt.Println("hello")
    if err != nil {
        return "", err
    }
    // close the file when you're done
    defer f.Close()
    h := sha256.New()
    // copy the entire file data
    if _, err := io.Copy(h, f); err != nil {
        return "", err
    }
    // return a hexadecimal representation of the checksum
    fmt.Println(hex.EncodeToString(h.Sum(nil)))
    return hex.EncodeToString(h.Sum(nil)), nil
}
// Verify checks whether or not hash is the same as the stored hash for entityID
func Verify(entityID, hash, token string, storage string) (bool, error) {
    // create a new http request object
    req, err := http.NewRequest(http.MethodGet, VerifyURL+entityID, nil)
    if err != nil {
        return false, err
    }
    // set the Content-Type header
    req.Header.Set("Content-Type", "application/json")
    // set an accessToken to the Authorization header: needed to be authorized to access the apis.
    req.Header.Set("Authorization", "Bearer " + token)
    // add the "hash" for document verification to query params
    q := req.URL.Query()
    q.Add("hash", hash)
    q.Add("storage", "private")
    req.URL.RawQuery = q.Encode()
    // create a new http client object
    client := &http.Client{}
    // make the request
    resp, err := client.Do(req)
    if err != nil {
        return false, err
    }
    // the client must close the response body when finished with it
    defer resp.Body.Close()
    // ready the body content
    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return false, err
    }
    // ensure that http status equals 200
    if resp.StatusCode != http.StatusOK {
        return false, fmt.Errorf("Unable to verify a document. StatusCode: %d Body: %s", resp.StatusCode, body)
    }
    // parse the response body as boolean value
    valid, err := strconv.ParseBool(string(body))
    if err != nil {
        return false, err
    }
    return valid, nil
}
func main() {
    fmt.Print("username: ")
    var username string
    fmt.Scanln(&username)
    fmt.Printf("password: ")
    // Silent. For printing *'s use gopass.GetPasswdMasked()
    bytePassword, err := terminal.ReadPassword(int(syscall.Stdin))
    if err == nil {
        //fmt.Println("\nPassword typed: " + string(bytePassword))
    }
    password := string(bytePassword)
    fmt.Println("")
    userInfo := UserInfo{
        Name:     username,
        Password: password,
    }
    fmt.Print("storage: ")
    var storage string
    fmt.Scanln(&storage)
    fmt.Print("assetId: ")
    var assetId string
    fmt.Scanln(&assetId)
    // First element in os.Args is always the program name,
    // So we need at least 2 arguments to have a file name argument.
    if len(os.Args) < 2 {
        fmt.Println("Missing parameter, provide file name!")
        return
    }
    filename := os.Args[1]
    // set your file path here and get a checksum for the file
    hash, err := PrepareFile(filename)
    if err != nil {
        log.Fatalln(err)
    }
    // authorize the user, get an accessToken for a register request
    token, err := GetToken(userInfo)
    if err != nil {
        // handle an error
        log.Fatalln(err)
    }
    //fmt.Println(token)
    // verify whether or not hash is the same as the stored hash for entityID
    valid, err := Verify(assetId, hash, token, storage)
    if err != nil {
        log.Fatalln(err)
    }
    fmt.Println("Got verification result = ", valid)
}
