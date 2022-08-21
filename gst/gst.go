// Package gst provides an easy API to create an appsink pipeline
package gst

/*
#cgo pkg-config: gstreamer-1.0 gstreamer-app-1.0

#include "gst.h"

*/
import "C"
import (
	"fmt"
	"log"
	"sync"
	"time"
	"unsafe"

	"github.com/pion/webrtc/v3"
	"github.com/pion/webrtc/v3/pkg/media"
)

func init() {
	go C.gstreamer_send_start_mainloop()
}

// Pipeline is a wrapper for a GStreamer Pipeline
type Pipeline struct {
	Pipeline  *C.GstElement
	tracks    []*webrtc.TrackLocalStaticSample
	id        int
	codecName string
	clockRate float32
}

var pipelines = make(map[int]*Pipeline)
var pipelinesLock sync.Mutex

const (
	videoClockRate = 90000
	audioClockRate = 48000
	pcmClockRate   = 8000
)

// CreatePipeline creates a GStreamer Pipeline
func CreatePipeline(codecName string, tracks []*webrtc.TrackLocalStaticSample, pipelineSrc string) *Pipeline {
	pipelineStr := "appsink name=appsink"
	var clockRate float32

	pipelineStr = "rpicamsrc annotation-mode=12 bitrate=256000 preview=false ! video/x-h264, framerate=1/1, stream-format=byte-stream, profile=constrained-baseline, level=3.0 ! queue max-size-time=100000000 ! h264parse ! video/x-h264, alignment=au, stream-format=byte-stream ! " + pipelineStr

	// works ok, but eats through cpu
	//pipelineStr = "rpicamsrc bitrate=256000 preview=false ! video/x-raw, width=640, height=480, framerate=30/1 ! x264enc speed-preset=ultrafast tune=zerolatency key-int-max=20  ! video/x-h264,stream-format=byte-stream ! " + pipelineStr

	// works ok, but its a test source
	//pipelineStr = "videotestsrc ! video/x-raw,format=I420 ! x264enc speed-preset=ultrafast tune=zerolatency key-int-max=20 ! video/x-h264,stream-format=byte-stream ! " + pipelineStr

	//pipelineStr = pipelineSrc + " ! " + pipelineStr
	clockRate = videoClockRate

	print("Pipeline created: " + pipelineStr)

	pipelineStrUnsafe := C.CString(pipelineStr)
	defer C.free(unsafe.Pointer(pipelineStrUnsafe))

	pipelinesLock.Lock()
	defer pipelinesLock.Unlock()

	pipeline := &Pipeline{
		Pipeline:  C.gstreamer_send_create_pipeline(pipelineStrUnsafe),
		tracks:    tracks,
		id:        len(pipelines),
		codecName: codecName,
		clockRate: clockRate,
	}

	pipelines[pipeline.id] = pipeline
	return pipeline
}

// Start starts the GStreamer Pipeline
func (p *Pipeline) Start() {
	C.gstreamer_send_start_pipeline(p.Pipeline, C.int(p.id))
}

// Stop stops the GStreamer Pipeline
func (p *Pipeline) Stop() {
	C.gstreamer_send_stop_pipeline(p.Pipeline)
}

//export goHandlePipelineBuffer
func goHandlePipelineBuffer(buffer unsafe.Pointer, bufferLen C.int, duration C.int, pipelineID C.int) {
	pipelinesLock.Lock()
	pipeline, ok := pipelines[int(pipelineID)]
	pipelinesLock.Unlock()

	if ok {
		for _, t := range pipeline.tracks {
			sampleData := C.GoBytes(buffer, bufferLen)
			sampleDuration := time.Duration(duration)
			log.Printf("Handling buffers, buffer len %d, duration %d mks", int(bufferLen), sampleDuration.Microseconds())
			if err := t.WriteSample(media.Sample{Data: sampleData, Duration: sampleDuration}); err != nil {
				panic(err)
			}
		}
	} else {
		fmt.Printf("discarding buffer, no pipeline with id %d", int(pipelineID))
	}
	C.free(buffer)
}
