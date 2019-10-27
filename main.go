package main

import "fmt"

var baseUrl = "https://app.asana.com/api/1.0"
var token = "0/140204d60a962b96ee4cfcb4ec9dff87"

func main() {
	client := NewAsanaSDK(token, "1114955233304260")
	// tasks := client.ProjectTasks(1143429561288049, 10)

	// for _, task := range *tasks {
	// 	tsk := client.Task(task.Id)
	// 	fmt.Println(tsk)
	// }

	task := client.Task(1144629399438001)

	stories := client.TaskStories(task.Id, 20)

	fmt.Println(stories)

	// client.ChangeTaskName(task.Id, "123"+task.Name)

}

// client id 1146727810660992

// client secret 6ab659430e63172a5ef7096ef632a0f9

// Sprint review token 0/140204d60a962b96ee4cfcb4ec9dff87
