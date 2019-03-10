package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/bitrise-io/go-utils/log"
	"github.com/kitasuke/go-swift/swift"
	"github.com/kitasuke/go-swift/utility"
)

const (
	BuildPathEnvKey   = "build_path"
	IsSkipBuildEnvKey = "is_skip_build"
	IsParallelEnvKey  = "is_parallel"
)

// ConfigModel ...
type ConfigModel struct {
	// Project Parameters
	BuildPath string

	// Test Run Configs
	isSkipBuild string
	isParallel  string
}

func (configs ConfigModel) print() {
	fmt.Println()

	log.Infof("Project Parameters:")
	log.Printf("- BuildPath: %s", configs.BuildPath)

	fmt.Println()
	log.Infof("Test Run Configs:")
	log.Printf("- IsSkipBuild: %s", configs.isSkipBuild)
	log.Printf("- IsParallel: %s", configs.isParallel)
}

func createConfigsModelFromEnvs() ConfigModel {
	return ConfigModel{
		// Project Parameters
		BuildPath: os.Getenv(BuildPathEnvKey),

		// Test Run Configs
		isSkipBuild: os.Getenv(IsSkipBuildEnvKey),
		isParallel:  os.Getenv(IsParallelEnvKey),
	}
}

func (configs ConfigModel) validate() error {
	if err := validateRequiredInputWithOptions(configs.isSkipBuild, IsSkipBuildEnvKey, []string{"yes", "no"}); err != nil {
		return err
	}

	if err := validateRequiredInputWithOptions(configs.isParallel, IsParallelEnvKey, []string{"yes", "no"}); err != nil {
		return err
	}

	return nil
}

//--------------------
// Functions
//--------------------

func validateRequiredInput(value, key string) error {
	if value == "" {
		return fmt.Errorf("Missing required input: %s", key)
	}
	return nil
}

func validateRequiredInputWithOptions(value, key string, options []string) error {
	if err := validateRequiredInput(value, key); err != nil {
		return err
	}

	found := false
	for _, option := range options {
		if option == value {
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("Invalid input: (%s) value: (%s), valid options: %s", key, value, strings.Join(options, ", "))
	}

	return nil
}

func failf(format string, v ...interface{}) {
	log.Errorf(format, v...)
	os.Exit(1)
}

//--------------------
// Main
//--------------------

func main() {
	configs := createConfigsModelFromEnvs()
	configs.print()
	if err := configs.validate(); err != nil {
		failf("Issue with input: %s", err)
	}

	fmt.Println()
	log.Infof("Other Configs:")

	isSkipBuild := configs.isSkipBuild == "yes"
	isParallel := configs.isParallel == "yes"

	swiftVersion, err := utility.GetSwiftVersion()
	if err != nil {
		failf("Failed to get the version of swift! Error: %s", err)
	}

	log.Printf("* swift_version: %s (%s)", swiftVersion.Version, swiftVersion.Target)

	fmt.Println()

	// setup CommandModel for test
	testCommandModel := swift.NewTestCommand()
	testCommandModel.SetBuildPath(configs.BuildPath)
	testCommandModel.SetSkipBuild(isSkipBuild)
	testCommandModel.SetIsParallel(isParallel)

	log.Infof("$ %s\n", testCommandModel.PrintableCmd())

	if err := testCommandModel.Run(); err != nil {
		failf("Test failed, error: %s", err)
	}
}
