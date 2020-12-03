package libs

import(
  "strings"
)

type Config struct{
  Timing            Timing            `toml "timing"`
  Settings          Settings          `toml "settings"`
  Exclusions        Exclusions        `toml "exclusions"`
  Inclusions        Inclusions        `toml "inclusions"`
  Advanced_Settings Advanced_Settings `toml "advanced_settings"`
}

type Timing struct{
  TimeOfDay TimeOfDay `toml "timeOfDay"`
  Type      Type      `toml "type"`
  Frequency Frequency `toml "times"`
  Dates     Dates     `toml "dates"`
}

type Settings struct{
  Full_Backup_Frequency    full_Backup_Frequency    `toml "full_Backup_Frequency"`
  Partial_Backup_Frequency partial_Backup_Frequency `toml "partial_Backup_Frequency"`
  Use_Exclusions           Use_Exclusions           `toml "Use_Exclusions "`
  Use_Inclusions           Use_Inclusions           `toml "Use_Inclusions"`
  Email_Extension          Email_Extension          `toml "Email_Extension"`
  Network_Path             Network_Path             `toml "Network_Path"`
}

type Exlcusions struct{
  Default_Exclusions   Default_Exclusions   `toml "Default_Exclusions"`
  General_Exclusions   General_Exclusions   `toml "General_Exclusions"`
  File_Type_Exclusions File_Type_Exclusions `toml "File_Type_Exclusions"`
  Profile_Exclusions   Profile_Exclusions   `toml "Profile_Exclusions"`
}

type Inclusions struct{
  Default_Inclusions Default_Inclusions `toml "Default_Inclusions"`
  General_Inclusions General_Inclusions `toml "General_Inclusions"`
  Profile_Inclusions Profile_Inclusions `toml "Profile_Inclusions"`
}

type Advanced_Settings struct{
  Use_Ecryption Use_Ecryption `toml "Use_Ecryption"`
  Domain        Domain        `toml "Domain"`
}
