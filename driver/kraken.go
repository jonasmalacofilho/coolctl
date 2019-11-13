package driver

import (
	"context"
	"flag"
	"fmt"
	"log"

	"github.com/google/gousb"
)

const (
	productID = 0x170e // Kraken X (X42, X52, X62 or X72)
	vendorID = 0x1e71 // NZXT
)

var (
	config    = flag.Int("config", 1, "Configuration number to use with the device.")
	iface     = flag.Int("interface", 0, "Interface to use on the device.")
	alternate = flag.Int("alternate", 0, "Alternate setting to use on the interface.")
	endpoint  = flag.Int("endpoint", 1, "Endpoint number to which to connect (without the leading 0x8).")
	debug     = flag.Int("debug", 0, "Debug level for libusb.")
	size      = flag.Int("read_size", 64, "Number of bytes of data to read in a single transaction.")
	bufSize   = flag.Int("buffer_size", 0, "Number of buffer transfers, for data prefetching.")
	timeout   = flag.Duration("timeout", 0, "Timeout for the command. 0 means infinite.")
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
		log.Fatalf("could not open a device: %v", err)
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

	d.InEndpoint, err = d.Interface.InEndpoint(*endpoint)
	if err != nil {
		log.Fatalf("dev.InEndpoint(): %s", err)
	}
}

// Disconnect closes the USB Context instance
func (d *KrakenDriver) Disconnect() {
	d.Context.Close()
}

func (d *KrakenDriver) Read() []byte {
	var rdr contextReader = d.InEndpoint
	if *bufSize > 1 {
		log.Print("Creating buffer...")
		s, err := d.InEndpoint.NewStream(*size, *bufSize)
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
	msg := make([]byte, *size)
	_, err := rdr.ReadContext(opCtx, msg)
	if err != nil {
		log.Fatalf("Reading from device failed: %v", err)
	}

	return msg
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
