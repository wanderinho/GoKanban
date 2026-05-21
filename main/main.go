package main

import (
	"fmt"
	"kanban/handlers"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/boards", handlers.CreateBoard)
	mux.HandleFunc("/boards/{boardID}/columns", handlers.CreateColumn)
	mux.HandleFunc("GET /boards/{boardID}", handlers.GetBoard)
	mux.HandleFunc("DELETE /boards/{boardID}", handlers.DeleteBoard)
	mux.HandleFunc("PATCH /boards/{boardID}", handlers.UpdateBoardTitle)
	mux.HandleFunc("GET /boards/{boardID}/columns/{columnID}", handlers.GetColumn)
	mux.HandleFunc("DELETE /boards/{boardID}/columns/{columnID}", handlers.DeleteColumn)
	mux.HandleFunc("/boards/{boardID}/tasks/move", handlers.MoveTask)

	err := http.ListenAndServe(":7890", mux)
	if err != nil {
		fmt.Println("err:", err.Error())
	}
}
