package nordshare

type Note struct {
	Content       []byte `json:"content"`
	ReadPassword  []byte `json:"readPassword"`
	WritePassword []byte `json:"writePassword"`
	TTL           int64  `json:"ttl"`
}
