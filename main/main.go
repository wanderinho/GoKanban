package main

import (
	"fmt"
	"kanban/handlers"
	"kanban/src"
	"log"
	"net/http"
	"os"
)

func main() {

	storage, err := src.NewPostgresStorage(os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
		return
	}
	defer storage.Close()
	
	handler := handlers.NewHandler(storage)
	
	mux := http.NewServeMux()

	mux.HandleFunc("POST /boards", handler.CreateBoard)
	mux.HandleFunc("POST /boards/{boardID}/columns", handler.CreateColumn)
	mux.HandleFunc("GET /boards/{boardID}", handler.GetBoard)
	mux.HandleFunc("DELETE /boards/{boardID}", handler.DeleteBoard)
	mux.HandleFunc("PATCH /boards/{boardID}", handler.UpdateBoardTitle)
	mux.HandleFunc("GET /boards/{boardID}/columns/{columnID}", handler.GetColumn)
	mux.HandleFunc("DELETE /boards/{boardID}/columns/{columnID}", handler.DeleteColumn)
	mux.HandleFunc("PATCH /boards/{boardID}/tasks/move", handler.MoveTask)
	mux.HandleFunc("POST /boards/{boardID}/columns/{columnID}/tasks", handler.CreateTask)
	mux.HandleFunc("GET /boards/{boardID}/columns/{columnID}/tasks/{taskID}", handler.GetTask)
	mux.HandleFunc("DELETE /boards/{boardID}/columns/{columnID}/tasks/{taskID}", handler.DeleteTask)
	mux.HandleFunc("PATCH /boards/{boardID}/columns/{columnID}", handler.UpdateColumnTitle)
	mux.HandleFunc("PUT /boards/{boardID}/columns/{columnID}/tasks/{taskID}", handler.UpdateTaskTitle)
	mux.HandleFunc("PATCH /boards/{boardID}/columns/{columnID}/tasks/{taskID}", handler.UpdateTaskDescription)

	err = http.ListenAndServe(":7890", mux)
	if err != nil {
		fmt.Println("err:", err.Error())
	}
}
