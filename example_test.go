package efp

import (
	"io/ioutil"
	"testing"
)

func TestLexisFile(t *testing.T) {
	bytes, _ := ioutil.ReadFile("test_files/lexis.efp")
	p, errs := PrototypeBytes(bytes)
	assert(t, errs == nil, "errs must be nil")
	assertNow(t, p != nil, "prototype must not be nil")
	e, errs := p.ValidateString(`
        name = "Example"
        publisher = "LexisTextUs"
        topic = "Immigration"

        question {
            ask = "What's your name"?
            regex = "[a-zA-Z]+"
        }
        question {
            ask = "When was the last time you entered Australia? (dd/mm/yyyy)?"
            regex = "(^(((0[1-9]|1[0-9]|2[0-8])[\/](0[1-9]|1[012]))|((29|30|31)[\/](0[13578]|1[02]))|((29|30)[\/](0[4,6,9]|11)))[\/](19|[2-9][0-9])\d\d$)|(^29[\/]02[\/](19|[2-9][0-9])(00|04|08|12|16|20|24|28|32|36|40|44|48|52|56|60|64|68|72|76|80|84|88|92|96)$)"
        }

        question {
            ask = "How many days will you be staying in Australia?"
            regex = "^(\d+)(days)?$"
        }

        `)
	assertNow(t, errs == nil, "errs must be nil")
	assert(t, e.Field("name", 0).Value(0) == "Example", "failed name test")
	assert(t, e.Field("publisher", 0).Value(0) == "LexisTextUs", "failed publisher test")
	assert(t, e.Field("topic", 0).Value(0) == "Immigration", "failed topic test")
	assert(t, len(e.Elements("question")) == 3, "wrong question length")

}
