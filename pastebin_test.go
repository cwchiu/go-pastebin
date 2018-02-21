package pastebin

import (
	pastebin "github.com/cwchiu/go-pastebin"
	"log"
	"os"
	"testing"
)

func init() {
	f, err := os.OpenFile("trace.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	log.SetOutput(f)
}

func TestExample(t *testing.T) {
	p := pastebin.Pastebin{Key: "dev-key", UserKey: "session-key"}
	// p.Login("guest", "1234")
	// fmt.Println( p.UserKey )

	id, err := p.Put(pastebin.Paste{
		Title:      "paste#3",
		Code:       "Content code",
		MemberOnly: true,
	})
	if err != nil {
		panic(err)
	}
	log.Println(id)

	id2, err := p.Put(*pastebin.CreateNormalPaste("NormalPaste", "NomalPasteContent"))
	if err != nil {
		panic(err)
	}
	log.Println(id2)

	items, err := p.ListByUser(1000)
	if err != nil {
		panic(err)
	}
	log.Println(items)

	// code, err := p.GetByUser(id)
	// if err != nil {
	// panic(err)
	// }
	// log.Println(code)
	// err = p.DelByUser(id)
	// if err != nil {
	// panic(err)
	// }
	// info, err := p.InfoByUser()
	// if err != nil {
	// panic(err)
	// }
	// log.Println(info)
}
