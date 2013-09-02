package user

import (
    "net/http"
    "strings"
    "fmt"
)

func init() {
    http.HandleFunc("/user", getUsers)
    http.HandleFunc("/user/", getUser)
    http.HandleFunc("/user/*/", errorResp)
}

func getUsers(response http.ResponseWriter, request *http.Request) {
	fmt.Fprintf(response, "Users")	

}

func getUser(response http.ResponseWriter, request *http.Request) {
	fmt.Fprintf(response, "Get user %v", strings.Replace(request.RequestURI, "/user/","",1))
}

func errorResp(response http.ResponseWriter, request *http.Request) {
	fmt.Fprintf(response, "error")
}

