package main

import (
	"crypto/rand"
	"encoding/json"
	"log"
	"math/big"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

const (
	baseUrl       = "localhost:8081"
	createPostfix = "/notes"
	getPostfix    = "/notes/{id}"
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

type SyncMap struct {
	elems map[int64]*Note
	m     sync.RWMutex
}

var notes = &SyncMap{
	elems: make(map[int64]*Note),
}

func createNoteHandler(w http.ResponseWriter, r *http.Request) {
	info := &NoteInfo{}
	err := json.NewDecoder(r.Body).Decode(info)
	if err != nil {
		http.Error(w, "Failed to parse body", http.StatusBadRequest)
		return
	}

	now := time.Now()
	note := &Note{
		ID:        randInt64(),
		Info:      *info,
		CreatedAt: now,
		UpdatedAt: now,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(note)
	if err != nil {
		http.Error(w, "Could not encode note", http.StatusInternalServerError)
		return
	}

	notes.m.Lock()
	defer notes.m.Unlock()

	notes.elems[note.ID] = note
}

func randInt64() int64 {
	a, _ := rand.Int(rand.Reader, big.NewInt(100))

	return a.Int64()
}

func getNoteHandler(w http.ResponseWriter, r *http.Request) {
	noteIDStr := chi.URLParam(r, "id")
	noteID, err := parseNoteID(noteIDStr)
	if err != nil {
		http.Error(w, "id param is not specified", http.StatusBadRequest)
		return
	}

	notes.m.RLock()
	defer notes.m.RUnlock()
	note, ok := notes.elems[noteID]
	if !ok {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(note)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

func parseNoteID(idStr string) (int64, error) {
	val, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return 0, err
	}

	return val, nil
}

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Post(createPostfix, createNoteHandler)
	r.Get(getPostfix, getNoteHandler)

	err := http.ListenAndServe(baseUrl, r)
	if err != nil {
		log.Fatal(err)
	}
}
