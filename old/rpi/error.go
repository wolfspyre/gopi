/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016
	All Rights Reserved

	For Licensing and Usage information, please see LICENSE.md
*/
package rpi

import (
	"errors"
)

var (
	ErrorInit           = errors.New("init error")
	ErrorGenCmd         = errors.New("vcgencmd error")
	ErrorResponse       = errors.New("unexpected response")
	ErrorUnknownProduct = errors.New("unknown product")
	ErrorUnknownPin     = errors.New("unknown pin")
)