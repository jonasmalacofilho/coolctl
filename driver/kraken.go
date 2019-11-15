package driver

import (
	"context"
	"flag"
	"fmt"
	"image/color"
	"log"
	"strings"

	"github.com/google/gousb"
)

const (
	productID = 0x170e // Kraken X (X42, X52, X62 or X72)
	vendorID  = 0x1e71 // NZXT

	readEndpoint = 1
	readLength   = 64

	writeEndpoint = 1
	writeLength   = 65

	totalLEDs = 9
)

var (
	config    = flag.Int("config", 1, "Configuration number to use with the device.")
	iface     = flag.Int("interface", 0, "Interface to use on the device.")
	alternate = flag.Int("alternate", 0, "Alternate setting to use on the interface.")
	debug     = flag.Int("debug", 0, "Debug level for libusb.")
	bufSize   = flag.Int("buffer_size", 0, "Number of buffer transfers, for data prefetching.")
	timeout   = flag.Duration("timeout", 0, "Timeout for the command. 0 means infinite.")

	speedChannels = map[string]int{
		"fan":  0x80, // 25, 100
		"pump": 0xc0, // 50, 100
	}

	colorChannels = map[string]int{
		"sync": 0x0,
		"logo": 0x1,
		"ring": 0x2,
	}

	colorModes = map[string][]int{
		// byte3/mode, byte2/reverse, byte4/modifier, min colors, max colors, only ring (0=no, 1=yes)
		"off":                          []int{0x00, 0x00, 0x00, 0, 0, 0},
		"fixed":                        []int{0x00, 0x00, 0x00, 1, 1, 0},
		"super-fixed":                  []int{0x00, 0x00, 0x00, 1, 9, 0}, // independent logo + ring leds
		"fading":                       []int{0x01, 0x00, 0x00, 2, 8, 0},
		"spectrum-wave":                []int{0x02, 0x00, 0x00, 0, 0, 0},
		"backwards-spectrum-wave":      []int{0x02, 0x10, 0x00, 0, 0, 0},
		"marquee-3":                    []int{0x03, 0x00, 0x00, 1, 1, 1},
		"marquee-4":                    []int{0x03, 0x00, 0x08, 1, 1, 1},
		"marquee-5":                    []int{0x03, 0x00, 0x10, 1, 1, 1},
		"marquee-6":                    []int{0x03, 0x00, 0x18, 1, 1, 1},
		"backwards-marquee-3":          []int{0x03, 0x10, 0x00, 1, 1, 1},
		"backwards-marquee-4":          []int{0x03, 0x10, 0x08, 1, 1, 1},
		"backwards-marquee-5":          []int{0x03, 0x10, 0x10, 1, 1, 1},
		"backwards-marquee-6":          []int{0x03, 0x10, 0x18, 1, 1, 1},
		"covering-marquee":             []int{0x04, 0x00, 0x00, 1, 8, 1},
		"covering-backwards-marquee":   []int{0x04, 0x10, 0x00, 1, 8, 1},
		"alternating":                  []int{0x05, 0x00, 0x00, 2, 2, 1},
		"moving-alternating":           []int{0x05, 0x08, 0x00, 2, 2, 1},
		"backwards-moving-alternating": []int{0x05, 0x18, 0x00, 2, 2, 1},
		"breathing":                    []int{0x06, 0x00, 0x00, 1, 8, 0}, // colors for each step
		"super-breathing":              []int{0x06, 0x00, 0x00, 1, 9, 0}, // one step, independent logo + ring leds
		"pulse":                        []int{0x07, 0x00, 0x00, 1, 8, 0},
		"tai-chi":                      []int{0x08, 0x00, 0x00, 2, 2, 1},
		"water-cooler":                 []int{0x09, 0x00, 0x00, 0, 0, 1},
		"loading":                      []int{0x0a, 0x00, 0x00, 1, 1, 1},
		"wings":                        []int{0x0c, 0x00, 0x00, 1, 1, 1},
		"super-wave":                   []int{0x0d, 0x00, 0x00, 1, 8, 1}, // independent ring leds
		"backwards-super-wave":         []int{0x0d, 0x10, 0x00, 1, 8, 1}, // independent ring leds
	}

	animationSpeeds = map[string]int{
		"slowest": 0x0,
		"slower":  0x1,
		"normal":  0x2,
		"faster":  0x3,
		"fastest": 0x4,
	}
)

type contextReader interface {
	ReadContext(context.Context, []byte) (int, error)
}

// KrakenDriver holds all driver relevant informations
type KrakenDriver struct {
	ProductID gousb.ID
	VendorID  gousb.ID
	*gousb.Context
	*gousb.Interface
	*gousb.InEndpoint
	*gousb.OutEndpoint
}

// NewKrakenDriver creates a new USB Context instance & returns a new KrakenDriver
func NewKrakenDriver() *KrakenDriver {
	flag.Parse()
	ctx := gousb.NewContext()
	ctx.Debug(*debug)

	return &KrakenDriver{
		ProductID: productID,
		VendorID:  vendorID,
		Context:   ctx,
	}
}

