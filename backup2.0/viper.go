package main

import (
	"fmt"

	"github.com/spf13/viper"
)

var (
	defaults = map[string]interface{}{
		"settings": map[string]interface{}{
			"Use_Exclusions":  true,
			"Use_Inclusions":  true,
			"Email_Extension": "liberty.edu",
			"Network_Path":    "",
		},
		"tktsystem": map[string]interface{}{
			"Provider": "ServiceNow",
			"URL":      "libertydev.service-now.com/",
		},
		"exclusions": map[string]interface{}{
			"Default_Exclusions": []string{},
			"General_Exclusions": []string{"My Music", "My Pictures", "My Videos", "My Documents", "OneDrive", "Start Menu", "Templates", ".bash_history", "ntuser.pol",
				"Recent", ".atom", ".cisco", ".idlerc", ".osquery", "Application Data", "Cookies", "Local Settings", "NetHood", "PrintHood", "SendTo", "NTUSER.DAT",
				"ntuser.dat.LOG1", "ntuser.dat.LOG2", "AppData", "Appdata", "Dropbox (Liberty University)"},
			"File_Type_Exclusions": []string{".lnk", ".ini", ".dat"},
			"Profile_Exclusions":   []string{},
		},
		"inclusions": map[string]interface{}{
			"Default_Inclusions": []string{},
			"General_Inclusions": []string{"AppData\\Local\\Google\\Chrome\\User Data\\Default\\Bookmarks", "AppData\\Roaming\\Microsoft\\Document Building Blocks",
				"AppData\\Roaming\\Stickies", "Appdata\\Local\\Microsoft\\Outlook", "AppData\\Local\\Microsoft\\Signatures",
				"Appdata\\Local\\Microsoft\\Windows Sidebar"},
			"Profile_Inclusions": []string{},
		},
		"advanced_settings": map[string]interface{}{
			"Use_Ecryption": false,
			"Domain":        "",
		},
	}

	configName = "Settings"
	configPath = []string{
		".",
	}
)

func test() {

	for k, v := range defaults {
		viper.SetDefault(k, v)
	}
	viper.SetConfigName(configName)
	viper.SetConfigType("toml")
	for _, p := range configPath {
		viper.AddConfigPath(p)
	}
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		fmt.Println(err)
		viper.Unmarshal(&conf)
	}
	viper.Unmarshal(&conf)
	viper.SafeWriteConfig()
}
