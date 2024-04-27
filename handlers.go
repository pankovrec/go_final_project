package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
)

func (t TaskService) ExtractParamsDateHandler(w http.ResponseWriter, r *http.Request) {
	now := r.URL.Query().Get("now")
	date := r.URL.Query().Get("date")
	repeat := r.URL.Query().Get("repeat")

	nowTime, err := time.Parse(timeformat, now)
	responseWithError(w, err, "parsing string now")

	dateTime, err := time.Parse(timeformat, date)
	responseWithError(w, err, "parsing string date")

	err = checkNextDate(repeat)
	responseWithError(w, err, "repeat not empty")

	nextDate, err := NextDate(dateTime, nowTime, repeat)
	responseWithError(w, err, "error calculating")

	log.Println("[Info] FOR now =", now, "date =", date, "repeat =", repeat)

	w.Write([]byte(nextDate.Format(timeformat)))
}

func responseWithError(w http.ResponseWriter, err error, message string) {
	if err != nil {
		log.Printf("%s: %s", err, message)
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	log.Printf("%s: %s", message)

}

func (t TaskService) TaskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		t.addTask(w, r)
	case http.MethodGet:
		t.getTask(w, r)
	case http.MethodPut:
		t.editTask(w, r)
	case http.MethodDelete:
		t.removeTask(w, r)
	}
}

func (t TaskService) addTask(w http.ResponseWriter, r *http.Request) {
	var input Task

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err.Error()), http.StatusBadRequest)
		log.Println("[WARN] Failed json decoding:", err)
		return
	}

	task, err := CheckTask(input)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err.Error()), http.StatusBadRequest)
		log.Println("[WARN] Failed input check:", err)
		return
	}

	log.Println("[TASK] : " + task.Date + task.Title + task.Comment + task.Repeat)

	id, err := t.storage.InsertTask(task)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err.Error()), http.StatusInternalServerError)
		log.Println("[WARN] Failed to add a task:", err)
		return
	}

	log.Println("[Info] Success: Task added with id = " + strconv.Itoa(id))
	w.Write([]byte(fmt.Sprintf(`{"id":"%d"}`, id)))
}

func (t TaskService) TasksHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	tasksFromDB, _ := t.storage.SelectTasks()
	tasks := MakeTasksList(tasksFromDB)
	response := ResponseTasks{Tasks: tasks}
	responseBody, err := json.Marshal(response)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err.Error()), http.StatusInternalServerError)
		log.Println("json Marshal:", err)

		return
	}
	log.Println("[Info] Success: tasks from DB are given")
	w.Write(responseBody)
}

func (t TaskService) getTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	num := r.FormValue("id")
	if num == "" {
		http.Error(w, `{"error":"true"}`, http.StatusBadRequest)
		log.Println("[WARN] empty param")
		return
	}

	id, err := strconv.Atoi(num)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err.Error()), http.StatusBadRequest)
		log.Println("[WARN] Failed convertation:", err)
		return
	}

	task, err := t.storage.SelectById(id)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err.Error()), http.StatusBadRequest)
		log.Println("[WARN] Failed convertation:", err)
		return
	}

	responseBody, err := json.Marshal(task)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err.Error()), http.StatusInternalServerError)
		log.Println("json Marshal:", err)

		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err = w.Write(responseBody); err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err.Error()), http.StatusInternalServerError)
	}
}

func (t TaskService) editTask(w http.ResponseWriter, r *http.Request) {

	var input Task
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err.Error()), http.StatusBadRequest)
		log.Println("json Decoder:", err)

		return
	}

	task, err := CheckTask(input)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err.Error()), http.StatusBadRequest)
		log.Println("CheckInput:", err)

		return
	}

	if task.ID == "0" {
		http.Error(w, `{"error":"wrong id"}`, http.StatusBadRequest)

		return
	}

	err = t.storage.UpdateTask(task)

	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err.Error()), http.StatusInternalServerError)
		log.Println("UpdateTask:", err)

		return
	}

	w.Write([]byte(`{}`))
}
func (t TaskService) removeTask(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		http.Error(w, `{"error":"wrong id"}`, http.StatusBadRequest)
		log.Println("Atoi:", err)

		return
	}

	err = t.storage.DeleteTask(id)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err.Error()), http.StatusInternalServerError)

		return
	}

	w.Write([]byte(`{}`))
}

func (t TaskService) DoneHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	num := r.FormValue("id")
	id, err := strconv.Atoi(num)
	if err != nil {
		http.Error(w, `{"error":"wrong id"}`, http.StatusBadRequest)
		log.Println("Atoi:", err)

		return
	}

	err = t.getTaskDone(id)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err.Error()), http.StatusInternalServerError)

		return
	}

	w.Write([]byte(`{}`))
}
