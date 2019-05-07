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
)

const (
	// TokensURL is the URL for getting an ID token
	TokensURL = "https://api.pencildata.com/token"
	// RegisterURL is the URL for a document registration
	RegisterURL = "https://api.pencildata.com/register/"
	// storage is private (default)
	Storage = "private"
)

// UserInfo represents a struct for a user info needed for auth
type UserInfo struct {
	Name     string `json:"username"`
	Password string `json:"password"`
}

// AuthResponse represents a struct for response of auth endpoint. Here only needed fields are parsed.
// Entire response struct is presented in API docs.
type AuthResponse struct {
	Result struct {
		AccessToken string `json:"accessToken"`
	} `json:"result"`
}

// RegisterReqBody represents a struct for request body for Register an entity
type RegisterReqBody struct {
	Hash string `json:"hash"`
	Storage string   `json:"storage"`
}

// GetAccessToken authorizes user by passed username and password, returns accessToken needed for future requests
func GetAccessToken(userInfo UserInfo) (string, error) {
	// marshal userInfo object to json
	userData, err := json.Marshal(userInfo)
	if err != nil {
		return "", err
	}

	// create a new http request object
	req, err := http.NewRequest(http.MethodPost, TokensURL, bytes.NewBuffer(userData))
	if err != nil {
		return "", err
	}
	// set the Content-Type header for the request
	req.Header.Set("Content-Type", "application/json")

	// create a new http client object. Notice: go’s http package doesn’t specify request timeouts by default
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

	// ensure that http status equals 200
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("Unable to auth user. StatusCode: %d Body: %s", resp.StatusCode, body)
	}

	authResp := AuthResponse{}
	// parse the response body to appropriate struct
	if err := json.Unmarshal(body, &authResp); err != nil {
		return "", err
	}
	// accessToken needed for an authorization in future requests
	return authResp.Result.AccessToken, nil
}

// PrepareFile reads the entire file by passed file path and returns the SHA256 checksum of the data
func PrepareFile(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	// close the file when you’re done
	defer f.Close()

	h := sha256.New()
	// copy the entire file data
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}

	// return a hexadecimal representation of the checksum
	return hex.EncodeToString(h.Sum(nil)), nil
}

// Register registers the hash. Returns an entity id for the registered hash once it has been stored in the blockchain
func Register(hash, token string) (string, error) {
	// marshal request object with passed hash to json
	reqData, err := json.Marshal(RegisterReqBody{Hash: hash,Storage: storage})
	if err != nil {
		return "", err
	}

	// create a new http request object
	req, err := http.NewRequest(http.MethodPost, RegisterURL, bytes.NewBuffer(reqData))
	if err != nil {
		return "", err
	}
	// set the Content-Type header
	req.Header.Set("Content-Type", "application/json")
	// set an accessToken to the Authorization header: needed to be authorized to access the apis.
	req.Header.Set("Authorization", "Bearer "+token)

	// create a new http client object
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
	// ensure that http status equals 200
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("Unable to register a document. StatusCode: %d Body: %s", resp.StatusCode, body)
	}

	// entity id for the registered hash
	return string(body), nil
}

func main() {
	// build an userInfo object with your credentials
	userInfo := UserInfo{
		Name:     "xxx",
		Password: "xxx",
	}

	// authorize the user, get an accessToken for a register request
	token, err := GetAccessToken(userInfo)
	if err != nil {
		// handle an error
		log.Fatalln(err)
	}

	// set your file path here and get a checksum for the file
	filePath := "file.txt"
	hash, err := PrepareFile(filePath)
	if err != nil {
		log.Fatalln(err)
	}
	// register the hash and get an entity id once it has been stored in the blockchain
	entityID, err := Register(hash, token)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println("Registered Entity ID = ", entityID)
}
