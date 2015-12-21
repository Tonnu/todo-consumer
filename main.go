package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/codegangsta/cli"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
)

func main() {
	app := cli.NewApp()

	app.Name = "todo"
	app.Flags = []cli.Flag{
		cli.IntFlag{
			Name:  "port",
			Value: 5000,
		},
	}

	app.Commands = []cli.Command{
		{
			Name:  "add",
			Usage: "add a taks to the list",
			Action: func(c *cli.Context) {
				fmt.Println("added task: ", c.Args().First())
				err := addTodo(c.Args().First())
				if err != nil {
					fmt.Println("An internal server error occurred")
					os.Exit(1)
				}
			},
		},
		{
			Name:  "get",
			Usage: "get task with id",
			Action: func(c *cli.Context) {
				fmt.Println("tasks: ")
				id, _ := strconv.Atoi(c.Args().First())
				todo, err := getTodo(id)
				if err != nil {
					fmt.Println("An internal server error occurred")
					os.Exit(1)
				}

				fmt.Println(fmt.Sprintf("ID: %s\nTodo: %s\nStatus: %s", todo.ID, todo.Title, todo.Status))
			},
		},
	}

	app.Run(os.Args)
}

type Todo struct {
	ID     string `json:"id,omitempty"`
	Title  string `json:"title,omitempty"`
	Status string `json:"status,omitempty"`
}

func getTodo(id int) (*Todo, error) {
	resp, err := http.Get(fmt.Sprintf("http://golang-api:5000/todo/%d", id))
	if err != nil {
		fmt.Println("got err: ", err)
		return nil, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	var todo Todo
	err = json.Unmarshal(body, &todo)

	if err != nil {
		return nil, err
	}

	return &todo, nil
}

func addTodo(title string) error {
	todo := &Todo{
		Title:  title,
		Status: "started",
	}

	data, _ := json.Marshal(todo)

	resp, err := http.Post(fmt.Sprintf("http://golang-api:5000/todo"), "application/json", bytes.NewReader(data))

	if err != nil {
		return err
	}

	defer resp.Body.Close()
	id, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(fmt.Sprintf("Task ID: %s", string(id)))
	return nil
}
