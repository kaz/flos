package audit

/*
	An auditor for dawin actually does nothing.
*/

type (
	Auditor struct {
		Event chan *Event
	}
	Event struct {
		Acts        []string
		FileName    string
		ProcessInfo string
	}
)

func NewAuditor(ignore string, perm bool) (*Auditor, error) {
	return &Auditor{make(chan *Event)}, nil
}

func (a *Auditor) watch(path string, addFlag int) error {
	return nil
}
func (a *Auditor) WatchFile(path string) error {
	return nil
}
func (a *Auditor) WatchMount(path string) error {
	return nil
}

func (a *Auditor) startAudit() {
	<-a.Event
}
