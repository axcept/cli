package main

import (
  "fmt"
  "os"
  "time"
  "github.com/spf13/cobra"
  "github.com/dghubble/sling"
)

type StartTestrunResponse struct {
	Success   bool   `json:"success,omitempty"`
	TestRunId int    `json:"testRunId,omitempty"`
}

type TestRunResponse struct {
	Success   bool   `json:"success,omitempty"`
	Running   bool   `json:"running,omitempty"`
}

type TestRunRequest struct {
    Token     string   `json:"token,omitempty"`
}

func getEnv(key, defaultValue string) string {
    value := os.Getenv(key)
    if len(value) == 0 {
        return defaultValue
	}
	value = value + "/api"
	fmt.Println("Using Axcept Endpoint: " + value)
    return value
}

var url = getEnv("AXCEPT_SERVICE_URL", "https://app.axcept.io/api") 

func main() {
  var environment string
  var token string

  var testRun = &cobra.Command{
    Use:   "testrun",
    Short: "Start a test run for a specific environment",
    Long: `Long.`,
    Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Starting test run for environment Id: " + environment)

		var endpoint = url + "/environment/" + environment + "/test-run"

		body := &TestRunRequest{
			Token: token,
		}

		responseData := new(StartTestrunResponse)

		_, err := sling.New().Post(endpoint).BodyJSON(body).ReceiveSuccess(responseData)

		if &err != nil && responseData.Success {
			fmt.Println("Test run started successfully with id ", responseData.TestRunId)
		} else {
			fmt.Println("Can't start test run. Please check your environment id and token")
			os.Exit(1)
		}
		start := time.Now()

		for i := 1; i <= 10; i++ {
			testRunResponse := new(TestRunResponse)
			_, err := sling.New().Get(endpoint+fmt.Sprintf("/%d", responseData.TestRunId)).ReceiveSuccess(testRunResponse)

			if &err != nil {
				if testRunResponse.Running == false {
					fmt.Println("Test run finished. Success:", testRunResponse.Success)
					if testRunResponse.Success {
						os.Exit(0)
					} else {
						os.Exit(1)
					}
					fmt.Printf("Took %s", time.Since(start))
				}
			}
			time.Sleep(2 * time.Second)		
		}

    },
  }

  testRun.Flags().StringVarP(&environment, "environment", "e","", "Environment Id that should be tested.")
  testRun.MarkFlagRequired("environment")

  testRun.Flags().StringVarP(&token, "token", "t","", "Axcept project token. Found in project settings")
  testRun.MarkFlagRequired("token")

  var rootCmd = &cobra.Command{Use: "axcept"}
  rootCmd.AddCommand(testRun)
  rootCmd.Execute()
}