// Connect connects to the USB device
func (d *KrakenDriver) Connect() {
	dev, err := d.Context.OpenDeviceWithVIDPID(d.VendorID, d.ProductID)
	if err != nil {
		log.Fatal("NXZT Kraken X (X42, X52, X62 or X72) not found")
	}
	defer dev.Close()

	dev.SetAutoDetach(true)

	cfg, err := dev.Config(*config)
	if err != nil {
		log.Fatalf("dev.Config(%d): %v", *config, err)
	}

	d.Interface, err = cfg.Interface(*iface, *alternate)
	if err != nil {
		log.Fatalf("cfg.Interface(%d, %d): %v", *iface, *alternate, err)
	}

	d.InEndpoint, err = d.Interface.InEndpoint(readEndpoint)
	if err != nil {
		log.Fatalf("dev.InEndpoint(): %s", err)
	}

	d.OutEndpoint, err = d.Interface.OutEndpoint(writeEndpoint)
	if err != nil {
		log.Fatalf("dev.OutEndpoint(): %s", err)
	}
}

// Read reads from the USB device
func (d *KrakenDriver) Read() []byte {
	var rdr contextReader = d.InEndpoint
	if *bufSize > 1 {
		log.Print("creating buffer...")
		s, err := d.InEndpoint.NewStream(readLength, *bufSize)
		if err != nil {
			log.Fatalf("ep.NewStream(): %v", err)
		}
		defer s.Close()
		rdr = s
	}

	opCtx := context.Background()
	if *timeout > 0 {
		var done func()
		opCtx, done = context.WithTimeout(opCtx, *timeout)
		defer done()
	}
	msg := make([]byte, readLength)
	_, err := rdr.ReadContext(opCtx, msg)
	if err != nil {
		log.Fatalf("reading from device failed: %v", err)
	}

	return msg
}

// Write writes to the USB device
func (d *KrakenDriver) Write(data []byte) {
	padding := make([]byte, writeLength-len(data))
	data = append(data, padding...)
	_, err := d.OutEndpoint.Write(data)
	if err != nil {
		log.Fatalf("could not write data %d to device", data)
	}
}

// GetStatus returns the current device status
func (d *KrakenDriver) GetStatus() {
	msg := d.Read()

	temperature := fmt.Sprintf("%d.%d", uint64(msg[1]), uint64(msg[2]))
	fanSpeed := uint64(msg[3])<<8 | uint64(msg[4])
	pumpSpeed := uint64(msg[5])<<8 | uint64(msg[6])
	firmwareVersion := fmt.Sprintf("%d.%d.%d", uint64(msg[0xb]), uint64(msg[0xc])<<8|uint64(msg[0xd]), uint64(msg[0xe]))

	fmt.Println("============================================")
	fmt.Println(fmt.Sprintf("  Liquid temperature %s °C", temperature))
	fmt.Println(fmt.Sprintf("  Fan speed %d rpm", fanSpeed))
	fmt.Println(fmt.Sprintf("  Pump speed %d rpm", pumpSpeed))
	fmt.Println(fmt.Sprintf("  Firmware Version: %s", firmwareVersion))
	fmt.Println("============================================")
}

// SetColor sets the color of a channel & mode
func (d *KrakenDriver) SetColor(channel, mode string, colors []string) {
	colorChannel, ok := colorChannels[channel]
	if !ok {
		log.Fatalf("channel %s not found", channel)
	}

	colorMode, ok := colorModes[mode]
	if !ok {
		log.Fatalf("mode %s not found", mode)
	}

	mval, mod2, mod4, mincolors, maxcolors, ringonly := colorMode[0], colorMode[1], colorMode[2], colorMode[3], colorMode[4], colorMode[5]
	if ringonly == 1 && channel != "ring" {
		log.Fatalf("mode %s unsupported with channel %s", mode, channel)
	}

	steps := generateSteps(paletteFromColors(colors), mincolors, maxcolors, mode, ringonly)

	for seq, step := range steps {
		logoRed, logoGreen, logoBlue, _ := step[0].RGBA()

		var buf []byte
		buf = append(buf, 0x2)
		buf = append(buf, 0x4c)
		buf = append(buf, byte(mod2|colorChannel))
		buf = append(buf, byte(mval))
		buf = append(buf, byte(animationSpeeds["normal"]|seq<<5|mod4))
		buf = append(buf, byte(logoGreen))
		buf = append(buf, byte(logoRed))
		buf = append(buf, byte(logoBlue))
		for _, leds := range step[1:] {
			red, green, blue, _ := leds.RGBA()
			buf = append(buf, byte(red))
			buf = append(buf, byte(green))
			buf = append(buf, byte(blue))
		}

		d.Write(buf)
	}
}

func generateSteps(colors color.Palette, mincolors, maxcolors int, mode string, ringonly int) []color.Palette {
	if len(colors) < mincolors {
		log.Fatalf("not enough colors for mode %s, at least %d required", mode, mincolors)
	} else if maxcolors == 0 {
		if len(colors) > 0 {
			log.Printf("too many colors for mode %s, none needed", mode)
			colors = color.Palette{color.RGBA{0, 0, 0, 1}} // discard the input but ensure at least one step
		}
	} else if len(colors) > maxcolors {
		log.Printf("too many colors for mode %s, dropping to %d", mode, maxcolors)
		colors = colors[:maxcolors]
	}

	if len(colors) == 0 {
		colors = color.Palette{color.RGBA{0, 0, 0, 1}}
	}

	var steps []color.Palette

	if !strings.Contains(mode, "super") {
		for colorNum := range colors {
			var colorPalette color.Palette
			for i := 0; i < totalLEDs; i++ {
				colorPalette = append(colorPalette, colors[colorNum])
			}
			steps = append(steps, colorPalette)
		}
	} else if ringonly == 1 {
		steps = append(steps, color.Palette{color.RGBA{0, 0, 0, 1}})
		steps = append(steps, colors)
	} else {
		steps = append(steps, colors)
	}

	return steps
}
