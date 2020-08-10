package devfile

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/openshift/odo/tests/helper"
	"github.com/openshift/odo/tests/integration/devfile/utils"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("odo devfile push command tests", func() {
	var namespace, context, cmpName, currentWorkingDirectory, originalKubeconfig string
	var sourcePath = "/projects"

	// Using program commmand according to cliRunner in devfile
	cliRunner := helper.GetCliRunner()

	// This is run after every Spec (It)
	var _ = BeforeEach(func() {
		SetDefaultEventuallyTimeout(10 * time.Minute)
		context = helper.CreateNewContext()
		os.Setenv("GLOBALODOCONFIG", filepath.Join(context, "config.yaml"))

		// Devfile push requires experimental mode to be set
		helper.CmdShouldPass("odo", "preference", "set", "Experimental", "true")

		originalKubeconfig = os.Getenv("KUBECONFIG")
		helper.LocalKubeconfigSet(context)
		namespace = cliRunner.CreateRandNamespaceProject()
		currentWorkingDirectory = helper.Getwd()
		cmpName = helper.RandString(6)
		helper.Chdir(context)
	})

	// Clean up after the test
	// This is run after every Spec (It)
	var _ = AfterEach(func() {
		cliRunner.DeleteNamespaceProject(namespace)
		helper.Chdir(currentWorkingDirectory)
		err := os.Setenv("KUBECONFIG", originalKubeconfig)
		Expect(err).NotTo(HaveOccurred())
		helper.DeleteDir(context)
		os.Unsetenv("GLOBALODOCONFIG")
	})

	Context("Pushing devfile without an .odo folder", func() {

		It("should be able to push based on metadata.name in devfile WITH a dash in the name", func() {
			// This is the name that's contained within `devfile-with-metadataname-foobar.yaml`
			name := "foobar"
			helper.CopyExample(filepath.Join("source", "devfiles", "springboot", "project"), context)
			helper.CopyExampleDevFile(filepath.Join("source", "devfiles", "springboot", "devfile-with-metadataname-foobar.yaml"), filepath.Join(context, "devfile.yaml"))

			output := helper.CmdShouldPass("odo", "push", "--namespace", namespace)
			Expect(output).To(ContainSubstring("Executing devfile commands for component " + name))
		})

		It("should be able to push based on name passed", func() {
			name := "springboot"
			helper.CopyExample(filepath.Join("source", "devfiles", "springboot", "project"), context)
			helper.CopyExampleDevFile(filepath.Join("source", "devfiles", "springboot", "devfile.yaml"), filepath.Join(context, "devfile.yaml"))

			output := helper.CmdShouldPass("odo", "push", "--namespace", namespace, name)
			Expect(output).To(ContainSubstring("Executing devfile commands for component " + name))
		})

	})

	Context("Verify devfile push works", func() {

		It("should have no errors when no endpoints within the devfile, should create a service when devfile has endpoints", func() {
			helper.CmdShouldPass("odo", "create", "nodejs", "--project", namespace, cmpName)

			helper.CopyExample(filepath.Join("source", "devfiles", "nodejs", "project"), context)
			helper.RenameFile("devfile.yaml", "devfile-old.yaml")
			helper.CopyExampleDevFile(filepath.Join("source", "devfiles", "nodejs", "devfile-no-endpoints.yaml"), filepath.Join(context, "devfile.yaml"))

			helper.CmdShouldPass("odo", "push", "--project", namespace)
			output := cliRunner.GetServices(namespace)
			Expect(output).NotTo(ContainSubstring(cmpName))

			helper.RenameFile("devfile-old.yaml", "devfile.yaml")
			output = helper.CmdShouldPass("odo", "push", "--project", namespace)

			Expect(output).To(ContainSubstring("Changes successfully pushed to component"))
			output = cliRunner.GetServices(namespace)
			Expect(output).To(ContainSubstring(cmpName))
		})

		It("checks that odo push works with a devfile", func() {
			helper.CmdShouldPass("odo", "create", "nodejs", "--project", namespace, cmpName)

			helper.CopyExample(filepath.Join("source", "devfiles", "nodejs", "project"), context)
			helper.CopyExampleDevFile(filepath.Join("source", "devfiles", "nodejs", "devfile.yaml"), filepath.Join(context, "devfile.yaml"))

			output := helper.CmdShouldPass("odo", "push", "--project", namespace)
			Expect(output).To(ContainSubstring("Changes successfully pushed to component"))

			// update devfile and push again
			helper.ReplaceString("devfile.yaml", "name: FOO", "name: BAR")
			helper.CmdShouldPass("odo", "push", "--project", namespace)
		})

		It("checks that odo push works with a devfile with sourcemapping set", func() {
			helper.CmdShouldPass("odo", "create", "nodejs", "--project", namespace, cmpName)

			helper.CopyExample(filepath.Join("source", "devfiles", "nodejs", "project"), context)
			helper.CopyExampleDevFile(filepath.Join("source", "devfiles", "nodejs", "devfileSourceMapping.yaml"), filepath.Join(context, "devfile.yaml"))

			output := helper.CmdShouldPass("odo", "push", "--project", namespace)
			Expect(output).To(ContainSubstring("Changes successfully pushed to component"))

			// Verify source code was synced to /test instead of /projects
			var statErr error
			podName := cliRunner.GetRunningPodNameByComponent(cmpName, namespace)
			cliRunner.CheckCmdOpInRemoteDevfilePod(
				podName,
				"runtime",
				namespace,
				[]string{"stat", "/test/server.js"},
				func(cmdOp string, err error) bool {
					statErr = err
					return true
				},
			)
			Expect(statErr).ToNot(HaveOccurred())
		})

		It("checks that odo push works with a devfile with composite commands", func() {
			helper.CmdShouldPass("odo", "create", "nodejs", "--project", namespace, cmpName)

			helper.CopyExample(filepath.Join("source", "devfiles", "nodejs", "project"), context)
			helper.CopyExampleDevFile(filepath.Join("source", "devfiles", "nodejs", "devfileCompositeCommands.yaml"), filepath.Join(context, "devfile.yaml"))

			output := helper.CmdShouldPass("odo", "push", "--context", context)
			Expect(output).To(ContainSubstring("Executing mkdir command"))

			// Verify the command executed successfully
			var statErr error
			podName := cliRunner.GetRunningPodNameByComponent(cmpName, namespace)
			cliRunner.CheckCmdOpInRemoteDevfilePod(
				podName,
				"runtime",
				namespace,
				[]string{"stat", "/projects/testfolder"},
				func(cmdOp string, err error) bool {
					statErr = err
					return true
				},
			)
			Expect(statErr).ToNot(HaveOccurred())
		})

		It("checks that odo push works with a devfile with parallel composite commands", func() {
			helper.CmdShouldPass("odo", "create", "nodejs", "--project", namespace, cmpName)

			helper.CopyExample(filepath.Join("source", "devfiles", "nodejs", "project"), context)
			helper.CopyExampleDevFile(filepath.Join("source", "devfiles", "nodejs", "devfileCompositeCommandsParallel.yaml"), filepath.Join(context, "devfile.yaml"))

			output := helper.CmdShouldPass("odo", "push", "--build-command", "buildAndMkdir", "--context", context)
			Expect(output).To(ContainSubstring("Executing mkdir command"))

			// Verify the command executed successfully
			var statErr error
			podName := cliRunner.GetRunningPodNameByComponent(cmpName, namespace)
			cliRunner.CheckCmdOpInRemoteDevfilePod(
				podName,
				"runtime",
				namespace,
				[]string{"stat", "/projects/testfolder"},
				func(cmdOp string, err error) bool {
					statErr = err
					return true
				},
			)
			Expect(statErr).ToNot(HaveOccurred())
		})

		It("checks that odo push works with a devfile with nested composite commands", func() {
			helper.CmdShouldPass("odo", "create", "nodejs", "--project", namespace, cmpName)

			helper.CopyExample(filepath.Join("source", "devfiles", "nodejs", "project"), context)
			helper.CopyExampleDevFile(filepath.Join("source", "devfiles", "nodejs", "devfileNestedCompCommands.yaml"), filepath.Join(context, "devfile.yaml"))

			// Verify nested command was executed
			output := helper.CmdShouldPass("odo", "push", "--context", context)
			Expect(output).To(ContainSubstring("Executing mkdir command"))

			// Verify the command executed successfully
			var statErr error
			podName := cliRunner.GetRunningPodNameByComponent(cmpName, namespace)
			cliRunner.CheckCmdOpInRemoteDevfilePod(
				podName,
				"runtime",
				namespace,
				[]string{"stat", "/projects/testfolder"},
				func(cmdOp string, err error) bool {
					statErr = err
					return true
				},
			)
			Expect(statErr).ToNot(HaveOccurred())
		})

		It("should throw a validation error for composite run commands", func() {
			helper.CmdShouldPass("odo", "create", "nodejs", "--project", namespace, cmpName)

			helper.CopyExample(filepath.Join("source", "devfiles", "nodejs", "project"), context)
			helper.CopyExampleDevFile(filepath.Join("source", "devfiles", "nodejs", "devfileCompositeRun.yaml"), filepath.Join(context, "devfile.yaml"))

			// Verify odo push failed
			output := helper.CmdShouldFail("odo", "push", "--context", context)
			Expect(output).To(ContainSubstring("composite commands of run Kind are not supported currently"))
		})

		It("should throw a validation error for composite command referencing non-existent commands", func() {
			helper.CmdShouldPass("odo", "create", "nodejs", "--project", namespace, cmpName)

			helper.CopyExample(filepath.Join("source", "devfiles", "nodejs", "project"), context)
			helper.CopyExampleDevFile(filepath.Join("source", "devfiles", "nodejs", "devfileCompositeNonExistent.yaml"), filepath.Join(context, "devfile.yaml"))

			// Verify odo push failed
			output := helper.CmdShouldFail("odo", "push", "--context", context)
			Expect(output).To(ContainSubstring("does not exist in the devfile"))
		})

		It("should throw a validation error for composite command indirectly referencing itself", func() {
			helper.CmdShouldPass("odo", "create", "nodejs", "--project", namespace, cmpName)

			helper.CopyExample(filepath.Join("source", "devfiles", "nodejs", "project"), context)
			helper.CopyExampleDevFile(filepath.Join("source", "devfiles", "nodejs", "devfileIndirectNesting.yaml"), filepath.Join(context, "devfile.yaml"))

			// Verify odo push failed
			output := helper.CmdShouldFail("odo", "push", "--context", context)
			Expect(output).To(ContainSubstring("cannot indirectly reference itself"))
		})

		It("should throw a validation error for composite command that has invalid exec subcommand", func() {
			helper.CmdShouldPass("odo", "create", "nodejs", "--project", namespace, cmpName)

			helper.CopyExample(filepath.Join("source", "devfiles", "nodejs", "project"), context)
			helper.CopyExampleDevFile(filepath.Join("source", "devfiles", "nodejs", "devfileCompositeInvalidComponent.yaml"), filepath.Join(context, "devfile.yaml"))

			// Verify odo push failed
			output := helper.CmdShouldFail("odo", "push", "--context", context)
			Expect(output).To(ContainSubstring("references an invalid command"))
		})

		It("checks that odo push works outside of the context directory", func() {
			helper.Chdir(currentWorkingDirectory)

			helper.CmdShouldPass("odo", "create", "nodejs", "--project", namespace, "--context", context, cmpName)

			helper.CopyExample(filepath.Join("source", "devfiles", "nodejs", "project"), context)
			helper.CopyExampleDevFile(filepath.Join("source", "devfiles", "nodejs", "devfile.yaml"), filepath.Join(context, "devfile.yaml"))

			output := helper.CmdShouldPass("odo", "push", "--context", context)
			Expect(output).To(ContainSubstring("Changes successfully pushed to component"))
		})

		It("should not build when no changes are detected in the directory and build when a file change is detected", func() {
			utils.ExecPushToTestFileChanges(context, cmpName, namespace)
		})

		It("checks that odo push with -o json displays machine readable JSON event output", func() {

			helper.CmdShouldPass("odo", "create", "nodejs", "--project", namespace, cmpName)

			helper.CopyExample(filepath.Join("source", "devfiles", "nodejs", "project"), context)
			helper.CopyExampleDevFile(filepath.Join("source", "devfiles", "nodejs", "devfile.yaml"), filepath.Join(context, "devfile.yaml"))

			output := helper.CmdShouldPass("odo", "push", "-o", "json", "--project", namespace)
			utils.AnalyzePushConsoleOutput(output)

			// update devfile and push again
			helper.ReplaceString("devfile.yaml", "name: FOO", "name: BAR")
			output = helper.CmdShouldPass("odo", "push", "-o", "json", "--project", namespace)
			utils.AnalyzePushConsoleOutput(output)

		})

		It("should be able to create a file, push, delete, then push again propagating the deletions", func() {
			newFilePath := filepath.Join(context, "foobar.txt")
			newDirPath := filepath.Join(context, "testdir")
			utils.ExecPushWithNewFileAndDir(context, cmpName, namespace, newFilePath, newDirPath)

			// Check to see if it's been pushed (foobar.txt abd directory testdir)
			podName := cliRunner.GetRunningPodNameByComponent(cmpName, namespace)

			stdOut := cliRunner.ExecListDir(podName, namespace, sourcePath)
			helper.MatchAllInOutput(stdOut, []string{"foobar.txt", "testdir"})

			// Now we delete the file and dir and push
			helper.DeleteDir(newFilePath)
			helper.DeleteDir(newDirPath)
			helper.CmdShouldPass("odo", "push", "--project", namespace, "-v4")

			// Then check to see if it's truly been deleted
			stdOut = cliRunner.ExecListDir(podName, namespace, sourcePath)
			helper.DontMatchAllInOutput(stdOut, []string{"foobar.txt", "testdir"})
		})

		It("should delete the files from the container if its removed locally", func() {
			helper.CmdShouldPass("odo", "create", "nodejs", "--project", namespace, cmpName)

			helper.CopyExample(filepath.Join("source", "devfiles", "nodejs", "project"), context)
			helper.CopyExampleDevFile(filepath.Join("source", "devfiles", "nodejs", "devfile.yaml"), filepath.Join(context, "devfile.yaml"))

			helper.CmdShouldPass("odo", "push", "--project", namespace)

			// Check to see if it's been pushed (foobar.txt abd directory testdir)
			podName := cliRunner.GetRunningPodNameByComponent(cmpName, namespace)

			var statErr error
			cliRunner.CheckCmdOpInRemoteDevfilePod(
				podName,
				"",
				namespace,
				[]string{"stat", "/projects/server.js"},
				func(cmdOp string, err error) bool {
					statErr = err
					return true
				},
			)
			Expect(statErr).ToNot(HaveOccurred())
			Expect(os.Remove(filepath.Join(context, "server.js"))).NotTo(HaveOccurred())
			helper.CmdShouldPass("odo", "push", "--project", namespace)

			cliRunner.CheckCmdOpInRemoteDevfilePod(
				podName,
				"",
				namespace,
				[]string{"stat", "/projects/server.js"},
				func(cmdOp string, err error) bool {
					statErr = err
					return true
				},
			)
			Expect(statErr).To(HaveOccurred())
			Expect(statErr.Error()).To(ContainSubstring("cannot stat '/projects/server.js': No such file or directory"))
		})

		It("should build when no changes are detected in the directory and force flag is enabled", func() {
			utils.ExecPushWithForceFlag(context, cmpName, namespace)
		})

		It("should execute the default build and run command groups if present", func() {
			utils.ExecDefaultDevfileCommands(context, cmpName, namespace)

			// Check to see if it's been pushed (foobar.txt abd directory testdir)
			podName := cliRunner.GetRunningPodNameByComponent(cmpName, namespace)

			var statErr error
			var cmdOutput string
			cliRunner.CheckCmdOpInRemoteDevfilePod(
				podName,
				"runtime",
				namespace,
				[]string{"ps", "-ef"},
				func(cmdOp string, err error) bool {
					cmdOutput = cmdOp
					statErr = err
					return true
				},
			)
			Expect(statErr).ToNot(HaveOccurred())
			Expect(cmdOutput).To(ContainSubstring("/myproject/app.jar"))
		})

		It("should execute PostStart commands if present and not execute when component already exists", func() {
			helper.CmdShouldPass("odo", "create", "nodejs", "--project", namespace, cmpName)

			helper.CopyExample(filepath.Join("source", "devfiles", "nodejs", "project"), context)
			helper.CopyExampleDevFile(filepath.Join("source", "devfiles", "nodejs", "devfile-with-valid-events.yaml"), filepath.Join(context, "devfile.yaml"))

			output := helper.CmdShouldPass("odo", "push", "--namespace", namespace)
			helper.MatchAllInOutput(output, []string{"Executing mypoststart command \"echo I am a PostStart\"", "Executing secondpoststart command \"echo I am also a PostStart\""})

			// Need to force so build and run get triggered again with the component already created.
			output = helper.CmdShouldPass("odo", "push", "--namespace", namespace, "-f")
			helper.DontMatchAllInOutput(output, []string{"Executing mypoststart command \"echo I am a PostStart\"", "Executing secondpoststart command \"echo I am also a PostStart\""})
			helper.MatchAllInOutput(output, []string{
				"Executing devbuild command",
				"Executing devrun command",
			})
		})

		It("should err out on an event not mentioned in the devfile commands", func() {
			helper.CmdShouldPass("odo", "create", "nodejs", "--project", namespace, cmpName)

			helper.CopyExample(filepath.Join("source", "devfiles", "nodejs", "project"), context)
			helper.CopyExampleDevFile(filepath.Join("source", "devfiles", "nodejs", "devfile-with-valid-events.yaml"), filepath.Join(context, "devfile.yaml"))

			helper.ReplaceString("devfile.yaml", "secondpoststart", "secondpoststart12345")

			output := helper.CmdShouldFail("odo", "push", "--namespace", namespace)
			helper.MatchAllInOutput(output, []string{"does not map to a valid devfile command"})
		})

		It("should err out on an event command not mapping to a devfile container component", func() {
			helper.CmdShouldPass("odo", "create", "nodejs", "--project", namespace, cmpName)

			helper.CopyExample(filepath.Join("source", "devfiles", "nodejs", "project"), context)
			helper.CopyExampleDevFile(filepath.Join("source", "devfiles", "nodejs", "devfile-with-valid-events.yaml"), filepath.Join(context, "devfile.yaml"))

			helper.ReplaceString("devfile.yaml", "secondpoststart", "wrongPostStart")

			output := helper.CmdShouldFail("odo", "push", "--namespace", namespace)
			helper.MatchAllInOutput(output, []string{"the command does not map to a supported component"})
		})

		It("should err out on an event composite command mentioning an invalid child command", func() {
			helper.CmdShouldPass("odo", "create", "nodejs", "--project", namespace, cmpName)

			helper.CopyExample(filepath.Join("source", "devfiles", "nodejs", "project"), context)
			helper.CopyExampleDevFile(filepath.Join("source", "devfiles", "nodejs", "devfile-with-valid-events.yaml"), filepath.Join(context, "devfile.yaml"))

			helper.ReplaceString("devfile.yaml", "secondpoststart", "myWrongCompCmd")

			output := helper.CmdShouldFail("odo", "push", "--namespace", namespace)
			helper.MatchAllInOutput(output, []string{"does not exist in the devfile"})
		})

		It("should be able to handle a missing build command group", func() {
			utils.ExecWithMissingBuildCommand(context, cmpName, namespace)
		})

		It("should error out on a missing run command group", func() {
			utils.ExecWithMissingRunCommand(context, cmpName, namespace)
		})

		It("should be able to push using the custom commands", func() {
			utils.ExecWithCustomCommand(context, cmpName, namespace)
		})

		It("should error out on a wrong custom commands", func() {
			utils.ExecWithWrongCustomCommand(context, cmpName, namespace)
		})

		It("should error out on multiple or no default commands", func() {
			utils.ExecWithMultipleOrNoDefaults(context, cmpName, namespace)
		})

		It("should execute commands with flags if there are more than one default command", func() {
			utils.ExecMultipleDefaultsWithFlags(context, cmpName, namespace)
		})

		It("should execute commands with flags if the command has no group kind", func() {
			utils.ExecCommandWithoutGroupUsingFlags(context, cmpName, namespace)
		})

		It("should error out if the devfile has an invalid command group", func() {
			utils.ExecWithInvalidCommandGroup(context, cmpName, namespace)
		})

		It("should not restart the application if restart is false", func() {
			utils.ExecWithRestartAttribute(context, cmpName, namespace)
		})

		It("should create pvc and reuse if it shares the same devfile volume name", func() {
			helper.CmdShouldPass("odo", "create", "nodejs", "--project", namespace, cmpName)

			helper.CopyExample(filepath.Join("source", "devfiles", "nodejs", "project"), context)
			helper.CopyExampleDevFile(filepath.Join("source", "devfiles", "nodejs", "devfile-with-volumes.yaml"), filepath.Join(context, "devfile.yaml"))

			output := helper.CmdShouldPass("odo", "push", "--namespace", namespace)
			helper.MatchAllInOutput(output, []string{
				"Executing devbuild command",
				"Executing devrun command",
			})

			// Check to see if it's been pushed (foobar.txt abd directory testdir)
			podName := cliRunner.GetRunningPodNameByComponent(cmpName, namespace)

			var statErr error
			var cmdOutput string

			cliRunner.CheckCmdOpInRemoteDevfilePod(
				podName,
				"runtime2",
				namespace,
				[]string{"cat", "/data/myfile.log"},
				func(cmdOp string, err error) bool {
					cmdOutput = cmdOp
					statErr = err
					return true
				},
			)
			Expect(statErr).ToNot(HaveOccurred())
			Expect(cmdOutput).To(ContainSubstring("hello"))

			cliRunner.CheckCmdOpInRemoteDevfilePod(
				podName,
				"runtime2",
				namespace,
				[]string{"stat", "/data2"},
				func(cmdOp string, err error) bool {
					statErr = err
					return true
				},
			)
			Expect(statErr).ToNot(HaveOccurred())

			volumesMatched := false

			// check the volume name and mount paths for the containers
			volNamesAndPaths := cliRunner.GetVolumeMountNamesandPathsFromContainer(cmpName, "runtime", namespace)
			volNamesAndPathsArr := strings.Fields(volNamesAndPaths)
			for _, volNamesAndPath := range volNamesAndPathsArr {
				volNamesAndPathArr := strings.Split(volNamesAndPath, ":")

				if strings.Contains(volNamesAndPathArr[0], "myvol") && volNamesAndPathArr[1] == "/data" {
					volumesMatched = true
				}
			}
			Expect(volumesMatched).To(Equal(true))
		})
	})

	Context("Verify devfile volume components work", func() {

		It("should error out when duplicate volume components exist", func() {
			helper.CmdShouldPass("odo", "create", "nodejs", "--project", namespace, cmpName)

			helper.CopyExample(filepath.Join("source", "devfiles", "nodejs", "project"), context)
			helper.RenameFile("devfile.yaml", "devfile-old.yaml")
			helper.CopyExampleDevFile(filepath.Join("source", "devfiles", "nodejs", "devfile-with-volume-components.yaml"), filepath.Join(context, "devfile.yaml"))

			helper.ReplaceString("devfile.yaml", "secondvol", "firstvol")

			output := helper.CmdShouldFail("odo", "push", "--project", namespace)
			Expect(output).To(ContainSubstring("duplicate volume components present in devfile"))
		})

		It("should error out when a wrong volume size is used", func() {
			helper.CmdShouldPass("odo", "create", "nodejs", "--project", namespace, cmpName)

			helper.CopyExample(filepath.Join("source", "devfiles", "nodejs", "project"), context)
			helper.RenameFile("devfile.yaml", "devfile-old.yaml")
			helper.CopyExampleDevFile(filepath.Join("source", "devfiles", "nodejs", "devfile-with-volume-components.yaml"), filepath.Join(context, "devfile.yaml"))

			helper.ReplaceString("devfile.yaml", "3Gi", "3Garbage")

			output := helper.CmdShouldFail("odo", "push", "--project", namespace)
			Expect(output).To(ContainSubstring("quantities must match the regular expression"))
		})

		It("should error out if a container component has volume mount that does not refer a valid volume component", func() {
			helper.CmdShouldPass("odo", "create", "nodejs", "--project", namespace, cmpName)

			helper.CopyExample(filepath.Join("source", "devfiles", "nodejs", "project"), context)
			helper.RenameFile("devfile.yaml", "devfile-old.yaml")
			helper.CopyExampleDevFile(filepath.Join("source", "devfiles", "nodejs", "devfile-with-invalid-volmount.yaml"), filepath.Join(context, "devfile.yaml"))

			output := helper.CmdShouldFail("odo", "push", "--project", namespace)
			Expect(output).To(ContainSubstring("unable to find volume mount"))
		})

		It("should successfully use the volume components in container components", func() {
			helper.CmdShouldPass("odo", "create", "nodejs", "--project", namespace, cmpName)

			helper.CopyExample(filepath.Join("source", "devfiles", "nodejs", "project"), context)
			helper.RenameFile("devfile.yaml", "devfile-old.yaml")
			helper.CopyExampleDevFile(filepath.Join("source", "devfiles", "nodejs", "devfile-with-volume-components.yaml"), filepath.Join(context, "devfile.yaml"))

			output := helper.CmdShouldPass("odo", "push", "--project", namespace)
			Expect(output).To(ContainSubstring("Changes successfully pushed to component"))

			// Verify the pvc size for firstvol
			storageSize := cliRunner.GetPVCSize(cmpName, "firstvol", namespace)
			// should be the default size
			Expect(storageSize).To(ContainSubstring("5Gi"))

			// Verify the pvc size for secondvol
			storageSize = cliRunner.GetPVCSize(cmpName, "secondvol", namespace)
			// should be the specified size in the devfile volume component
			Expect(storageSize).To(ContainSubstring("3Gi"))
		})

	})

	Context("when .gitignore file exists", func() {
		It("checks that .odo/env exists in gitignore", func() {
			helper.CmdShouldPass("odo", "create", "nodejs", "--project", namespace, cmpName)

			ignoreFilePath := filepath.Join(context, ".gitignore")

			helper.FileShouldContainSubstring(ignoreFilePath, filepath.Join(".odo", "env"))

		})
	})

	Context("exec commands with environment variables", func() {
		It("Should be able to exec command with single environment variable", func() {
			helper.CmdShouldPass("odo", "create", "nodejs", "--project", namespace, cmpName)
			helper.CopyExample(filepath.Join("source", "devfiles", "nodejs", "project"), context)
			helper.CopyExampleDevFile(filepath.Join("source", "devfiles", "nodejs", "devfile-with-multiple-defaults.yaml"), filepath.Join(context, "devfile.yaml"))
			output := helper.CmdShouldPass("odo", "push", "--build-command", "firstbuild", "--run-command", "singleenv", "--namespace", namespace, "--context", context)
			Expect(output).To(ContainSubstring("mkdir $ENV1"))

			podName := cliRunner.GetRunningPodNameByComponent(cmpName, namespace)
			output = cliRunner.ExecListDir(podName, namespace, sourcePath)
			Expect(output).To(ContainSubstring("test_env_variable"))

		})

		It("Should be able to exec command with multiple environment variables", func() {
			helper.CmdShouldPass("odo", "create", "nodejs", "--project", namespace, cmpName)
			helper.CopyExample(filepath.Join("source", "devfiles", "nodejs", "project"), context)
			helper.CopyExampleDevFile(filepath.Join("source", "devfiles", "nodejs", "devfile-with-multiple-defaults.yaml"), filepath.Join(context, "devfile.yaml"))
			output := helper.CmdShouldPass("odo", "push", "--build-command", "firstbuild", "--run-command", "multipleenv", "--namespace", namespace, "--context", context)
			Expect(output).To(ContainSubstring("mkdir $ENV1 $ENV2"))

			podName := cliRunner.GetRunningPodNameByComponent(cmpName, namespace)
			output = cliRunner.ExecListDir(podName, namespace, sourcePath)
			helper.MatchAllInOutput(output, []string{"test_env_variable1", "test_env_variable2"})

		})

		It("Should be able to exec command with environment variable with spaces", func() {
			helper.CmdShouldPass("odo", "create", "nodejs", "--project", namespace, cmpName)
			helper.CopyExample(filepath.Join("source", "devfiles", "nodejs", "project"), context)
			helper.CopyExampleDevFile(filepath.Join("source", "devfiles", "nodejs", "devfile-with-multiple-defaults.yaml"), filepath.Join(context, "devfile.yaml"))
			output := helper.CmdShouldPass("odo", "push", "--build-command", "firstbuild", "--run-command", "envwithspace", "--namespace", namespace, "--context", context)
			Expect(output).To(ContainSubstring("mkdir \\\"$ENV1\\\""))

			podName := cliRunner.GetRunningPodNameByComponent(cmpName, namespace)
			output = cliRunner.ExecListDir(podName, namespace, sourcePath)
			helper.MatchAllInOutput(output, []string{"env with space"})

		})
	})

	Context("push with listing the devfile component", func() {

		It("checks components in a specific app and all apps", func() {

			// component created in "app" application
			helper.CmdShouldPass("odo", "create", "nodejs", "--project", namespace, "--context", context, cmpName)

			helper.CopyExample(filepath.Join("source", "devfiles", "nodejs", "project"), context)
			helper.CopyExampleDevFile(filepath.Join("source", "devfiles", "nodejs", "devfile.yaml"), filepath.Join(context, "devfile.yaml"))
			output := helper.CmdShouldPass("odo", "list", "--context", context)
			Expect(helper.Suffocate(output)).To(ContainSubstring(helper.Suffocate(fmt.Sprintf("%s%s%s%sUnpushed", "app", cmpName, namespace, "nodejs"))))

			output = helper.CmdShouldPass("odo", "push", "--context", context)
			Expect(output).To(ContainSubstring("Changes successfully pushed to component"))

			// component created in different application
			context2 := helper.CreateNewContext()
			cmpName2 := helper.RandString(6)
			appName := helper.RandString(6)

			helper.CmdShouldPass("odo", "create", "nodejs", "--project", namespace, "--app", appName, "--context", context2, cmpName2)

			helper.CopyExample(filepath.Join("source", "devfiles", "nodejs", "project"), context2)
			helper.CopyExampleDevFile(filepath.Join("source", "devfiles", "nodejs", "devfile.yaml"), filepath.Join(context2, "devfile.yaml"))

			output = helper.CmdShouldPass("odo", "list", "--context", context2)
			Expect(helper.Suffocate(output)).To(ContainSubstring(helper.Suffocate(fmt.Sprintf("%s%s%s%sUnpushed", appName, cmpName2, namespace, "nodejs"))))
			output2 := helper.CmdShouldPass("odo", "push", "--context", context2)
			Expect(output2).To(ContainSubstring("Changes successfully pushed to component"))

			output = helper.CmdShouldPass("odo", "list", "--project", namespace)
			Expect(output).To(ContainSubstring(cmpName))
			Expect(output).ToNot(ContainSubstring(cmpName2))

			output = helper.CmdShouldPass("odo", "list", "--all-apps", "--project", namespace)

			Expect(output).To(ContainSubstring(cmpName))
			Expect(output).To(ContainSubstring(cmpName2))

			helper.CmdShouldPass("odo", "preference", "set", "Experimental", "false")
			helper.DeleteDir(context2)

		})

		It("checks devfile and s2i components together", func() {
			if os.Getenv("KUBERNETES") == "true" {
				Skip("Skipping test because s2i image is not supported on Kubernetes cluster")
			}

			// component created in "app" application
			helper.CmdShouldPass("odo", "create", "nodejs", "--project", namespace, "--context", context, cmpName)

			helper.CopyExample(filepath.Join("source", "devfiles", "nodejs", "project"), context)
			helper.CopyExampleDevFile(filepath.Join("source", "devfiles", "nodejs", "devfile.yaml"), filepath.Join(context, "devfile.yaml"))

			output := helper.CmdShouldPass("odo", "list", "--context", context)
			Expect(helper.Suffocate(output)).To(ContainSubstring(helper.Suffocate(fmt.Sprintf("%s%s%s%sUnpushed", "app", cmpName, namespace, "nodejs"))))

			output = helper.CmdShouldPass("odo", "push", "--context", context)
			Expect(output).To(ContainSubstring("Changes successfully pushed to component"))

			// component created in different application
			context2 := helper.CreateNewContext()
			cmpName2 := helper.RandString(6)
			appName := helper.RandString(6)
			helper.CmdShouldPass("odo", "preference", "set", "--force", "Experimental", "false")
			helper.CopyExample(filepath.Join("source", "nodejs"), context2)
			helper.CmdShouldPass("odo", "create", "nodejs", "--project", namespace, "--app", appName, "--context", context2, cmpName2)

			output2 := helper.CmdShouldPass("odo", "push", "--context", context2)
			Expect(output2).To(ContainSubstring("Changes successfully pushed to component"))

			helper.CmdShouldPass("odo", "preference", "set", "--force", "Experimental", "true")

			output = helper.CmdShouldPass("odo", "list", "--all-apps", "--project", namespace)

			Expect(output).To(ContainSubstring(cmpName))
			Expect(output).To(ContainSubstring(cmpName2))

			output = helper.CmdShouldPass("odo", "list", "--app", appName, "--project", namespace)
			Expect(output).To(Not(ContainSubstring(cmpName))) // cmpName component hasn't been created under appName
			Expect(output).To(ContainSubstring(cmpName2))

			helper.DeleteDir(context2)
		})

	})

	/*
		Disabled test due to issue https://github.com/openshift/odo/issues/3638

		Context("Handle devfiles with parent", func() {
			It("should handle a devfile with a parent and add a extra command", func() {
				utils.ExecPushToTestParent(context, cmpName, namespace)
				podName := cliRunner.GetRunningPodNameByComponent(cmpName, namespace)
				listDir := cliRunner.ExecListDir(podName, namespace, "/projects/nodejs-starter")
				Expect(listDir).To(ContainSubstring("blah.js"))
			})

			It("should handle a parent and override/append it's envs", func() {
				utils.ExecPushWithParentOverride(context, cmpName, namespace)

				envMap := cliRunner.GetEnvsDevFileDeployment(cmpName, namespace)

				value, ok := envMap["MODE2"]
				Expect(ok).To(BeTrue())
				Expect(value).To(Equal("TEST2-override"))

				value, ok = envMap["myprop-3"]
				Expect(ok).To(BeTrue())
				Expect(value).To(Equal("myval-3"))

				value, ok = envMap["myprop2"]
				Expect(ok).To(BeTrue())
				Expect(value).To(Equal("myval2"))
			})


				It("should handle a multi layer parent", func() {
					utils.ExecPushWithMultiLayerParent(context, cmpName, namespace)

					podName := cliRunner.GetRunningPodNameByComponent(cmpName, namespace)
					listDir := cliRunner.ExecListDir(podName, namespace, "/projects/user-app")
					helper.MatchAllInOutput(listDir, []string{"blah.js", "new-blah.js"})

					envMap := cliRunner.GetEnvsDevFileDeployment(cmpName, namespace)

					value, ok := envMap["MODE2"]
					Expect(ok).To(BeTrue())
					Expect(value).To(Equal("TEST2-override"))

					value, ok = envMap["myprop3"]
					Expect(ok).To(BeTrue())
					Expect(value).To(Equal("myval3"))

					value, ok = envMap["myprop2"]
					Expect(ok).To(BeTrue())
					Expect(value).To(Equal("myval2"))

					value, ok = envMap["myprop4"]
					Expect(ok).To(BeTrue())
					Expect(value).To(Equal("myval4"))
				})
		})
	*/
})
