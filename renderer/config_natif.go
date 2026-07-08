//go:build natif

package renderer

import (
    "time"

	"github.com/klippa-app/go-pdfium"
	"github.com/klippa-app/go-pdfium/single_threaded"
)

// Be sure to close pools/instances when you're done with them.
var pool pdfium.Pool
var instance pdfium.Pdfium
const DPI = 200 // The DPI to render the page in. You can change this to your needs.

func Init() (error) {
	// Init the PDFium library and return the instance to open documents.
	// You can tweak these configs to your need. Be aware that workers can use quite some memory.
	pool = single_threaded.Init(single_threaded.Config{})

	var err error
	instance, err = pool.GetInstance(time.Second * 30)
	if err != nil {
		return err
	}
	return nil
}

func Close() error {
	var err error
	err = pool.Close()
	if err != nil {
		return err
	}
	err = instance.Close()
		if err != nil {
		return err
	}
	return nil
}