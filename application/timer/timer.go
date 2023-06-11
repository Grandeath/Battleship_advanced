// timer create a clock which count turn time
package timer

import (
	"context"
	"fmt"
	"time"

	gui "github.com/grupawp/warships-gui/v2"
)

// Timer contain
// ClockField field text,
// Time current time
// TimeGoes if time flowing
type Timer struct {
	ClockField *gui.Text
	Time       int
	TimeGoes   bool
}

// NewTimer create timer with text field in appropriate coordinates
func NewTimer() Timer {
	clockField := gui.NewText(100, 2, "Time to move:", nil)
	return Timer{ClockField: clockField}
}

// setClock update timer text with current time
func (t *Timer) setClock() {
	timeText := fmt.Sprintf("Time to move: %d", t.Time)
	t.ClockField.SetText(timeText)
}

// StartClock starting ticking goroutine
func (t *Timer) StartClock(ctx context.Context) {
	t.TimeGoes = false
	go t.Ticking(ctx)
}

// Ticking take a context which break a loop upon cancelation.
// If TimeGoes is true time is decreasing by one every one second.
// If Not loop is just checking context.
func (t *Timer) Ticking(ctx context.Context) {
mainloop:
	for {
		if t.TimeGoes {
			select {
			case <-ctx.Done():
				break mainloop
			default:
				if t.Time > 0 {
					t.Time--
					t.setClock()
				}
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

// TimeStop stoping a timer by setting TimeGoes to false and changing clock text to "Enemy turn"
func (t *Timer) TimeStop() {
	t.TimeGoes = false
	t.ClockField.SetText("Enemy turn")
}
