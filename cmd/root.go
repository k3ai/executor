package cmd

import (
<<<<<<< HEAD
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"time"

	// "io/fs"
	"strings"

	"net/url"

	"github.com/go-git/go-git/v5"
	"github.com/spf13/cobra"
)

var directory = "/tmp/git/"

var (
	// Used for flags.
	strURI     string
	strBackend string
	strExtras  []string
	strEntry   string

	rootCmd = &cobra.Command{
		Use:   "run",
		Short: "K3ai executor for one-click applications",
		Run: func(cmd *cobra.Command, args []string) {
			getURI, _ := cmd.Flags().GetString("source")
			getBackend, _ := cmd.Flags().GetString("backend")
			getExtras, _ := cmd.Flags().GetStringArray("extras")
			getEntry, _ := cmd.Flags().GetString("entrypoint")
			getURI = strings.ToLower(getURI)
			getBackend = strings.ToLower(getBackend)
			getEntry = strings.ToLower(getEntry)

			if getURI != "" {
				boolURI, strHost, strPath := isValidUrl(getURI)
				var gitFolder string
				var gitFolderArr []string
				if boolURI {
					//strRepos will split the url in "http","","domain","owner","repo","everything else"
					strRepos := strings.Split(strHost+strPath, "/")
					// strRepos := strings.Split(strPath,"/")
					strRepo := append(strRepos, getEntry)
					lenRepo := len(strRepo)
					gitOwner := strRepo[3]
					gitRepo := strRepo[4]
					if lenRepo > 6 {
						for i := 4; i <= lenRepo-2; i++ {
							gitFolderArr = append(gitFolderArr, strRepos[i])
						}
						gitFolder = gitFolderArr[0]
					} else {
						gitFolder = strRepo[4]
					}

					strGit := strHost + "/" + gitOwner + "/" + gitRepo
					_ = gitClone(strGit, gitFolder)
					if getBackend == "mlflow" {
						if len(getExtras) > 0 {
							strCmd := getExtras[0]
							err := executionRun(strCmd)
							if err != nil {
								log.Println(err)
							}
						}
						if len(gitFolderArr) > 1 {
							gitFolder = strings.Join(gitFolderArr, "/")
						}
						strCmd := "sed -i '$ s/$/    - boto3/' /tmp/git/" + gitFolder + "/conda.yaml"
						err := executionRun(strCmd)
						if err != nil {
							log.Println(err)
						}

						time.Sleep(5 * time.Second)
						strCmd = "mlflow run /tmp/git/" + gitFolder + "/"
						err = executionRun(strCmd)
						if err != nil {
							log.Println(err)
						}

					}
					if getBackend == "kfp" {
						if len(gitFolderArr) > 1 {
							gitFolder = strings.Join(gitFolderArr, "/")
						}
						fileName := getEntry
						name := strconv.FormatInt(time.Now().UTC().UnixNano(), 10)
						strCmd := "dsl-compile --py /tmp/git/" + gitFolder + "/" + fileName + " --output /tmp/git/" + gitFolder + "/" + name + ".tar.gz"
						err := executionRun(strCmd)
						if err != nil {
							log.Println(err)
						}
						log.Println("..Preparing pipeline...")
						time.Sleep(3 * time.Second)
						log.Println("..Pushing pipeline...")
						time.Sleep(3 * time.Second)
						strCmd = "kfp --endpoint http://ml-pipeline.kubeflow.svc.cluster.local:8888 run submit -e default -r " + name + " -f /tmp/git/" + gitFolder + "/" + name + ".tar.gz"
						err = executionRun(strCmd)
						if err != nil {
							log.Println(err)
						}
					}

					strCmd := "rm -r /tmp/git/" + gitRepo
					err := executionRun(strCmd)
					if err != nil {
						log.Println(err)
					}

				} else {
					fmt.Println(" this is a local path")
				}
			}

=======
	// "os"
	// "fmt"

	"github.com/spf13/cobra"
)

const cliName = "k3ai"

var version bool

var (
	rootCmd = &cobra.Command{

		Use:   cliName + "[options]",
		Short: "K3ai is a very fast tool to run AI Infrastructure stacks",
		// By default (no Run/RunE in parent command) for typos in subcommands, cobra displays the help of parent command but exit(0) !
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return cmd.Help()
			}
			if version {
				return versionCommand().Execute()
			}
			_ = cmd.Help()
			return nil
>>>>>>> 540ac7528167d5b9f51fac2cdc50f0d7291ddfed
		},
	}
)

<<<<<<< HEAD
// Execute executes the root command.
=======
>>>>>>> 540ac7528167d5b9f51fac2cdc50f0d7291ddfed
func Execute() error {
	return rootCmd.Execute()
}

func init() {
<<<<<<< HEAD

	rootCmd.Flags().StringVarP(&strURI, "source", "s", "", "URI of the code, can be either a remote URL in the minimal form as [https://{domain}/{owner}/{user}]")
	rootCmd.Flags().StringVarP(&strBackend, "backend", "b", "", "Backend to be used may be only of one of the supported types [MLFLow,Kubeflow Pipelines (KFP),KFP DSL,Airflow,Argo,Tensorflow,Pytorch]")
	rootCmd.Flags().StringArrayVarP(&strExtras, "extras", "x", strExtras, "Represent the extra")
	rootCmd.Flags().StringVarP(&strEntry, "entrypoint", "e", "", "Entrypoint for KFP kind of workloads")

}

// isValidUrl tests a string to determine if it is a well-structured url or not.
func isValidUrl(toTest string) (boolURI bool, uri string, path string) {
	_, err := url.ParseRequestURI(toTest)
	if err != nil {
		return false, "", ""
	}

	u, err := url.Parse(toTest)
	if err != nil || u.Scheme == "" || u.Host == "" {
		return false, "", ""
	}
	uri = u.Scheme + "://" + u.Host
	return true, uri, u.Path
}

// Clone the repository into memory and return HEAD as a billy Filesystem.
func gitClone(url string, name string) error {
	os.MkdirAll(directory+"/"+name, 0755)
	_, err := git.PlainClone(directory+name, false, &git.CloneOptions{
		URL:               url,
		RecurseSubmodules: git.DefaultSubmoduleRecursionDepth,
	})
	if err != nil && err != git.ErrRepositoryAlreadyExists {
		log.Println(err)
	}

	return nil
}

func executionRun(strCmd string) error {
	cmd := exec.Command("/bin/bash", "--noprofile", "-c", "source /root/.profile && "+strCmd)
	r, err := cmd.StdoutPipe()
	cmd.Stderr = cmd.Stdout
	done := make(chan struct{})
	scanner := bufio.NewScanner(r)
	go func() {
		// Read line by line and process it
		for scanner.Scan() {
			line := scanner.Text()
			log.Println(line)
		}
		done <- struct{}{}
	}()
	// Start the command and check for errors
	_ = cmd.Start()
	<-done
	_ = cmd.Wait()
	return err
=======
	cobra.OnInitialize()
	rootCmd.CompletionOptions.DisableDefaultCmd = true
	rootCmd.Flags().SortFlags = false
	rootCmd.DisableFlagsInUseLine = true
	rootCmd.PersistentFlags().BoolP("help", "h", false, "Help usage")
	rootCmd.PersistentFlags().Lookup("help").Hidden = true
	rootCmd.AddCommand(
		upCommand(),
		downCommand(),
		clusterCommand(),
		pluginCommand(),
		runCommand(),
		versionCommand(),
	)
>>>>>>> 540ac7528167d5b9f51fac2cdc50f0d7291ddfed
}
