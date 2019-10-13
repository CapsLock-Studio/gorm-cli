package migrate

import (
	"os"
	"os/exec"
	"plugin"
	"strings"

	"github.com/pkg/errors"
)

func buildPlugin(goFileName string) (string, error) {
	pluginName := strings.Replace(goFileName, ".go", ".so", -1)
	cmd := exec.Command(
		"go",
		"build",
		"-buildmode=plugin",
		"-o",
		MigrationTargetFolder+"/"+pluginName,
		MigrationTargetFolder+"/"+goFileName,
	)

	if err := cmd.Run(); err != nil {
		return "", err
	}

	return pluginName, nil
}

func getMigration(pluginName string) (Migration, error) {
	plug, err := plugin.Open(MigrationTargetFolder + "/" + pluginName)
	if err != nil {
		return nil, errors.Wrap(err, "Open plugin filed.")
	}

	symMigration, err := plug.Lookup("Migration")
	if err != nil {
		return nil, errors.Wrap(err, "Migration wrong format - missing Migration declaration.")
	}

	var migration Migration
	migration, ok := symMigration.(Migration)
	if !ok {
		return nil, errors.Wrap(err, "Unexpected type from module symbol.")
	}
	return migration, nil
}

func removePlugin(goFileName string) error {
	pluginName := strings.Replace(goFileName, ".go", ".so", -1)
	return os.Remove(MigrationTargetFolder + "/" + pluginName)
}

// MigrationTargetFolder is the migration folder target.
var MigrationTargetFolder string
