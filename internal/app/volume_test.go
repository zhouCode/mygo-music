package app

import (
	"math"
	"testing"

	"github.com/egoist/mygo/ui"
)

// A wheel over the speaker or its slider changes the volume by four
// percent, rather than a platform's scroll distance making it jump by 25.
func TestVolumeWheel(t *testing.T) {
	for _, slider := range []bool{false, true} {
		name := "speaker"
		if slider {
			name = "slider"
		}
		t.Run(name, func(t *testing.T) {
			a := newTestApp()
			tt := ui.NewTester(func(c *ui.Context) {
				if slider {
					a.nativeVolumeSlider(c, popoverColors(false))
				} else {
					a.volumeButton(c, white)
				}
			}, 200, 200)
			tt.Frame()
			r, ok := tt.Find("Volume")
			if !ok {
				t.Fatal("no volume control")
			}
			x, y := r.X+r.W/2, r.Y+r.H/2
			for _, c := range []struct {
				name       string
				from, want float64
				dx, dy     float32
				precise    bool
			}{
				{name: "windows up", from: 0.5, dy: -99, want: 0.54},
				{name: "windows down", from: 0.5, dy: 99, want: 0.46},
				{name: "mac wheel", from: 0.5, dy: -40, want: 0.54},
				{name: "linux wheel", from: 0.5, dy: 50, want: 0.46},
				{name: "upper limit", from: 0.98, dy: -99, want: 1},
				{name: "lower limit", from: 0.02, dy: 99, want: 0},
				{name: "horizontal", from: 0.5, dx: 99, want: 0.5},
				{name: "no movement", from: 0.5, want: 0.5},
				{name: "touchpad up", from: 0.5, dy: -4, precise: true, want: 0.51},
				{name: "touchpad down", from: 0.5, dy: 4, precise: true, want: 0.49},
			} {
				t.Run(c.name, func(t *testing.T) {
					a.volume = c.from
					tt.Frame()
					if c.precise {
						tt.Scroll(x, y, c.dx, c.dy)
					} else {
						tt.Wheel(x, y, c.dx, c.dy)
					}
					if math.Abs(a.volume-c.want) > 1e-9 {
						t.Fatalf("volume = %v, want %v", a.volume, c.want)
					}
				})
			}
		})
	}
}

// Hovering the speaker shows the volume while scrolling, without the
// click it used to need, and still lets the listener use and dismiss it.
func TestVolumePopoverHover(t *testing.T) {
	a := newTestApp()
	a.volume = 0.5
	tt := ui.NewTester(a.view, 1000, 600)
	tt.Frame()
	r, ok := tt.Find("Volume")
	if !ok {
		t.Fatal("no volume button")
	}
	x, y := r.X+r.W/2, r.Y+r.H/2
	if err := tt.Click("Volume"); err != nil {
		t.Fatal(err)
	}
	if !a.volumeOpen {
		t.Fatal("clicking without a preceding hover did not open the popover")
	}
	tt.ClickAt(x, y)
	if a.volumeOpen {
		t.Fatal("a second click left the popover open")
	}
	tt.Move(10, 10)
	tt.Move(x, y)
	if !a.volumeOpen || !tt.HasText("50%") {
		t.Fatal("hover did not show the volume")
	}
	tt.Wheel(x, y, 0, -99)
	if !tt.HasText("54%") {
		t.Fatal("scrolling did not update the shown volume")
	}
	tt.Move(x, r.Y-5) // The gap between the speaker and the slider.
	if !a.volumeOpen {
		t.Fatal("the popover closed on the way to its slider")
	}
	tt.Key(0, ui.KeyUp)
	if math.Abs(a.volume-0.55) > 1e-9 {
		t.Fatalf("the slider did not take the arrow key: %v", a.volume)
	}
	tt.Move(x, y)
	tt.Key(0, ui.KeyEscape)
	tt.Frame()
	if a.volumeOpen {
		t.Fatal("Escape left the popover open under the pointer")
	}
	tt.Move(10, 10)
	tt.Move(x, y)
	if !a.volumeOpen {
		t.Fatal("hover did not open the popover again")
	}
	tt.ClickAt(x, y)
	tt.Frame()
	if a.volumeOpen {
		t.Fatal("clicking the speaker left the popover open")
	}
	tt.ClickAt(x, y)
	if !a.volumeOpen {
		t.Fatal("clicking the speaker did not open the popover")
	}
	tt.ClickAt(10, 10)
	tt.Frame()
	if a.volumeOpen {
		t.Fatal("clicking outside left the popover open")
	}
}
