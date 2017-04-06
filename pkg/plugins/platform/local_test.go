package platform

import (
	"errors"
	"fmt"
	"log"
	"path"
	"testing"

	"os"

	"github.com/drud/ddev/pkg/testcommon"
	"github.com/drud/drud-go/utils/dockerutil"
	"github.com/drud/drud-go/utils/system"
	docker "github.com/fsouza/go-dockerclient"
	"github.com/stretchr/testify/assert"
)

var (
	localDBContainerName  = "local-%s-db"
	localWebContainerName = "local-%s-web"
	TestSites             = []testcommon.TestSite{
		{
			Name:      "drupal8",
			SourceURL: "https://github.com/drud/drupal8/archive/v0.2.2.tar.gz",
			FileURL:   "https://github.com/drud/drupal8/releases/download/v0.2.2/files.tar.gz",
			DBURL:     "https://github.com/drud/drupal8/releases/download/v0.2.2/db.tar.gz",
		},
		{
			Name:      "wp",
			SourceURL: "https://github.com/drud/wordpress/archive/v0.1.0.tar.gz",
			FileURL:   "https://github.com/drud/wordpress/releases/download/v0.1.0/files.tar.gz",
			DBURL:     "https://github.com/drud/wordpress/releases/download/v0.1.0/db.tar.gz",
		},
		{
			Name:      "kickstart",
			SourceURL: "https://github.com/drud/drupal-kickstart/archive/v0.1.0.tar.gz",
			FileURL:   "https://github.com/drud/drupal-kickstart/releases/download/v0.1.0/files.tar.gz",
			DBURL:     "https://github.com/drud/drupal-kickstart/releases/download/v0.1.0/db.tar.gz",
		},
	}
)

const netName = "ddev_default"

func TestMain(m *testing.M) {
	for i := range TestSites {
		TestSites[i].Prepare()
	}

	fmt.Println("Running tests.")
	testRun := m.Run()

	for i := range TestSites {
		TestSites[i].Cleanup()
	}

	os.Exit(testRun)
}

// ContainerCheck determines if a given container name exists and matches a given state
func ContainerCheck(checkName string, checkState string) (bool, error) {
	// ensure we have docker network
	client, _ := dockerutil.GetDockerClient()
	err := EnsureNetwork(client, netName)
	if err != nil {
		log.Fatal(err)
	}

	containers, err := client.ListContainers(docker.ListContainersOptions{All: true})
	if err != nil {
		log.Fatal(err)
	}

	for _, container := range containers {
		name := container.Names[0][1:]
		if name == checkName {
			if container.State == checkState {
				return true, nil
			}
			return false, errors.New("container " + name + " returned " + container.State)
		}
	}

	return false, errors.New("unable to find container " + checkName)
}

// TestLocalStart tests the functionality that is called when "ddev start" is executed
func TestLocalStart(t *testing.T) {

	// ensure we have docker network
	client, _ := dockerutil.GetDockerClient()
	err := EnsureNetwork(client, netName)
	if err != nil {
		log.Fatal(err)
	}

	assert := assert.New(t)
	app := PluginMap["local"]

	for _, site := range TestSites {
		webContainer := fmt.Sprintf(localWebContainerName, site.Name)
		dbContainer := fmt.Sprintf(localDBContainerName, site.Name)
		cleanup := site.Chdir()

		testcommon.ClearDockerEnv()
		app.Init(site.Dir)

		err = app.Start()
		assert.NoError(err)

		_, err = app.Wait()
		assert.NoError(err)

		// ensure docker-compose.yaml exists inside .ddev site folder
		composeFile := system.FileExists(app.DockerComposeYAMLPath())
		assert.True(composeFile)

		check, err := ContainerCheck(webContainer, "running")
		assert.NoError(err)
		assert.True(check)

		check, err = ContainerCheck(dbContainer, "running")
		assert.NoError(err)
		assert.True(check)

		cleanup()
	}
}

// TestLocalImportDB tests the functionality that is called when "ddev import-db" is executed
func TestLocalImportDB(t *testing.T) {
	assert := assert.New(t)
	app := PluginMap["local"]

	for _, site := range TestSites {
		cleanup := site.Chdir()
		dbPath := path.Join(os.TempDir(), "db.tar.gz")

		err := system.DownloadFile(dbPath, site.DBURL)
		assert.NoError(err)

		testcommon.ClearDockerEnv()
		app.Init(site.Dir)

		err = app.ImportDB(dbPath)
		assert.NoError(err)

		os.Remove(dbPath)

		cleanup()
	}
}

// TestLocalStop tests the functionality that is called when "ddev stop" is executed
func TestLocalStop(t *testing.T) {
	assert := assert.New(t)

	app := PluginMap["local"]

	for _, site := range TestSites {
		webContainer := fmt.Sprintf(localWebContainerName, site.Name)
		dbContainer := fmt.Sprintf(localDBContainerName, site.Name)
		cleanup := site.Chdir()

		testcommon.ClearDockerEnv()
		app.Init(site.Dir)

		err := app.Stop()
		assert.NoError(err)

		check, err := ContainerCheck(webContainer, "exited")
		assert.NoError(err)
		assert.True(check)

		check, err = ContainerCheck(dbContainer, "exited")
		assert.NoError(err)
		assert.True(check)

		cleanup()
	}
}

// TestLocalRemove tests the functionality that is called when "ddev rm" is executed
func TestLocalRemove(t *testing.T) {
	assert := assert.New(t)

	app := PluginMap["local"]

	for _, site := range TestSites {
		webContainer := fmt.Sprintf(localWebContainerName, site.Name)
		dbContainer := fmt.Sprintf(localDBContainerName, site.Name)
		cleanup := site.Chdir()

		testcommon.ClearDockerEnv()
		app.Init(site.Dir)

		// start the previously stopped containers -
		// stopped/removed have the same state
		err := app.Start()
		assert.NoError(err)

		_, err = app.Wait()
		assert.NoError(err)

		if err == nil {
			err = app.Down()
			assert.NoError(err)
		}

		check, err := ContainerCheck(webContainer, "running")
		assert.Error(err)
		assert.False(check)

		check, err = ContainerCheck(dbContainer, "running")
		assert.Error(err)
		assert.False(check)

		cleanup()
	}
}
