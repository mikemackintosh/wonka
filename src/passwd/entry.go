package passwd

type Entry struct {
	Password string
	Username string
	UID      int
	GID      int
	Info     string
	HomeDir  string
	Shell    string
	Error    []error
}

func NewEntry(username, password string, uid, gid int, info, homedir, shell string) (Entry, error) {
	e := Entry{
		Username: username,
		Password: "x",
		UID:      uid,
		GID:      gid,
		Info:     info,
		HomeDir:  homedir,
		Shell:    shell,
	}

	return e, nil
}
