package libs

import(
  
)

type Config struct{
  Timing            Timing            `toml "timing"`
  Settings          Settings          `toml "settings"`
  Exclusions        Exclusions        `toml "exclusions"`
  Inclusions        Inclusions        `toml "inclusions"`
  Advanced_Settings Advanced_Settings `toml "advanced_settings"`
}

type Timing struct{
  TimeOfDay string `toml "timeOfDay"`
  Type      string    `toml "type"`
  Frequency int       `toml "times"`
  Dates     []string  `toml "dates"`
}

type Settings struct{
  Full_Backup_Frequency    int    `toml "full_Backup_Frequency"`
  Partial_Backup_Frequency int    `toml "partial_Backup_Frequency"`
  Use_Exclusions           bool   `toml "Use_Exclusions "`
  Use_Inclusions           bool   `toml "Use_Inclusions"`
  Email_Extension          string `toml "Email_Extension"`
  Network_Path             string `toml "Network_Path"`
}

type Exclusions struct{
  Default_Exclusions   []string  `toml "Default_Exclusions"`
  General_Exclusions   []string  `toml "General_Exclusions"`
  File_Type_Exclusions []string  `toml "File_Type_Exclusions"`
  Profile_Exclusions   []string  `toml "Profile_Exclusions"`
}

type Inclusions struct{
  Default_Inclusions []string `toml "Default_Inclusions"`
  General_Inclusions []string `toml "General_Inclusions"`
  Profile_Inclusions []string `toml "Profile_Inclusions"`
}

type Advanced_Settings struct{
  Use_Ecryption bool   `toml "Use_Ecryption"`
  Domain        string `toml "Domain"`
}
