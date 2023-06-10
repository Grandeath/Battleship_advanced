package timer

import (
	"context"
	"fmt"
	"time"

	gui "github.com/grupawp/warships-gui/v2"
)

type Timer struct {
	ClockField *gui.Text
	Time       int
	TimeGoes   bool
}

func NewTimer() Timer {
	clockField := gui.NewText(100, 2, "Time to move:", nil)
	return Timer{ClockField: clockField}
}

func (t *Timer) setClock() {
	timeText := fmt.Sprintf("Time to move: %d", t.Time)
	t.ClockField.SetText(timeText)
}

func (t *Timer) StartClock(ctx context.Context) {
	t.TimeGoes = false
	go t.Ticking(ctx)
}

func (t *Timer) Ticking(ctx context.Context) {
mainloop:
	for {
		if t.TimeGoes {
			select {
			case <-ctx.Done():
				break mainloop
			default:
				t.Time--
				t.setClock()
			}
		} else {
			select {
			case <-ctx.Done():
				break mainloop
			default:
			}
		}
		time.Sleep(time.Millisecond * 1000)
	}
}

func (t *Timer) TimeStop() {
	t.TimeGoes = false
	t.ClockField.SetText("Enemy turn")
}
