package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/google/uuid"
	"github.com/spf13/cobra"
	"github.com/zivlakmilos/taskmaster/db"
)

type TaskConfig struct {
	projectId string
}

var taskCfg TaskConfig

var taskCmd = &cobra.Command{
	Use:   "task",
	Short: "Task operations",
	Run: func(cmd *cobra.Command, args []string) {
	},
}

var taskListCmd = &cobra.Command{
	Use:   "ls",
	Short: "List tasks",
	Run: func(cmd *cobra.Command, args []string) {
		listTasks()
	},
}

var taskAddCmd = &cobra.Command{
	Use:   "add",
	Short: "Add task",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		addTask(args[0])
	},
}

var taskRmCmd = &cobra.Command{
	Use:   "rm",
	Short: "Remove task",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		removeTask(args[0])
	},
}

func init() {
	rootCmd.AddCommand(taskCmd)

	rootCmd.PersistentFlags().StringVarP(&taskCfg.projectId, "project", "p", "", "Project id/name")

	taskCmd.AddCommand(taskListCmd)
	taskCmd.AddCommand(taskAddCmd)
	taskCmd.AddCommand(taskRmCmd)
}

func listTasks() {
	con, err := openDb()
	if err != nil {
		exitWithError(err)
	}

	store := db.NewTaskStore(con)

	var tasks []*db.Task

	if taskCfg.projectId != "" {
		_, err = uuid.Parse(taskCfg.projectId)
		if err != nil {
			tasks, err = store.GetAllByProjectName(taskCfg.projectId)
			if err != nil {
				exitWithError(err)
			}
		} else {
			tasks, err = store.GetAllByProjectId(taskCfg.projectId)
			if err != nil {
				exitWithError(err)
			}
		}
	} else {
		tasks, err = store.GetAll()
		if err != nil {
			exitWithError(err)
		}
	}

	w := tabwriter.NewWriter(os.Stdout, 1, 1, 1, ' ', 0)
	fmt.Fprintf(w, "id\tdescription\tpriority\tprogress\tstatus\n")
	for _, task := range tasks {
		fmt.Fprintf(w, "%s\t%s\t%d\t%d%%\t%s\n", task.Id, task.Description, task.Priority, task.Progress, task.Status)
	}
	w.Flush()
}

func addTask(description string) {
	if taskCfg.projectId == "" {
		exitWithError(fmt.Errorf("project id is required"))
	}

	con, err := openDb()
	if err != nil {
		exitWithError(err)
	}

	store := db.NewTaskStore(con)

	task := db.NewTask()
	task.DateGet = db.Now()
	task.Description = description
	task.ProjectId = taskCfg.projectId

	err = store.Save(task)
	if err != nil {
		exitWithError(err)
	}
}

func removeTask(id string) {
	con, err := openDb()
	if err != nil {
		exitWithError(err)
	}

	store := db.NewTaskStore(con)
	err = store.Delete(id)
	if err != nil {
		exitWithError(err)
	}
}
