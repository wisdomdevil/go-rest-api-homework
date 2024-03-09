package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// Task ...
type Task struct {
	ID           string   `json:"id"`
	Description  string   `json:"description"`
	Note         string   `json:"note"`
	Applications []string `json:"applications"`
}

var tasks = map[string]Task{
	"1": {
		ID:          "1",
		Description: "Сделать финальное задание темы REST API",
		Note:        "Если сегодня сделаю, то завтра будет свободный день. Ура!",
		Applications: []string{
			"VS Code",
			"Terminal",
			"git",
		},
	},
	"2": {
		ID:          "2",
		Description: "Протестировать финальное задание с помощью Postmen",
		Note:        "Лучше это делать в процессе разработки, каждый раз, когда запускаешь сервер и проверяешь хендлер",
		Applications: []string{
			"VS Code",
			"Terminal",
			"git",
			"Postman",
		},
	},
}

// Ниже напишите обработчики для каждого эндпоинта
// ...
func postTasksHandle(res http.ResponseWriter, req *http.Request) {
	var newTask Task
	var buf bytes.Buffer

	_, err := buf.ReadFrom(req.Body)
	//Проверка на не пустую ошибку, при чтении тела запроса, при ошибке возвращаем 400 Bad Request
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}
	//Проверка на не пустую ошибку, при дессериализации JSON входящего в теле запроса, при наличии ошибки возвращаем 400 Bad Request
	if err = json.Unmarshal(buf.Bytes(), &newTask); err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}
	// Добавляем в мапу новую таску
	tasks[newTask.ID] = newTask
	//Задаем Header
	res.Header().Set("Content-Type", "application/json")
	//Задаем StatusCode
	res.WriteHeader(http.StatusCreated)
}
func getTasksHandle(res http.ResponseWriter, req *http.Request) {

	jTasks, err := json.Marshal(tasks)
	//Проверка на ошибку во время сериализации json, если ошибка не пуста возвращаем 400 Bad Request и сообщение ошибки
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
	}
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
	res.Write(jTasks)
}
func getTaskByIdHandle(res http.ResponseWriter, req *http.Request) {
	id := chi.URLParam(req, "id")
	task, ok := tasks[id]
	//Если в мапе нет нужного id возвращаем 400 Bad Request и сообщение
	if !ok {
		http.Error(res, "Такой задачи не существет", http.StatusBadRequest)
		return
	}
	//Сереализуем нужную нам таску по id в json
	resp, err := json.Marshal(task)
	//Проверка на ошибку при сериализации, если ошибка есть - возвращаем 400 Bad Request
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
	}
	//Задаем заголовок с типом контента
	res.Header().Set("Content-Type", "application/json")
	//Возвращаем 200 статус ответа
	res.WriteHeader(http.StatusOK)
	//Возвращаем сериализованные данные клиенту
	res.Write(resp)
}
func deleteTaskByIdHandle(res http.ResponseWriter, req *http.Request) {
	id := chi.URLParam(req, "id")

	_, ok := tasks[id]
	//Проверка на ошибку в мапе(например отсуствие элемента), можно было бы вернуть информативное с
	if !ok {
		http.Error(res, "Такой задачи не существует", http.StatusBadRequest)
		return
	}
	//Удаляем таску из мапы
	delete(tasks, id)
	//Выставляем Header
	res.Header().Set("Content-Type", "application/json")
	//Возвращаем 200 код
	res.WriteHeader(http.StatusOK)
}
func main() {
	r := chi.NewRouter()

	r.Get("/tasks", getTasksHandle)
	r.Post("/tasks", postTasksHandle)
	r.Get("/tasks/{id}", getTaskByIdHandle)
	r.Delete("/tasks/{id}", deleteTaskByIdHandle)
	if err := http.ListenAndServe(":8080", r); err != nil {
		fmt.Printf("Ошибка при запуске сервера: %s", err.Error())
		return
	}
}
