package yolo

/*
#cgo LDFLAGS: -ldl
#cgo LDFLAGS: -L${SRCDIR}/../../resources/libs/x86_64 -lyolo
#cgo arm64 LDFLAGS: -L${SRCDIR}/../../resources/libs/arm64-v8a -lyolo
#cgo amd64 LDFLAGS: -L${SRCDIR}/../../resources/libs/x86_64 -lyolo
#cgo 386 LDFLAGS: -L${SRCDIR}/../../resources/libs/x86 -lyolo
#include <stdlib.h>

extern void* newYolo();
extern char* loadModelYolo(void* Obj, const char* Version, const char* Param, const char* Bin, const char* Labels, int Thread);
extern char* detectYolo(void* Obj, unsigned char *Byte, int W, int H, float Nms, float Prob, int Size);
extern void closeYolo(void* Obj);
*/
import "C"

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"image"
	"image/draw"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"strings"
	"unsafe"

	"github.com/Dasongzi1366/AutoGo/images"
)

const (
	defaultThreshold    = 0.25
	defaultNmsThreshold = 0.45
	defaultSize         = 640
)

type Yolo struct {
	cPtr         unsafe.Pointer
	img          *image.NRGBA
	threshold    float32
	nmsThreshold float32
	size         int
}

// Result 表示对象检测的结果，包括位置、标签和置信度。
type Result struct {
	X       int     `json:"X"`
	Y       int     `json:"Y"`
	Width   int     `json:"宽"`
	Height  int     `json:"高"`
	Label   string  `json:"标签"`
	Score   float64 `json:"精度"`
	CenterX int     `json:"-"` // 中心坐标X
	CenterY int     `json:"-"` // 中心坐标Y
}

// New 创建一个新的 YOLO 实例，并加载模型和标签。
func New(version string, cpuThreadNum int, paramPath, binPath, labels string) *Yolo {
	cPtr := C.newYolo()
	if cPtr == nil {
		return nil
	}

	cVersion := C.CString(version)
	cParam := C.CString(paramPath)
	cBin := C.CString(binPath)
	cLabels := C.CString(normalizeLabels(labels))
	defer C.free(unsafe.Pointer(cVersion))
	defer C.free(unsafe.Pointer(cParam))
	defer C.free(unsafe.Pointer(cBin))
	defer C.free(unsafe.Pointer(cLabels))

	errText := takeCString(C.loadModelYolo(cPtr, cVersion, cParam, cBin, cLabels, C.int(cpuThreadNum)))
	if errText != "OK" {
		C.closeYolo(cPtr)
		return nil
	}

	return &Yolo{
		cPtr:         cPtr,
		threshold:    defaultThreshold,
		nmsThreshold: defaultNmsThreshold,
		size:         defaultSize,
	}
}

// SetImage 设置一个图片对象作为下次 Detect 方法的原始图像。
func (y *Yolo) SetImage(img *image.NRGBA) {
	if y == nil {
		return
	}
	y.img = img
}

// Detect 在指定的屏幕区域执行目标检测。
func (y *Yolo) Detect(x1, y1, x2, y2, displayId int) []Result {
	if y == nil {
		return nil
	}
	if y.img != nil {
		return y.DetectFromImage(y.img)
	}
	img := images.CaptureScreen(x1, y1, x2, y2, displayId)
	return y.DetectFromImage(img)
}

// DetectFromImage 从内存中的图像进行识别。
func (y *Yolo) DetectFromImage(img *image.NRGBA) []Result {
	if y == nil || y.cPtr == nil || img == nil {
		return nil
	}
	w := img.Rect.Dx()
	h := img.Rect.Dy()
	if w <= 0 || h <= 0 || len(img.Pix) == 0 {
		return nil
	}
	text := takeCString(C.detectYolo(
		y.cPtr,
		(*C.uchar)(unsafe.Pointer(&img.Pix[0])),
		C.int(w),
		C.int(h),
		C.float(y.nmsThreshold),
		C.float(y.threshold),
		C.int(y.size),
	))
	return parseResults(text)
}

// DetectFromBase64 从 Base64 编码的图像进行识别。
func (y *Yolo) DetectFromBase64(b64 string) []Result {
	data, err := base64.StdEncoding.DecodeString(b64)
	if err != nil {
		return nil
	}
	img, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return nil
	}
	return y.DetectFromImage(toNRGBA(img))
}

// DetectFromPath 从文件路径进行识别。
func (y *Yolo) DetectFromPath(path string) []Result {
	file, err := os.Open(path)
	if err != nil {
		return nil
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		return nil
	}
	return y.DetectFromImage(toNRGBA(img))
}

// Close 关闭 YOLO 模型实例，释放相关资源。
func (y *Yolo) Close() {
	if y == nil || y.cPtr == nil {
		return
	}
	C.closeYolo(y.cPtr)
	y.cPtr = nil
}

func takeCString(ptr *C.char) string {
	if ptr == nil {
		return ""
	}
	text := C.GoString(ptr)
	C.free(unsafe.Pointer(ptr))
	return text
}

func normalizeLabels(labels string) string {
	text := strings.TrimSpace(labels)
	if text == "" {
		return text
	}
	if data, err := os.ReadFile(text); err == nil {
		lines := strings.Split(string(data), "\n")
		out := make([]string, 0, len(lines))
		for _, line := range lines {
			line = strings.TrimSpace(strings.TrimRight(line, "\r"))
			if line != "" {
				out = append(out, line)
			}
		}
		if len(out) > 0 {
			return strings.Join(out, ",")
		}
	}
	return text
}

func parseResults(text string) []Result {
	text = strings.TrimSpace(text)
	if text == "" || !strings.HasPrefix(text, "[") {
		return nil
	}
	var results []Result
	if err := json.Unmarshal([]byte(text), &results); err != nil {
		return nil
	}
	for i := range results {
		results[i].CenterX = results[i].X + results[i].Width/2
		results[i].CenterY = results[i].Y + results[i].Height/2
	}
	return results
}

func toNRGBA(src image.Image) *image.NRGBA {
	if src == nil {
		return nil
	}
	if nrgba, ok := src.(*image.NRGBA); ok {
		return nrgba
	}
	dst := image.NewNRGBA(src.Bounds())
	draw.Draw(dst, dst.Bounds(), src, src.Bounds().Min, draw.Src)
	return dst
}
