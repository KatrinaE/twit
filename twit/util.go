package twit

import(
	"encoding/json"
)

func writeJson (??, ??) {
	b, err := json.Marshal(tweet)
	if err != nil {
		log.Fatal(err)
	}
	w.Write(b)
}
