package logging

type UserData struct {
	Username string
	Files    int64
	Folders  int64
	Size     int64
	Src      BreakDown
	Dst      BreakDown
}

type BreakDown struct {
	Desk int64
	Docs int64
	Down int64
	Pics int64
}
