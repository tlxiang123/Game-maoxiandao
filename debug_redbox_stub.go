//go:build !android

package main

func addDebugRedBox(x1, y1, x2, y2 int) {}

func addDebugPointBox(x, y int) {}

func renderDebugRedBoxes() {}

func (z *Zg) markScreenRect(x1, y1, x2, y2 int) {}

func (z *Zg) markScreenPoint(x, y int) {}

func (z *Zg) markClickPoint(x, y int) {}
