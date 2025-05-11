// main.go
package main

import (
	"encoding/json"      // JSON 인코딩/디코딩
	"fmt"                // 콘솔 출력용
	"log"                // 에러 로그 출력
	"net/http"           // HTTP 서버 구성
	"strconv"            // 문자열 → 숫자 변환

	"go-rest-api/models" // Todo 데이터 구조 (구현한 모델 import)
	"github.com/gorilla/mux" // 라우터 패키지 (경로 핸들링에 사용)
)

var todos []models.Todo // Todo 리스트를 저장할 슬라이스

func getTodos(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(todos)
}

func createTodo(w http.ResponseWriter, r *http.Request) {
	var todo models.Todo
	_ = json.NewDecoder(r.Body).Decode(&todo)
	todo.ID = len(todos) + 1
	todos = append(todos, todo)
	json.NewEncoder(w).Encode(todo)
}

func updateTodo(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, _ := strconv.Atoi(params["id"])

	for index, item := range todos {
		if item.ID == id {
			var updated models.Todo
			_ = json.NewDecoder(r.Body).Decode(&updated)
			updated.ID = id
			todos[index] = updated
			json.NewEncoder(w).Encode(updated)
			return
		}
	}
	w.WriteHeader(http.StatusNotFound)
}

func deleteTodo(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, _ := strconv.Atoi(params["id"])

	for index, item := range todos {
		if item.ID == id {
			todos = append(todos[:index], todos[index+1:]...)
			w.WriteHeader(http.StatusNoContent)
			return
		}
	}
	w.WriteHeader(http.StatusNotFound)
}

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/todos", getTodos).Methods("GET")
	r.HandleFunc("/todos", createTodo).Methods("POST")
	r.HandleFunc("/todos/{id}", updateTodo).Methods("PUT")
	r.HandleFunc("/todos/{id}", deleteTodo).Methods("DELETE")

	fmt.Println("🚀 Server started at http://localhost:8000")
	log.Fatal(http.ListenAndServe(":8000", r))
}