package script

var Registry = map[string]map[string]Func{
	"ubuntu": {
		"install-docker": New(UbuntuInstallDocker()),
	},
}
