package eval

import (
	"fmt"
	"github.com/teacoder/gron/manager"
	"os/exec"
	"time"
)

const GRON_WILDCARD = 255

// Execute runs the command passed to it
func execute(cmd string) error {
	cmdExec := exec.Command(cmd)
	op, err := cmdExec.Output()

	if err != nil {
		fmt.Println("Error : " + err.Error())
		return err
	}

	fmt.Print(string(op))

	return nil
}

// Poll will be called periodicaly from the main loop.
// It walks the execution queue and run any commands whose time has come
func poll(entries [manager.MaxEntries]manager.Entry, count int) {
	curTime := time.Now()
	curHour := curTime.Hour()
	curMin := curTime.Minute()
	curDow := int(curTime.Weekday())
	_, _, curDay := curTime.Date()
	curMonth := int(curTime.Month())

	for i := 0; i < count; i++ {
		e := entries[i].Get()
		min := int(e["min"].(uint8))
		hr := int(e["hr"].(uint8))
		cmd := e["cmd"].(string)
		dow := int(e["dow"].(uint8))
		dom := int(e["dom"].(uint8))
		mon := int(e["mon"].(uint8))

		if (mon == curMonth || mon == GRON_WILDCARD) &&
			(dom == curDay || dom == GRON_WILDCARD) &&
			(dow == curDow || dow == GRON_WILDCARD) &&
			(hr == curHour || hr == GRON_WILDCARD) &&
			(min == curMin || min == GRON_WILDCARD) {
			go execute(cmd)
			fmt.Printf("%d : Executed [%s] \n", i, cmd)
		}
	}
}

// Run is the entry function into the evaluator. It runs for ever, polling the
// execution queue every minute
func Run(entries [manager.MaxEntries]manager.Entry, count int) {
	for {
		poll(entries, count)
		time.Sleep(time.Minute)
	}
}
