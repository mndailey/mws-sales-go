package util

import (
	"os"
	"os/exec"
	"os/signal"
	"sync"
	"time"
)

// ExitInfo structured used for exit processing
type ExitInfo struct {
	mutex           sync.Mutex
	init            bool
	reboot          bool
	sigs            chan os.Signal
	done            chan bool
	doneClosed      bool
	onExitListeners []func()
}

var exitInfo ExitInfo

// exitInit - initialize processing for exit handling
//==============================================================================
func exitInit() {
	exitInfo.mutex.Lock()
	defer exitInfo.mutex.Unlock()
	if !exitInfo.init {
		exitInfo.init = true
		exitInfo.sigs = make(chan os.Signal, 1)
		exitInfo.done = make(chan bool)
		signal.Notify(exitInfo.sigs, os.Interrupt)

		go func() {
			<-exitInfo.sigs
			reboot := exitInfo.reboot
			exitInfo.mutex.Lock()
			fl := make([]func(), len(exitInfo.onExitListeners), len(exitInfo.onExitListeners))
			for idx := range fl {
				fl[idx] = exitInfo.onExitListeners[idx]
			}
			if !exitInfo.doneClosed {
				exitInfo.doneClosed = true
				close(exitInfo.done)
			}
			exitInfo.mutex.Unlock()
			for idx := range fl {
				if fl[idx] != nil {
					fl[idx]()
				}
			}
			if reboot {
				exec.Command("sudo", "reboot", "now").Run()
			}
			os.Exit(0)
		}()

	}
}

// CallOnExit - Add a callback on exit
//==============================================================================
func CallOnExit(callOnExit func()) {
	if callOnExit != nil {
		exitInit()
		exitInfo.mutex.Lock()
		defer exitInfo.mutex.Unlock()
		exitInfo.onExitListeners = append(exitInfo.onExitListeners, callOnExit)
	}
}

// Exit - Call this to cleanly exit program
//==============================================================================
func Exit(reboot bool) {
	exitInit()
	exitInfo.reboot = reboot
	exitInfo.sigs <- os.Interrupt
	time.Sleep(time.Hour)
}

// WaitForExit - block until program exits
//==============================================================================
func WaitForExit() {
	<-exitInfo.done
}
