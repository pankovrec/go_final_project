package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

const (
	timeformat = "20060102"
)

type TaskService struct {
	storage *Storage
}

func NewTaskService(storage *Storage) TaskService {
	return TaskService{storage: storage}
}

func NextDate(date, now time.Time, repeat string) (time.Time, error) {

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

func (t TaskService) getTaskDone(id int) error {
	task, err := t.storage.SelectById(id)
	if err != nil {
		return fmt.Errorf(`{"error":"%s"}`, err.Error())
	}

	if task.Repeat == "" {
		return t.storage.DeleteTask(id)
	}

	date, err := time.Parse(timeformat, task.Date)
	if err != nil {
		return fmt.Errorf("invalid date format in db")
	}

	nextDate, _ := NextDate(date, time.Now(), task.Repeat)

	task.Date = nextDate.Format(timeformat)

	err = t.storage.UpdateTask(task)
	if err != nil {
		return err
	}

	return nil
}

func checkNextDate(s string) error {
	if s == "" {
		return errors.New("empty repeat")
	}
	return nil
}

func CheckTask(task Task) (Task, error) {
	if task.Title == "" {
		return Task{}, fmt.Errorf("empty title")
	}

	y, m, d := time.Now().Date()
	today := time.Date(y, m, d, 0, 0, 0, 0, time.UTC)
	date := today

	var err error

	if task.Date != "" {
		date, err = time.Parse(timeformat, task.Date)
		if err != nil {
			return Task{}, fmt.Errorf("invalid date format")
		}
	}

	if date.Before(today) {
		if task.Repeat == "" {
			date = today
		} else {
			date, err = NextDate(date, today, task.Repeat)
			if err != nil {
				return Task{}, fmt.Errorf("can't get next date: %w", err)
			}
		}
	}

	id := "0"
	if task.ID != "" {
		id = task.ID
	}

	return Task{
		ID:      id,
		Date:    date.Format(timeformat),
		Title:   task.Title,
		Comment: task.Comment,
		Repeat:  task.Repeat,
	}, nil
}

func MakeTasksList(tasks []Task) []Task {
	list := make([]Task, 0, len(tasks))

	for _, task := range tasks {
		list = append(list, Task{
			ID:      task.ID,
			Date:    task.Date,
			Title:   task.Title,
			Comment: task.Comment,
			Repeat:  task.Repeat,
		})
	}
	return list
}
