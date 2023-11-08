//go:generate go install -v github.com/josephspurrier/goversioninfo/cmd/goversioninfo
//go:generate goversioninfo -icon=res/papp.ico -manifest=res/papp.manifest
package main

import (
	"os"

	"github.com/portapps/portapps/v3"
	"github.com/portapps/portapps/v3/pkg/log"
	"github.com/portapps/portapps/v3/pkg/utl"
)

var (
	app *portapps.App
)

func init() {
	var err error

	// Init app
	if app, err = portapps.New("sqlitestudio-portable", "SQLiteStudio"); err != nil {
		log.Fatal().Err(err).Msg("Cannot initialize application. See log file for more info.")
	}
}


func main() {
	utl.CreateFolder(app.DataPath)
	app.Process = utl.PathJoin(app.AppPath, "SQLiteStudio.exe")

	// SQLiteStudio paths	
	sqlitestudioCfg := utl.PathJoin(app.AppPath, "sqlitestudio-cfg")

	// Copy existing files from data to roaming folder for the current user
	utl.CreateFolder(sqlitestudioCfg)
	if _, err := os.Stat(app.DataPath); err == nil {
		if err := utl.CopyFolder(app.DataPath, sqlitestudioCfg); err != nil {
			log.Error().Err(err).Msgf("Cannot copy %s", app.DataPath)
		}
	}

	// On exit
	defer func() {

		if _, err := os.Stat(sqlitestudioCfg); err == nil {
			if err = utl.CopyFolder(sqlitestudioCfg, app.DataPath); err != nil {
				log.Warn().Err(err).Msgf("Cannot copy back %s", sqlitestudioCfg)
			}
		}
	}()
	
	defer app.Close()
	app.Launch(os.Args[1:])
}
