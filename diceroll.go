package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"html"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

// DiceRolls is the returntype for the REST-api
type DiceRolls struct {
	Status string `json:"Status"` // Status of the call, error string or "ok"
	Result []int  `json:"Result"` // THe resulting array
}

// RollMode is used to identify how probabilities are calculated
type RollMode int

const (
	ODDS  RollMode = 1 + iota // Odds mode
	PROBS                     // Probabilities mode
	NOTSUPPLIED
)

func rollDice(probs []float64) int {
	diceRoll := rand.Float64()
	side := 1
	ack := 0.0
	for _, p := range probs {
		ack += p
		if diceRoll < ack {
			break
		} else {
			side++
		}
	}
	return side
}

// odds, probabilites
func calculateProbablities(data *string, mode RollMode) ([]float64, error) {
	var tmpString []string

	if mode == NOTSUPPLIED {
		return []float64{1 / 6.0, 1 / 6.0, 1 / 6.0, 1 / 6.0, 1 / 6.0, 1 / 6.0}, nil
	} else if mode != ODDS && mode != PROBS {
		return nil, errors.New("Supplied mode for calculate probablities is unknown")
	}
	tmpString = strings.Split(*data, ",")
	probs := make([]float64, len(tmpString))

	sum := 0.0
	for i, s := range tmpString {
		f, err := strconv.ParseFloat(s, 64)
		if err != nil {
			if strings.Index(s, "/") > -1 {
				fraction := strings.Split(s, "/")
				numerator, err := strconv.ParseFloat(fraction[0], 64)
				denominator, err2 := strconv.ParseFloat(fraction[1], 64)
				if err != nil || err2 != nil {
					return nil, errors.New("Fraction can't be converted to float: Fraction=" + s)
				}
				f = numerator / denominator
			} else {
				return nil, errors.New("Unable to parse number: " + s)
			}
		}
		probs[i] = f
		sum += f
	}

	if mode == ODDS {
		for i := 0; i < len(probs); i++ {
			probs[i] /= sum
		}
	}
	return probs, nil
}

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "diceroller API - call /1")
}

func srvRollHandler(w http.ResponseWriter, r *http.Request, mode RollMode) {
	var srvResult DiceRolls
	var probs []float64
	var probsError error

	vars := mux.Vars(r)

	nrOfRolls, err := strconv.Atoi(vars["rolls"])

	if err != nil {
		srvResult.Status = "The amount of rolls could not be interpreted as a string: " + html.EscapeString(vars["rolls"])
		json.NewEncoder(w).Encode(srvResult)
		return
	}
	data := vars["data"]
	if nrOfRolls <= 0 {
		probsError = errors.New("Nr of rolls was zero or lower")
	} else if mode == PROBS || mode == ODDS {
		probs, probsError = calculateProbablities(&data, mode)
	} else if mode == NOTSUPPLIED {
		probs, probsError = calculateProbablities(nil, mode)
	}
	if probsError != nil {
		srvResult.Status = probsError.Error()
	} else {

		if nrOfRolls > 1000000 {
			srvResult.Status = "No more than 1 000 000 dice rolls is allowed"
		} else {
			var rolls = make([]int, nrOfRolls)
			for x := 0; x < nrOfRolls; x++ {
				rolls[x] = rollDice(probs)
			}
			srvResult.Result = rolls
			srvResult.Status = "ok"
		}
	}
	json.NewEncoder(w).Encode(srvResult)
}

func handleRequests(port int) {
	strPort := ":" + strconv.Itoa(port)
	r := mux.NewRouter()
	r.HandleFunc("/", homePage)
	r.HandleFunc("/{rolls}", func(w http.ResponseWriter, r *http.Request) { srvRollHandler(w, r, NOTSUPPLIED) })
	r.HandleFunc("/{rolls}/probs/{data}", func(w http.ResponseWriter, r *http.Request) { srvRollHandler(w, r, PROBS) })
	r.HandleFunc("/{rolls}/odds/{data}", func(w http.ResponseWriter, r *http.Request) { srvRollHandler(w, r, ODDS) })
	fmt.Println("Starting server on port " + strPort)
	log.Fatal(http.ListenAndServe(strPort, r))
}

func main() {
	var mode RollMode

	rand.Seed(time.Now().UTC().UnixNano())

	// Command line flags
	probablities := flag.String("probs", "", "Probablities space separated in a string \"0.5 0.5\"")
	odds := flag.String("odds", "", "Odds space separated in a string  \"1 1 1 4 \"")
	nrOfRolls := flag.Int("rolls", 1, "Nr of dice rolls")
	srv := flag.Bool("srv", false, "Run as rest-server")
	port := flag.Int("port", 10000, "Port for rest-server")
	flag.Parse()

	if *srv {
		// Servermode
		handleRequests(*port)
	} else {
		// Commandline mode
		var data *string

		if *odds != "" {
			mode, data = ODDS, odds
		} else if *probablities != "" {
			mode, data = PROBS, probablities
		} else {
			mode, data = NOTSUPPLIED, nil
		}

		probs, err := calculateProbablities(data, mode)

		if err != nil {
			fmt.Println(err.Error())
			return

		} else if *nrOfRolls <= 0 {
			fmt.Println("Zero or negative amount of dice rolls")
			return
		}
		var rolls = make([]string, *nrOfRolls)

		for x := 0; x < *nrOfRolls; x++ {
			rolls[x] = strconv.Itoa(rollDice(probs))
		}
		fmt.Println(strings.Join(rolls, ", "))
	}
}
