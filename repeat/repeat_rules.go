package repeat

import (

	"net/http"
	"log"
	"time"
	"errors"
	"strconv"
	"strings"
)

const(
	timeformat = "20060102"
)	

func ExtractParamsDate(w http.ResponseWriter, r *http.Request) {
	now := r.URL.Query().Get("now")
	date := r.URL.Query().Get("date")
	repeat := r.URL.Query().Get("repeat")

	nowTime, err := time.Parse(timeformat, now)
	failOnError(w, err, "parsing string now")

	dateTime, err := time.Parse(timeformat, date)
	failOnError(w, err, "parsing string date")

	err = checkError(repeat)
	failOnError(w, err, "repeat not empty")

	nextDate, err := NextDate(dateTime, nowTime, repeat)
	failOnError(w, err, "error calculating")

	log.Println("[Info] FOR now =", now, "date =", date, "repeat =", repeat)

	w.Write([]byte(nextDate.Format(timeformat)))
}

func failOnError(w http.ResponseWriter, err error, message string) {
	if err != nil {
		log.Printf("%s: %s", err, message)
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	log.Printf("%s: %s", message)
}

func NextDate(date, now time.Time, repeat string) (time.Time, error){

	nextDate := date
		repeat_array := strings.Split(repeat, " ")
		switch repeat_array[0] {
		case "y":
			for {
				nextDate = nextDate.AddDate(1, 0, 0)
				if nextDate.After(now) {
					break
				}
			}
		case "d":
			if len(repeat_array) != 2 {
				return nextDate, errors.New("invalid repeat rule format")
			}
			days, err := strconv.Atoi(repeat_array[1])
			if err != nil {
				return nextDate, err
			}
			if days > 400 {
				return nextDate, errors.New("max days in repeat rule must be 400")
			}
			for {
				nextDate = nextDate.AddDate(0, 0, days)
				if nextDate.After(now) {
					break
				}
			}
		default:
			return nextDate, errors.New("invalid repeat format")
		}
		return nextDate, nil
	}
	
	func checkError(s string) error {
		if s == "" {
			return errors.New("empty repeat")
		}
		return nil
	}
	
