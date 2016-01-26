package config

import (
	"log"
	"os"
	"text/template"
)

const instructions = `
The first time you start {{.SysName}} you will need to supply a few details
about your host environment, the administrator and the git repository hosting
environment. Currently, we only support GitHub for hosting git repositories.

{{.SysName}} can either read a configuration file with the necessary
information (see the example below), or you can provide these details as
command line arguments (also shown below).

Here is an example {{.ConfigFileName}} file:

{
  "HomepageURL": "http://example.com/",
  "ClientID": "123456789",
  "ClientSecret": "123456789abcdef",
  "BasePath": "/usr/share/{{.SysNameLC}}/"
}

Before you can start you will need to register the {{.SysName}} application
at GitHub; you will need to do this from the administrator account.

1. Go to https://github.com/settings/applications/new
2. Enter the information requested.
   - Application name: e.g. "{{.SysName}} at University of Stavanger"
   - Homepage URL: e.g. "http://{{.SysNameLC}}.ux.uis.no"
   - Authorization callback URL: e.g. "http://{{.SysNameLC}}.ux.uis.no/oauth"

Note that, the Homepage URL must be a fully qualified URL, including http://.
This must be the hostname (or an alias) of server running the '{{.SysNameLC}}'
program. This server must have a public IP address, since GitHub will make calls
to this server to support {{.SysName}}'s functionality. Further, {{.SysName}}
requires that the Authorization callback URL is the same as the Homepage URL
with the added "/oauth" path.

Once you have completed the above steps, the Client ID and Client Secret will be
available from the GitHub web interface. Simply copy each of these OAuth tokens
and paste them into the configuration file, or on the command line when starting
{{.SysName}} for the first time. You will not need to repeat this process
when starting {{.SysName}} in the future.

If you need to obtain the OAuth tokens at a later time, e.g. if you have deleted
the configuration file, go to: https://github.com/settings/developers and
select your Application to be able to view the OAuth tokens again.

`

// PrintInstructions prints the configuration instructions.
func PrintInstructions() {
	data := struct {
		SysName, SysNameLC, ConfigFileName string
	}{
		SysName, SysNameLC, FileName,
	}
	t := template.Must(template.New("instructions").Parse(instructions))
	err := t.Execute(os.Stdout, data)
	if err != nil {
		log.Fatalln(err)
	}
}
