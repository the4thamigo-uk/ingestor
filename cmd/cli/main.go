package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/the4thamigo-uk/ingestor/pkg/client"
	"os"
)

var addr string
var file string
var id string

func main() {
	var root = &cobra.Command{}
	root.PersistentFlags().StringVarP(&addr, "server", "s", ":8080", "Address of the ingestor service.")
	root.MarkFlagRequired("server")

	var add = &cobra.Command{
		Use: "add a source data file to the server.",
		Run: addSource,
	}

	add.Flags().StringVarP(&file, "file", "f", "", "Filename of source data file for server to publish.")
	add.MarkFlagRequired("file")

	var list = &cobra.Command{
		Use: "list source ids on server.",
		Run: listSources,
	}

	var read = &cobra.Command{
		Use: "read a source data file from the server.",
		Run: readSource,
	}

	read.Flags().StringVarP(&id, "id", "i", "", "id of source data file")
	read.MarkFlagRequired("id")

	root.AddCommand(add, list, read)
	root.Execute()
}

func addSource(cmd *cobra.Command, args []string) {
	c, err := client.New(addr)
	if err != nil {
		handleError(err)
		return
	}
	defer c.Close()

	id, err := c.AddSource(file)
	if err != nil {
		handleError(err)
		return
	}
	fmt.Println(id)
}

func listSources(cmd *cobra.Command, args []string) {
	c, err := client.New(addr)
	if err != nil {
		handleError(err)
		return
	}
	defer c.Close()

	srcs, err := c.ListSources()
	if err != nil {
		handleError(err)
		return
	}
	for _, src := range srcs {
		fmt.Println(src)
	}
}

func readSource(cmd *cobra.Command, args []string) {
	c, err := client.New(addr)
	if err != nil {
		handleError(err)
		return
	}
	defer c.Close()

	err = c.ReadSource(id,
		func(c client.Contact) error {
			fmt.Printf("%s, %s, %s\n", c.Mobile(), c.Name(), c.Email())
			return nil
		})
	if err != nil {
		handleError(err)
		return
	}

}

func handleError(err error) {
	fmt.Println(err)
	pflag.PrintDefaults()
	os.Exit(1)
}
