package main

import(
  "golang.org/x/sys/windows/registry"
  "fmt"
  "log"
  "os"
  "path/filepath"
)

func GetRegistryInt(loc, key string)string{
  k, err := registry.OpenKey(registry.CURRENT_USER, loc, registry.QUERY_VALUE)
  if err != nil {
    fmt.Println("ERROR HERE",key)
	  log.Fatal(err)
   }
   defer k.Close()

   s, _, err := k.GetIntegerValue(key)
   if err != nil {
     fmt.Println("NEW ERROR",key)
	   log.Fatal(err)
    }
    // formats as KEY : HEX
    return fmt.Sprintf(`%s = "%x"`,key, s)
}

func GetRegistryStr(loc, key string)string{
  k, err := registry.OpenKey(registry.CURRENT_USER, loc, registry.QUERY_VALUE)
  if err != nil {
    fmt.Println("ERROR HERE",key)
	  log.Fatal(err)
   }
   defer k.Close()

   s, _, err := k.GetStringValue(key)
   if err != nil {
     fmt.Println("NEW ERROR",key)
	   log.Fatal(err)
    }
    return fmt.Sprintf(`%s = "%s"`,key, s)
}

func GetRegistryBin(loc, key string)string{
  k, err := registry.OpenKey(registry.CURRENT_USER, loc, registry.QUERY_VALUE)
  if err != nil {
    fmt.Println("ERROR HERE",key)
	  log.Fatal(err)
   }
   defer k.Close()

   s, _, err := k.GetBinaryValue(key)
   if err != nil {
     fmt.Println("NEW ERROR",key)
	   log.Fatal(err)
    }
    return fmt.Sprintf(`%s = "%x"`,key, s)
}

func main(){
  // Gets all color settings for the user
  AccentColorMenu := GetRegistryInt(`SOFTWARE\Microsoft\Windows\CurrentVersion\Explorer\Accent`,"AccentColorMenu")
  StartColorMenu := GetRegistryInt(`SOFTWARE\Microsoft\Windows\CurrentVersion\Explorer\Accent`,"StartColorMenu")
  AccentPalette := GetRegistryBin(`SOFTWARE\Microsoft\Windows\CurrentVersion\Explorer\Accent`,"AccentPalette")
  // Select theme Settings
  SystemUsesLightTheme := GetRegistryInt(`SOFTWARE\Microsoft\Windows\CurrentVersion\Themes\Personalize`,"SystemUsesLightTheme")
  AppsUseLightTheme := GetRegistryInt(`SOFTWARE\Microsoft\Windows\CurrentVersion\Themes\Personalize`,"AppsUseLightTheme")
  ColorPrevalence := GetRegistryInt(`SOFTWARE\Microsoft\Windows\CurrentVersion\Themes\Personalize`,"ColorPrevalence")
  EnableTransparency := GetRegistryInt(`SOFTWARE\Microsoft\Windows\CurrentVersion\Themes\Personalize`,"EnableTransparency")
  //Bulk of theme Settings
  ColorSetFromTheme := GetRegistryInt(`SOFTWARE\Microsoft\Windows\CurrentVersion\Themes`,"ColorSetFromTheme")
  CurrentTheme := GetRegistryStr(`SOFTWARE\Microsoft\Windows\CurrentVersion\Themes`,"CurrentTheme")
  InstallVisualStyleColor := GetRegistryStr(`SOFTWARE\Microsoft\Windows\CurrentVersion\Themes`,"InstallVisualStyleColor")
  InstallVisualStyleSize := GetRegistryStr(`SOFTWARE\Microsoft\Windows\CurrentVersion\Themes`,"InstallVisualStyleSize")
  LastHighContrastTheme := GetRegistryStr(`SOFTWARE\Microsoft\Windows\CurrentVersion\Themes`,"LastHighContrastTheme")
  SetupVersion := GetRegistryStr(`SOFTWARE\Microsoft\Windows\CurrentVersion\Themes`,"SetupVersion")
  ThemeChangesDesktopIcons := GetRegistryInt(`SOFTWARE\Microsoft\Windows\CurrentVersion\Themes`,"ThemeChangesDesktopIcons")
  ThemeChangesMousePointers := GetRegistryInt(`SOFTWARE\Microsoft\Windows\CurrentVersion\Themes`,"ThemeChangesMousePointers")
  WallpaperSetFromTheme := GetRegistryInt(`SOFTWARE\Microsoft\Windows\CurrentVersion\Themes`,"WallpaperSetFromTheme")

  creationLoc := `C:\Users\micha\OneDrive\Desktop\`
  user := "dmcanady"
  f,err := os.Create(filepath.Join(creationLoc,fmt.Sprintf(`%s_Preferences.toml`,user)))
  if err != nil{
    fmt.Println("Creation Error:",err)
  }
  //Prints out gathered values
  f.WriteString(`[SOFTWARE\Microsoft\Windows\CurrentVersion\Explorer\Accent]`+"\n")
  f.WriteString(" "+AccentColorMenu+"\n")
  f.WriteString(" "+StartColorMenu+"\n")
  f.WriteString(" "+AccentPalette+"\n")
  f.WriteString(`[SOFTWARE\Microsoft\Windows\CurrentVersion\Themes\Personalize]`+"\n")
  f.WriteString(" "+SystemUsesLightTheme+"\n")
  f.WriteString(" "+AppsUseLightTheme+"\n")
  f.WriteString(" "+ColorPrevalence+"\n")
  f.WriteString(" "+EnableTransparency+"\n")
  f.WriteString(`[SOFTWARE\Microsoft\Windows\CurrentVersion\Themes]`+"\n")
  f.WriteString(" "+ColorSetFromTheme+"\n")
  f.WriteString(" "+CurrentTheme+"\n")
  f.WriteString(" "+InstallVisualStyleColor+"\n")
  f.WriteString(" "+InstallVisualStyleSize+"\n")
  f.WriteString(" "+LastHighContrastTheme+"\n")
  f.WriteString(" "+SetupVersion+"\n")
  f.WriteString(" "+ThemeChangesDesktopIcons+"\n")
  f.WriteString(" "+ThemeChangesMousePointers+"\n")
  f.WriteString(" "+WallpaperSetFromTheme+"\n")
}
