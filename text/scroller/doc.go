// Tideland Go Library - Text - Scroller
//
// Copyright (C) 2018 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

// Package scroller helps analyzing a continuously written line by line text
// content, e.g. at the monitoring of log files. Here the Scroller is working
// in the background and allows to read out of any ReadSeeker (which may be
// a File) from beginning, end or a given number of lines before the end,
// filter the output by a filter function and write it into a Writer. If
// a number of lines and a filter are passed the Scroller tries to find that
// number of lines matching to the filter.
package scroller

// EOF
