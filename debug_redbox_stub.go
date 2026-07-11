//go:build !android

package main

func addDebugRedBox(x1, y1, x2, y2 int) {}

func addDebugPointBox(x, y int) {}

func 暂停调试红框() {}

func 恢复调试红框() {}

func renderDebugRedBoxes() {}

func (z *Zg) markScreenRect(x1, y1, x2, y2 int) {}

func (z *Zg) markScreenPoint(x, y int) {}

func (z *Zg) markClickPoint(x, y int) {}
