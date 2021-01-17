package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
)

type responsegetNewToken struct {
	AccessToken string `json:"access_token"`
	Tokentype   string `json:"token_type"`
}

var token string = "eyJhbGciOiJSUzI1NiIsInR5cCIgOiAiSldUIiwia2lkIiA6ICJ0T2YyTmRKVXh5X1ZGMmY4OTlnWWVPUWhWNUg1cWZOUmdYRENIS2hSc29nIn0.eyJleHAiOjE2MDk0OTU2MzQsImlhdCI6MTYwOTQ5NTMzNCwiYXV0aF90aW1lIjoxNjA5NDk1MzM0LCJqdGkiOiIyMWE2Y2I4ZC05YWI4LTQyOGYtODZjYi02MjhlOWU5ODI1NmQiLCJpc3MiOiJodHRwOi8vMTAuMTUyLjE4My4xMTY6ODA4MC9hdXRoL3JlYWxtcy9leGFtcGxlIiwiYXVkIjoiYWNjb3VudCIsInN1YiI6IjQ0MDZjNGY0LWFkYjgtNGRjZi1hY2MzLTM3ZjQ1NzViZDM2OCIsInR5cCI6IkJlYXJlciIsImF6cCI6ImpzLWNvbnNvbGUiLCJub25jZSI6ImE0NTMwOGFkLWE3NDMtNDIxYy04MmQyLTRiYzk2N2UxMDI3NCIsInNlc3Npb25fc3RhdGUiOiI4YmViZjhkNy1mNTU4LTRiZjQtOGE4Ny01ZDU2OTcwNmVkMGQiLCJhY3IiOiIxIiwiYWxsb3dlZC1vcmlnaW5zIjpbImh0dHA6Ly8xMC4wLjAuMzI6ODA4MCIsImh0dHA6Ly8xMC4wLjAuMzI6MzAwMSIsImh0dHA6Ly8xMC4wLjAuMzIiXSwicmVhbG1fYWNjZXNzIjp7InJvbGVzIjpbImpzY29uc29sZWFkbWluIiwidXNlciJdfSwicmVzb3VyY2VfYWNjZXNzIjp7ImFjY291bnQiOnsicm9sZXMiOlsidmlldy1wcm9maWxlIl19fSwic2NvcGUiOiJvcGVuaWQgZW1haWwgdG91cm5hbWVudG1hbmFnZXIgcHJvZmlsZSBzdGVhbWlkIiwic3RlYW1pZCI6InllcyIsImVtYWlsX3ZlcmlmaWVkIjpmYWxzZSwibmFtZSI6IkVtaWwgVm9sY2ttYXIgUnkiLCJwcmVmZXJyZWRfdXNlcm5hbWUiOiJlbWlseW8iLCJnaXZlbl9uYW1lIjoiRW1pbCBWb2xja21hciIsImZhbWlseV9uYW1lIjoiUnkiLCJlbWFpbCI6Im15ZW1haWxAZm9vLmJhciJ9.fDFznnVf7kuwZ-0Ra2DDwcR74DIJd2HFwExBZ3x9adOLeuq0RY8eLA6Gw9hhfidECKpMH13Pq2ABXNVZiAqk5DfBxzuNbiZsiG4WiSQiFFyeJlCzLpCBhRS1sqCpPm_CA-XwbvGl5LlxmqJO2QSFttZ8rftHQmWMMurZsBTIg8Wqek20tkEXAhKclcehDWBM8x4OfGXq9W4DngRcwuO6NI3CJoGRS445am0WdEuR2BcD8RHtZEVxrJuElvg8WVBDBNOTkMRqRaAGDp1um_sz08A2tvBi3SlBuyfkkN7I014i2HIGw4Lb_icG-xnKMMz4gh_ru-r3oPH8g0Jldg6DYg"

func registerTeam() {
	url := "http://localhost:27015/api/addTeam"
	// f, err := os.Open("example_match.json.bak")
	f, err := os.Open("compdiffteam.json")

	if err != nil {
		// handle err
	}

	req, err := http.NewRequest("GET", url, f)
	if err != nil {
		fmt.Println(err)
	}

	req.Header.Add("authorization", "Bearer "+token)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println(err)
	}

	defer res.Body.Close()

	body, _ := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(body))

}

func deleteTeam() {
	url := "http://localhost:27015/api/removeTeam"
	f, err := os.Open("compdiffteam.json")
	if err != nil {
		// handle err
	}

	req, err := http.NewRequest("GET", url, f)
	if err != nil {
		fmt.Println(err)
	}

	req.Header.Add("authorization", "Bearer "+token)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println(err)
	}

	defer res.Body.Close()

	body, _ := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(body))

}

func readTeams() {
	url := "http://localhost:27015/api/listTeams"

	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("authorization", "Bearer "+token)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	fmt.Println(string(body))

}

func readMatch() {

}

func makeTournament() {
	url := "http://localhost:27015/api/addTournaments"
	f, err := os.Open("addTournament.json")
	if err != nil {
		// handle err
	}

	req, err := http.NewRequest("GET", url, f)
	if err != nil {
		fmt.Println(err)
	}

	req.Header.Add("authorization", "Bearer "+token)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println(err)
	}

	defer res.Body.Close()

	body, _ := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(body))
}

func getAllTournaments() {
	url := "http://localhost:27015/api/getAllTournaments"

	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("authorization", "Bearer "+token)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	fmt.Println(string(body))

}

func getSpecificTournament(TournamentID string) {
	resp, err := http.PostForm("http://localhost:27015/api/getTournament",
		url.Values{"TournamentID": {TournamentID}})
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	fmt.Println(string(body))
}

func getAllTeams() {
	url := "http://localhost:27015/api/getAllTeams"

	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("authorization", "Bearer "+token)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	fmt.Println(string(body))

}

func main() {
	// getValidToken()
	// registerTeam()
	// deleteTeam()
	// readTeams()

	// makeTournament()
	// getAllTournaments()
	getSpecificTournament("72")
}
