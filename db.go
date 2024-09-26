package main

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func ConnectDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open(config.DBPath), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to database!")
	}

	return db
}

func addTodo(user, content string) error {
	newTodo := todo{User: user, Task: content, Completed: false}

	result := db.Create(&newTodo)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func getTodos(user string) ([]todo, error) {
	var todo []todo
	result := db.Where("user = ? AND completed = ?", user, false).Find(&todo)
	if result.Error != nil {
		return nil, result.Error
	}

	return todo, nil
}

func getNumCompleted(user string) (int64, error) {
	var count int64
	result := db.Where("user = ? AND completed = ?", user, true).Table("todos").Count(&count)
	if result.Error != nil {
		return 0, result.Error
	}

	return count, nil
}

func getAllTodos() ([]todo, error) {
	var todo []todo

	result := db.Where("completed = ?", false).Find(&todo)
	if result.Error != nil {
		return nil, result.Error
	}
	return todo, nil
}

func completeTodoById(id int) (todo, error) {
	var todoRes todo
	result := db.Where("Id = ?", id).First(&todoRes)
	if result.Error != nil {
		return todo{}, result.Error
	}

	result = db.Model(&todo{}).Where("Id = ?", id).Update("completed", true)
	if result.Error != nil {
		return todo{}, result.Error
	}

	return todoRes, nil
}

func deleteTodos(user string) error {
	return db.Where("user = ?", user).Delete(&todo{}).Error
}

func deleteTodoById(id int) error {
	return db.Where("Id = ?", id).Delete(&todo{}).Error
}

func updateTodoById(id int, content string) error {
	return db.Model(&todo{}).Where("Id = ?", id).Update("task", content).Error
}
