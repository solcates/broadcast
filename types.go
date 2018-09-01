package broadcast

//Discoverer defines the API contract for discovering devices.
type Discoverer interface {
	Discover() (found []string, err error)
	SetFindself(setting bool)
	SetTimeout(setting int)
	//SetPayload(payload string)
	//SetPort(port int)
}
