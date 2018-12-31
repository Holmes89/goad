package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"io/ioutil"
	"os"
	"smail/internal/commands"
	"strings"
)

var rootCmd = &cobra.Command{
	Use:   "smail",
	Short: "Simple Mail",
	Long: `A Simple Mail client built as an example on how to use GRPC, Cobra, and Mongo`,
}

func main() {

	//start()

	h, _ := homedir.Dir()
	viper.SetConfigType("json")
	viper.SetConfigName("config")
	viper.AddConfigPath(h+"/.smail/")
	if err := viper.ReadInConfig(); err != nil {
		conf := createConfig()
		viper.ReadConfig(bytes.NewBuffer(conf))
	}

	sendCmd := commands.NewSendCommand()
	inboxCmd := commands.NewInboxCommand()
	readCmd := commands.NewReadCommand()
	markCmd := commands.NewMarkCommand()
	deleteCmd := commands.NewDeleteCommand()
	rootCmd.AddCommand(sendCmd, inboxCmd, readCmd, markCmd, deleteCmd)
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func createConfig() []byte{
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Please Enter Name: ")
	text, _ := reader.ReadString('\n')
	text = strings.TrimSuffix(text, "\n")

	config := &Config{text}
	configJson, err := json.Marshal(config)
	if err != nil {
		panic(err)
	}
	h, _ := homedir.Dir()
	if _, err := os.Stat(h+"/.smail"); os.IsNotExist(err) {
		os.Mkdir(h+"/.smail", 0744)
	}

	if err := ioutil.WriteFile(h+"/.smail/config.json", configJson, 0744); err !=nil {
		panic(err)
	}

	return configJson
}

type Config struct {
	Username string `json:"username"`
}