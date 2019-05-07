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
)

const (
	// TokensURL is the URL for getting an ID token
	TokensURL = "https://api.pencildata.com/token"
	// VerifyURL is the URL for a document verification
	VerifyURL = "https://api.pencildata.com/verify/"
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

// Verify checks whether or not hash is the same as the stored hash for entityID
func Verify(entityID, hash, token string) (bool, error) {
	// create a new http request object
	req, err := http.NewRequest(http.MethodGet, VerifyURL+entityID, nil)
	if err != nil {
		return false, err
	}
	// set the Content-Type header
	req.Header.Set("Content-Type", "application/json")
	// set an accessToken to the Authorization header: needed to be authorized to access the apis.
	req.Header.Set("Authorization", "Bearer "+token)

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

	// set the entity id for a verification
	entityID := "1526156852285"

	// verify whether or not hash is the same as the stored hash for entityID
	valid, err := Verify(entityID, hash, token)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println("Got verification result = ", valid)
}
