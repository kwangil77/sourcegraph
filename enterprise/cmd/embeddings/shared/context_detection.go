package shared

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/sourcegraph/sourcegraph/internal/lazyregexp"
)

var NO_CONTEXT_MESSAGES_REGEXPS = []*lazyregexp.Regexp{
	// Common greetings
	lazyregexp.New(`^(hello|hey|hi|what['’]s up|how's it going)( Cody)?[!\.\?]?$`),

	// Clear reference to previous message
	lazyregexp.New(`(previous|above)\s+(message|code|text)`),
	lazyregexp.New(
		`(translate|convert|change|for|make|refactor|rewrite|ignore|describe|explain|fix|try|show)\s+(that|this|above|previous|it|again)`,
	),
	lazyregexp.New(`i don['’]t understand`),
	lazyregexp.New(`what you just said`),
	lazyregexp.New(`(explain|describe).*in more detail`),

	// Correcting previous message
	lazyregexp.New(
		`(this|that).*?\s+(is|seems|looks)\s+(wrong|incorrect|bad|good)`,
	),
	lazyregexp.New(
		`(this|that).*?\s+(does not|doesn't work)`,
	),
	lazyregexp.New(`(is not|isn['’]t) (correct|right)`),
	lazyregexp.New(`i don['’]t think that['’]s (correct|right)`),
	lazyregexp.New(`(does not|doesn['’]t) (look|seem) (correct|right)`),
	lazyregexp.New(`are you (sure|certain)`),
	lazyregexp.New(`you're (incorrect|not right|wrong)`),

	// Clearly moving on to new topic
	lazyregexp.New(`^(yes|no|correct|wrong|nope|yep|now|cool)(\s|.|,|!)`),

	// User provided their own code context in the form of a Markdown code block.
	lazyregexp.New("```"),
}

func isContextRequiredForChatQuery(query string) bool {
	queryTrimmed := strings.TrimSpace(query)
	payload := map[string]string{
		"text": queryTrimmed,
	}
	jsonPayload, _ := json.Marshal(payload)

	resp, _ := http.Post("http://127.0.0.1:8000/context-probabilities", "application/json", bytes.NewBuffer(jsonPayload))
	probabilities := map[string]float64{}
	json.NewDecoder(resp.Body).Decode(&probabilities)
	return probabilities["no_context_probability"] < 0.5
}
