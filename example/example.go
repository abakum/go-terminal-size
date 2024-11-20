package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"

	"golang.org/x/sys/windows"
)

func main() {
	con, err := windows.GetStdHandle(windows.STD_INPUT_HANDLE)
	if err != nil {
		log.Fatalf("get stdin handle: %s", err)
	}

	var originalConsoleMode uint32

	err = windows.GetConsoleMode(con, &originalConsoleMode)
	if err != nil {
		log.Fatalf("get console mode: %s", err)
	}

	newConsoleMode := uint32(windows.ENABLE_MOUSE_INPUT) |
		windows.ENABLE_WINDOW_INPUT |
		windows.ENABLE_PROCESSED_INPUT |
		windows.ENABLE_EXTENDED_FLAGS

	err = windows.SetConsoleMode(con, newConsoleMode)
	if err != nil {
		log.Fatalf("set console mode: %s", err)
	}

	defer func() {
		resetErr := windows.SetConsoleMode(con, originalConsoleMode)
		if err == nil && resetErr != nil {
			log.Fatalf("reset console mode: %s", resetErr)
		}
	}()

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	var read uint32
	events := [16]windows.InputRecord{}
	for {
		if ctx.Err() != nil {
			break
		}

		if err := windows.ReadConsoleInput(con, &events[0], uint32(len(events)), &read); err != nil {
			log.Fatalf("read input events: %s", err)
		}

		fmt.Printf("Read %d events:\n", len(events))
		for _, event := range events[:read] {
			var e any
			switch event.EventType {
			case windows.KEY_EVENT:
				e = event.KeyEvent()
			case windows.MOUSE_EVENT:
				e = event.MouseEvent()
			case windows.WINDOW_BUFFER_SIZE_EVENT:
				e = event.WindowBufferSizeEvent()
			case windows.FOCUS_EVENT:
				e = event.FocusEvent()
			case windows.MENU_EVENT:
				e = event.MenuEvent()
			}
			fmt.Printf("%#v\n", e)
		}
	}
}
