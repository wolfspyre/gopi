/*
Go Language Raspberry Pi Interface
(c) Copyright David Thorpe 2016
All Rights Reserved

For Licensing and Usage information, please see LICENSE.md
*/

package util /* import "github.com/djthorpe/gopi/util" */

import (
	"errors"
)

////////////////////////////////////////////////////////////////////////////////
// GLOBAL VARIABLES

var (
	ErrUnsupportedType = errors.New("Unsupported type")
	ErrParseError      = errors.New("Syntax error in input")
)
