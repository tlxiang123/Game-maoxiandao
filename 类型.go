package main

type Rect struct {
	X1 int
	Y1 int
	X2 int
	Y2 int
}

type Point struct {
	X int
	Y int
}

type Pic struct {
	Name        string
	X1          int
	Y1          int
	X2          int
	Y2          int
	PicPath     string
	Gray        bool
	Transparent bool
	Sim         float32
	Dir         int
}

type CColor struct {
	Name  string
	X     int
	Y     int
	Color string
	Sim   float32
}

type FColor struct {
	Name  string
	X1    int
	Y1    int
	X2    int
	Y2    int
	Color string
	Sim   float32
	Dir   int
}

type CCRegion struct {
	Name  string
	X1    int
	Y1    int
	X2    int
	Y2    int
	Color string
	Sim   float32
}

type DMColor struct {
	Name   string
	Colors string
	Sim    float32
}

type FMColor struct {
	Name        string
	X1          int
	Y1          int
	X2          int
	Y2          int
	MainColor   string
	OffsetColor string
	Sim         float32
	Dir         int
}

type FStr struct {
	Name        string
	X1          int
	Y1          int
	X2          int
	Y2          int
	String      string
	ColorFormat string
	Sim         float32
	DictName    string
}

type SOcr struct {
	Name        string
	X1          int
	Y1          int
	X2          int
	Y2          int
	ColorFormat string
	Sim         float32
	DictName    string
}

type PPOcrRegion struct {
	Name     string
	X1       int
	Y1       int
	X2       int
	Y2       int
	Color    string
	Version  string
	Contains string
	MinScore float64
}
