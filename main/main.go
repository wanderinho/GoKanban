package main

import (
	"fmt"
	"kanban/handlers"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /boards", handlers.CreateBoard)
	mux.HandleFunc("POST /boards/{boardID}/columns", handlers.CreateColumn)
	mux.HandleFunc("GET /boards/{boardID}", handlers.GetBoard)
	mux.HandleFunc("DELETE /boards/{boardID}", handlers.DeleteBoard)
	mux.HandleFunc("PATCH /boards/{boardID}", handlers.UpdateBoardTitle)
	mux.HandleFunc("GET /boards/{boardID}/columns/{columnID}", handlers.GetColumn)
	mux.HandleFunc("DELETE /boards/{boardID}/columns/{columnID}", handlers.DeleteColumn)
	mux.HandleFunc("PATCH /boards/{boardID}/tasks/move", handlers.MoveTask)
	mux.HandleFunc("POST /boards/{boardID}/columns/{columnID}/tasks", handlers.CreateTask)
	mux.HandleFunc("GET /boards/{boardID}/columns/{columnID}/tasks/{taskID}", handlers.GetTask)
	mux.HandleFunc("DELETE /boards/{boardID}/columns/{columnID}/tasks/{taskID}", handlers.DeleteTask)
	mux.HandleFunc("PATCH /boards/{boardID}/columns/{columnID}", handlers.UpdateColumnTitle)
	mux.HandleFunc("PATCH /boards/{boardID}/columns/{columnID}/tasks/{taskID}", handlers.UpdateTaskDescription)

	err := http.ListenAndServe(":7890", mux)
	if err != nil {
		fmt.Println("err:", err.Error())
	}
}
