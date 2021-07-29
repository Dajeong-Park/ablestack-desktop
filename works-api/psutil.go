package main

import (
	"fmt"
	"github.com/ycyun/go-powershell"
	"os"
	"strings"
)

type PSCMD struct {
	CMD string `uri:"cmd" binding:"required"`
	ARG string `uri:"arg" binding:"required"`
}
type APPVAL struct {
	Name string `json:"name"`
	Path string `json:"path"`
	Desc string `json:"desc"`
}
func getApps(shell *powershell.Shell) (apps []*APPVAL){
	stdout, err := shell.Exec("$WScript = New-Object -ComObject WScript.Shell")
	//stdout, stderr, err := shell.Execute("$WScript = New-Object -ComObject WScript.Shell")
	if err != nil {
		panic(err)
	}
	stdout, err = shell.Exec("Get-ChildItem -Path \"C:\\ProgramData\\Microsoft\\Windows\\Start Menu\\Programs\\**\\*.lnk\" | ForEach-Object {$WScript.CreateShortcut($_.FullName).TargetPath} | select-string -pattern exe | sort-object -unique")
	//stdout, stderr, err = shell.Execute("Get-ChildItem -Path \"C:\\ProgramData\\Microsoft\\Windows\\Start Menu\\Programs\\**\\*.lnk\" | ForEach-Object {$WScript.CreateShortcut($_.FullName).TargetPath} | select-string -pattern exe | sort-object -unique")
	if err != nil {
		panic(err)
	}

	//fmt.Println(stdout)
	//fmt.Println(stderr)
	applist := strings.Split(stdout, "\r\n")
	//apps := make([]APPVAL, len(applist))
	for _, app := range applist {

		if _, err := os.Stat(app); err == nil {
			cmd := fmt.Sprintf("Get-ItemPropertyValue \"%v\" -Name versionInfo | format-list", app)
			//fmt.Println(cmd)
			stdout, err = shell.Exec(cmd)
			//stdout, stderr, err = shell.Execute(cmd)
			if err != nil {
				panic(err)
			}
			appdetails := strings.Split(strings.TrimSpace(stdout), "\r\n")
			appdetaild := map[string]string{}
			for _, appdetail := range appdetails {
				kv := strings.Split(appdetail, ":")
				key := strings.TrimSpace(kv[0])
				if len(kv) >= 2 {
					val := strings.TrimSpace(kv[1])
					appdetaild[key] = val
				} else {
					appdetaild[key] = ""
				}
			}
			desc := app
			name := app
			val, exests := appdetaild["FileDescription"]
			if exests {
				desc = val
			}

			val, exests = appdetaild["ProductName"]
			if exests {
				name = val
			}

			fmt.Println(stdout)
			a := &APPVAL{
				Name: name,
				Path: app,
				Desc: desc,
			}
			apps = append(apps, a)

		} else if os.IsNotExist(err) {
			fmt.Println("notExist")
			continue
		} else {
			fmt.Println("other")
			continue
		}
	}
	return apps
}
/*

type PSCMD struct {
	CMD string `uri:"cmd" binding:"required"`
	ARG string `uri:"arg" binding:"required"`
}
type APPVAL struct {
	Name string `json:"name"`
	Path string `json:"path"`
	Desc string `json:"desc"`
}
var pscmd PSCMD

shell, err := powershell.New()
if err != nil {
panic(err)
}
defer shell.Exit()

 */