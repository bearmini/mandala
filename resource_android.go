// +build android

package mandala

import (
	"archive/zip"
	"fmt"
	"io"
	"path/filepath"
	"unsafe"
)

// #include <android/native_activity.h>
// #include "resource_android.h"
import "C"

func loadResource(activity unsafe.Pointer, filename string) ([]byte, error) {
	apkPath := C.GoString(C.getPackageName((*C.ANativeActivity)(activity)))

	// Open a zip archive for reading.
	r, err := zip.OpenReader(apkPath)
	if err != nil {
		return nil, err
	}
	defer r.Close()

	// Iterate through the files in the archive.
	for _, f := range r.File {
		if f.Name == filepath.Join("res", filename) {
			rc, err := f.Open()
			if err != nil {
				return nil, err
			}
			buffer := make([]byte, f.UncompressedSize64)
			_, err = io.ReadFull(rc, buffer)
			if err != nil {
				return nil, err
			}
			rc.Close()
			return buffer, nil
		}
	}
	return nil, fmt.Errorf(`Resource "%s" was not found!`, filename)
}
