package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	jwtmiddleware "github.com/auth0/go-jwt-middleware"
	"github.com/codegangsta/negroni"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/joho/godotenv"
	"github.com/rs/cors"
	"gorm.io/datatypes"
)

type Response struct {
	Message string `json:"message"`
}

type Jwks struct {
	Keys []JSONWebKeys `json:"keys"`
}

type JSONWebKeys struct {
	Kty string   `json:"kty"`
	Kid string   `json:"kid"`
	Use string   `json:"use"`
	N   string   `json:"n"`
	E   string   `json:"e"`
	X5c []string `json:"x5c"`
}

type player struct {
	Nickname datatypes.JSON
	SteamID  datatypes.JSON
}

type Team struct {
	gorm.Model
	TeamName datatypes.JSON
	Players  datatypes.JSON
}

type Match struct {
	Team1       Team `json:"Team`
	Team2       Team `json:"Team`
	Score_team1 int  `json:"Score_team1"`
	Score_team2 int  `json:"Score_team2"`
}
type Tournament struct {
	gorm.Model
	TournamentName datatypes.JSON
	Size           int    `json:"Size"`
	Teams          []Team `json:"Teams"`
}

type Tournaments struct {
	gorm.Model
	TournamentList []Tournament `json:"TournamentList"`
}

var allTeams []Team

func run() {

}

type CustomClaims struct {
	Scope string `json:"scope"`
	jwt.StandardClaims
}

func checkScope(scope string, tokenString string) bool {
	token, _ := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		cert, err := getPemCert(token)
		if err != nil {
			return nil, err
		}
		result, _ := jwt.ParseRSAPublicKeyFromPEM([]byte(cert))
		return result, nil
	})

	claims, ok := token.Claims.(*CustomClaims)
	hasScope := false
	if ok && token.Valid {
		result := strings.Split(claims.Scope, " ")
		for i := range result {
			if result[i] == scope {
				hasScope = true

			}
		}
	}

	return hasScope
}

func getPemCert(token *jwt.Token) (string, error) {
	cert := ""
	resp, err := http.Get(os.Getenv("KEYCLOAK_CERT"))
	if err != nil {
		return cert, err
	}
	defer resp.Body.Close()

	var jwks = Jwks{}
	err = json.NewDecoder(resp.Body).Decode(&jwks)
	if err != nil {
		return cert, err
	}
	for k, _ := range jwks.Keys {
		if token.Header["kid"] == jwks.Keys[k].Kid {
			cert = "-----BEGIN CERTIFICATE-----\n" + jwks.Keys[k].X5c[0] + "\n-----END CERTIFICATE-----"
		}
	}

	if cert == "" {
		err := errors.New("Unable to find appropriate key.")
		return cert, err
	}
	return cert, nil
}

