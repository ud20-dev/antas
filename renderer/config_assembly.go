//go:build !natif

package renderer

import (
    "time"

	"github.com/klippa-app/go-pdfium"
	"github.com/klippa-app/go-pdfium/webassembly"
)

// Be sure to close pools/instances when you're done with them.
var pool pdfium.Pool
var instance pdfium.Pdfium
const DPI = 200 // The DPI to render the page in. You can change this to your needs.

func Init() (error) {
	// Init the PDFium library and return the instance to open documents.
	// You can tweak these configs to your need. Be aware that workers can use quite some memory.
    var err error
	pool, err = webassembly.Init(webassembly.Config{
		MinIdle:  1, // Makes sure that at least x workers are always available
		MaxIdle:  1, // Makes sure that at most x workers are ever available
		MaxTotal: 1, // Maxium amount of workers in total, allows the amount of workers to grow when needed, items between total max and idle max are automatically cleaned up, while idle workers are kept alive so they can be used directly.
	})
    if err != nil {
        return err
    }

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