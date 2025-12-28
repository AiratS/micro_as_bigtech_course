package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/fatih/color"
)

const (
	baseUrl       = "http://localhost:8081"
	createPostfix = "/notes"
	getPostfix    = "/notes/%d"
)

type NoteInfo struct {
	Title    string `json:"title"`
	Context  string `json:"context"`
	Author   string `json:"author"`
	IsPublic bool   `json:"is_public"`
}

type Note struct {
	ID        int64     `json:"id"`
	Info      NoteInfo  `json:"info"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func createNote() (Note, error) {
	note := &NoteInfo{
		Title:    gofakeit.BeerName(),
		Context:  gofakeit.IPv4Address(),
		Author:   gofakeit.Name(),
		IsPublic: gofakeit.Bool(),
	}

	data, err := json.Marshal(note)
	if err != nil {
		return Note{}, err
	}

	resp, err := http.Post(baseUrl+createPostfix, "application/json", bytes.NewBuffer(data))
	if err != nil {
		return Note{}, err
	}
	defer func() {
		err = resp.Body.Close()
		if err != nil {
			log.Print(err)
		}
	}()

	if resp.StatusCode != http.StatusCreated {
		return Note{}, err
	}

	createdNote := &Note{}
	if err = json.NewDecoder(resp.Body).Decode(createdNote); err != nil {
		return Note{}, err
	}

	return *createdNote, nil
}

func getNote(ID int64) (Note, error) {
	url := fmt.Sprintf(baseUrl+getPostfix, ID)
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return Note{}, err
	}

	if resp.StatusCode != http.StatusOK {
		return Note{}, errors.New("Is not Ok")
	}

	note := &Note{}
	err = json.NewDecoder(resp.Body).Decode(note)
	if err != nil {
		return Note{}, err
	}

	return *note, nil
}

func main() {
	note, err := createNote()
	if err != nil {
		log.Fatal("Failed to create note", err)
	}

	log.Printf(color.RedString("Note created\n"), color.GreenString("Note: %+v\n", note))

	note, err = getNote(note.ID)
	if err != nil {
		log.Fatal("Failed to create note", err)
	}

	log.Printf(color.RedString("Note get info\n"), color.GreenString("Note: %+v\n", note))
}