func responseJSON(message string, w http.ResponseWriter, statusCode int) {
	response := Response{message}

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(jsonResponse)
}

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Print("Error loading .env file")
	}
	var db *gorm.DB

	psqlConnectString := os.Getenv("psqlConnectString")
	db, err = gorm.Open("postgres", psqlConnectString)
	if err != nil {

		panic("failed to connect database")

	}

	defer db.Close()

	jwtMiddleware := jwtmiddleware.New(jwtmiddleware.Options{
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			// Verify 'aud' claim
			aud := "account"
			checkAud := token.Claims.(jwt.MapClaims).VerifyAudience(aud, false)
			if !checkAud {
				return token, errors.New("Invalid audience.")
			}
			// Verify 'iss' claim
			iss := "http://10.152.183.116:8080/auth/realms/example"
			checkIss := token.Claims.(jwt.MapClaims).VerifyIssuer(iss, false)
			if !checkIss {
				return token, errors.New("Invalid issuer.")
			}

			cert, err := getPemCert(token)
			if err != nil {
				panic(err.Error())
			}
			result, _ := jwt.ParseRSAPublicKeyFromPEM([]byte(cert))
			return result, nil
		},
		SigningMethod: jwt.SigningMethodRS256,
	})

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://10.0.0.32", "http://10.0.0.32:8080"},
		AllowCredentials: true,
		AllowedHeaders:   []string{"Authorization"},
	})

	r := mux.NewRouter()

	// This route is always accessible
	r.Handle("/api/getAllTournaments", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var allTournaments Tournaments
		db.Find(&allTournaments.TournamentList)
		json.NewEncoder(w).Encode(allTournaments)
	}))

	/*
		methods for team administration:
		listTeams
		addTeam
		RemoveTeam
		requestToJoinTeam
		getRequestedJoiners
	*/

	// This route is only accessible if the user has a valid access_token
	// We are chaining the jwtmiddleware middleware into the negroni handler function which will check
	// for a valid token.
	r.Handle("/api/listTeams", negroni.New(
		negroni.HandlerFunc(jwtMiddleware.HandlerWithNext),
		negroni.Wrap(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var allTeams []Team
			// var allPlayers []player
			db.Find(&allTeams)
			if len(allTeams) == 0 {
				message := "no teams registered at this time"
				json.NewEncoder(w).Encode(message)
				json.NewEncoder(w).Encode(allTeams)

			} else {
				json.NewEncoder(w).Encode(allTeams)
			}
		}))))

	// This route is only accessible if the user has a valid access_token and a scope
	// We are chaining the jwtmiddleware middleware into the negroni handler function which will check
	// for a valid token and scope.
	r.Handle("/api/addTeam", negroni.New(
		negroni.HandlerFunc(jwtMiddleware.HandlerWithNext),
		negroni.Wrap(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeaderParts := strings.Split(r.Header.Get("Authorization"), " ")
			token := authHeaderParts[1]

			hasScope := checkScope("tournamentmanager", token)

			if !hasScope {
				message := "Insufficient scope."
				responseJSON(message, w, http.StatusForbidden)
				return
			}
			var team Team
			decoder := json.NewDecoder(r.Body)
			err := decoder.Decode(&team)
			if err != nil {
				responseJSON(err.Error(), w, http.StatusOK)
				fmt.Println(err.Error())
			} else {
				responseJSON(fmt.Sprintf("Registered team: %s", team.TeamName), w, http.StatusOK)
				db.Create(&team)

			}
		}))))

	r.Handle("/api/removeTeam", negroni.New(
		negroni.HandlerFunc(jwtMiddleware.HandlerWithNext),
		negroni.Wrap(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeaderParts := strings.Split(r.Header.Get("Authorization"), " ")
			token := authHeaderParts[1]

			hasScope := checkScope("tournamentmanager", token)

			if !hasScope {
				message := "Insufficient scope."
				responseJSON(message, w, http.StatusForbidden)
				return
			}
			var team Team
			decoder := json.NewDecoder(r.Body)
			err := decoder.Decode(&team)
			if err != nil {
				responseJSON(err.Error(), w, http.StatusOK)
				fmt.Println(err.Error())
			} else {
				db.Delete(&team)
				responseJSON(fmt.Sprintf("Deleted team: %s", &team.TeamName), w, http.StatusOK)

			}
		}))))

	/* tournamentEndpoints
	getAllTournaments
	addTournaments
	getTournament
	deleteTournament
	addTeamToTournament

	*/

	// This route is always accessible
	r.Handle("/api/getAllTournaments", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var allTournaments Tournaments
		db.Find(&allTournaments.TournamentList)
		json.NewEncoder(w).Encode(allTournaments)
	}))

	r.Handle("/api/getTournament", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		var tournament Tournament
		if r.PostForm["TournamentID"] != nil {
			db.First(&tournament, "id = ?", r.PostForm["TournamentID"])
			json.NewEncoder(w).Encode(tournament)

		}
	}))

	r.Handle("/api/addTournaments", negroni.New(
		negroni.HandlerFunc(jwtMiddleware.HandlerWithNext),
		negroni.Wrap(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeaderParts := strings.Split(r.Header.Get("Authorization"), " ")
			token := authHeaderParts[1]

			hasScope := checkScope("tournamentmanager", token)

			if !hasScope {
				message := "Insufficient scope."
				responseJSON(message, w, http.StatusForbidden)
				return
			}
			var Tournaments Tournaments
			decoder := json.NewDecoder(r.Body)
			err := decoder.Decode(&Tournaments)
			if err != nil {
				responseJSON(err.Error(), w, http.StatusOK)
				fmt.Println(err.Error())
			}
			var addedTournaments []Tournament
			for _, tournament := range Tournaments.TournamentList {
				db.Create(&tournament)
				addedTournaments = append(addedTournaments, tournament)
			}
			json.NewEncoder(w).Encode(addedTournaments)

		}))))

	r.Handle("/api/deleteTournament", negroni.New(
		negroni.HandlerFunc(jwtMiddleware.HandlerWithNext),
		negroni.Wrap(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeaderParts := strings.Split(r.Header.Get("Authorization"), " ")
			token := authHeaderParts[1]

			hasScope := checkScope("tournamentmanager", token)

			if !hasScope {
				message := "Insufficient scope."
				responseJSON(message, w, http.StatusForbidden)
				return
			}
			r.ParseForm()
			var tournament Tournament
			if r.PostForm["TournamentID"] != nil {
				db.Where(&tournament, "id = ?", r.PostForm["TournamentID"]).Delete(&tournament)
				json.NewEncoder(w).Encode(tournament)
			}

		}))))

	r.Handle("/api/addTeamToTournament", negroni.New(
		negroni.HandlerFunc(jwtMiddleware.HandlerWithNext),
		negroni.Wrap(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeaderParts := strings.Split(r.Header.Get("Authorization"), " ")
			token := authHeaderParts[1]

			hasScope := checkScope("tournamentmanager", token)

			if !hasScope {
				message := "Insufficient scope."
				responseJSON(message, w, http.StatusForbidden)
				return
			}
			r.ParseForm()
			var tournament Tournament
			if r.PostForm["TournamentID"] != nil {
				db.Where(&tournament, "id = ?", r.PostForm["TournamentID"]).Create(&tournament.Teams)
				json.NewEncoder(w).Encode(tournament.Teams)
			}

		}))))

	handler := c.Handler(r)
	http.Handle("/", r)
	fmt.Println("Listening on http://localhost:27015")
	http.ListenAndServe("0.0.0.0:27015", handler)
}
