package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

type Input struct {
	Inputs     string `json:"inputs"`
	Parameters struct {
		MaxNewTokens   int  `json:"max_new_tokens"`
		ReturnFullText bool `json:"return_full_text"`
	} `json:"parameters"`
	Options struct {
		WaitForModel bool `json:"wait_for_model"`
	} `json:"options"`
}

type ServerOutput struct {
	GeneratedText string `json:"generated_text"`
}

func toInput(ctx string) Input {
	result := Input{
		Inputs: ctx,
		Parameters: struct {
			MaxNewTokens   int  `json:"max_new_tokens"`
			ReturnFullText bool `json:"return_full_text"`
		}{
			MaxNewTokens:   100,
			ReturnFullText: false,
		},
		Options: struct {
			WaitForModel bool `json:"wait_for_model"`
		}{
			WaitForModel: true,
		},
	}
	return result
}

func request(input Input) []ServerOutput {

	inputJSON, err := json.Marshal(input)
	if err != nil {
		fmt.Println("Error marshalling input data:", err)
		return nil
	}

	resp, err := http.Post("https://api-inference.huggingface.co/models/OpenAssistant/oasst-sft-1-pythia-12b", "application/json", bytes.NewBuffer(inputJSON))
	if err != nil {
		fmt.Println("Error sending request to server:", err)
		return nil
	}

	defer resp.Body.Close()

	respData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading server response:", err)
		return nil
	}
	var serverResp []ServerOutput
	err = json.Unmarshal(respData, &serverResp)
	if err != nil {
		fmt.Println("Error unmarshalling server response:", err)
		fmt.Println("Server response: ", string(respData))
		return nil
	}
	return serverResp
}

func interact(ctx string) string {
	inputData := toInput(ctx)
	ctx = ctx + "<|assistant|>"
	for _, output := range request(inputData) {
		serverOutput := output.GeneratedText
		// TODO: not print, but return this as array of strings along with the
		// ctx
		fmt.Println(serverOutput)

		ctx = ctx + serverOutput
	}
	return ctx + "<|endoftext|>"
}

func main() {
	headless := flag.Bool("headless", false, "Whether to run the tool in headless mode or not")
	oneshot := flag.Bool("oneshot", false, "Wheter read the full input before sending it only once")
	flag.Parse()

	ctx := "<|prompter|>"

	scanner := bufio.NewScanner(os.Stdin)
	if *oneshot {
		for scanner.Scan() {
			userInput := scanner.Text()

			ctx = ctx + userInput
		}

		ctx = ctx + "<|endoftext|>"

		ctx = interact(ctx)
	} else {

		if !*headless {
			fmt.Print("User: ")
		}

		for scanner.Scan() {
			userInput := scanner.Text()

			ctx = ctx + userInput

			if userInput == "" {

				ctx = ctx + "<|endoftext|>"

				if !*headless {
					fmt.Print("Helper: ")
				}

				ctx = interact(ctx) + "<|prompter|>"

				if !*headless {
					fmt.Println()
					fmt.Print("User: ")
				}

			}
		}
	}

}
