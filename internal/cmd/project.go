package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"
	"github.com/zivlakmilos/taskmaster/db"
)

var projectCmd = &cobra.Command{
	Use:   "project",
	Short: "Project operations",
	Run: func(cmd *cobra.Command, args []string) {
	},
}

var projectListCmd = &cobra.Command{
	Use:   "ls",
	Short: "List projects",
	Run: func(cmd *cobra.Command, args []string) {
		listProjects()
	},
}

var projectAddCmd = &cobra.Command{
	Use:   "add",
	Short: "Add project",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		addProject(args[0])
	},
}

var projectRmCmd = &cobra.Command{
	Use:   "rm",
	Short: "Remove project",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		removeProject(args[0])
	},
}

func init() {
	rootCmd.AddCommand(projectCmd)

	projectCmd.AddCommand(projectListCmd)
	projectCmd.AddCommand(projectAddCmd)
	projectCmd.AddCommand(projectRmCmd)
}

func listProjects() {
	con, err := openDb()
	if err != nil {
		exitWithError(err)
	}
	defer con.Close()

	store := db.NewProjectStore(con)

	projects, err := store.GetAll()
	if err != nil {
		exitWithError(err)
	}

	w := tabwriter.NewWriter(os.Stdout, 1, 1, 1, ' ', 0)
	fmt.Fprintf(w, "id\tname\tstatus\n")
	for _, project := range projects {
		fmt.Fprintf(w, "%s\t%s\t%s\n", project.Id, project.Name, project.Status)
	}
	w.Flush()
}

func addProject(name string) {
	con, err := openDb()
	if err != nil {
		exitWithError(err)
	}
	defer con.Close()

	store := db.NewProjectStore(con)

	project := db.NewProject()
	project.Name = name

	err = store.Save(project)
	if err != nil {
		exitWithError(err)
	}
}

func removeProject(it string) {
}
