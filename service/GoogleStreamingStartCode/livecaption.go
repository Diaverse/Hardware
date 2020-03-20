package main

import (
	speech "cloud.google.com/go/speech/apiv1"
	"context"
	speechpb "google.golang.org/genproto/googleapis/cloud/speech/v1"
	"io"
	"log"
	"os"
)

// Scripts to handle audio transcription and maintenance of TTS files.
// This directory only contains a single go file, livecaption, and as such you will
// need to inspect any bash scripts directly from the github repository
//
//
// The single go file, livecaption, was initially provided from the golang samples repository
// https://github.com/GoogleCloudPlatform/golang-samples/tree/master/speech/livecaption
//
// The file has been modified such that only a single transcription is processed before the command exits.
// This command relies on several audio libraries such as gstreamer and pulse audio.
// You will need to install said libraries before this file can be used.
// The command to start the transcription process is as follows:
//
//    $ gst-launch-1.0 -v pulsesrc ! audioconvert ! audioresample ! audio/x-raw,channels=1,rate=16000 ! filesink location=/dev/stdout | livecaption

func main() ***REMOVED***
	ctx := context.Background()

	client, err := speech.NewClient(ctx)
	if err != nil ***REMOVED***
		log.Fatal(err)
	***REMOVED***
	stream, err := client.StreamingRecognize(ctx)
	if err != nil ***REMOVED***
		log.Fatal(err)
	***REMOVED***
	// Send the initial configuration message.
	if err := stream.Send(&speechpb.StreamingRecognizeRequest***REMOVED***
		StreamingRequest: &speechpb.StreamingRecognizeRequest_StreamingConfig***REMOVED***
			StreamingConfig: &speechpb.StreamingRecognitionConfig***REMOVED***
				Config: &speechpb.RecognitionConfig***REMOVED***
					Encoding:        speechpb.RecognitionConfig_LINEAR16,
					SampleRateHertz: 16000,
					LanguageCode:    "en-US",
				***REMOVED***,
			***REMOVED***,
		***REMOVED***,
	***REMOVED***); err != nil ***REMOVED***
		log.Fatal(err)
	***REMOVED***

	go func() ***REMOVED***
		// Pipe stdin to the API.
		buf := make([]byte, 1024)
		for ***REMOVED***
			n, err := os.Stdin.Read(buf)
			if n > 0 ***REMOVED***
				if err := stream.Send(&speechpb.StreamingRecognizeRequest***REMOVED***
					StreamingRequest: &speechpb.StreamingRecognizeRequest_AudioContent***REMOVED***
						AudioContent: buf[:n],
					***REMOVED***,
				***REMOVED***); err != nil ***REMOVED***
					log.Printf("Could not send audio: %v", err)
				***REMOVED***
			***REMOVED***
			if err == io.EOF ***REMOVED***
				// Nothing else to pipe, close the stream.
				if err := stream.CloseSend(); err != nil ***REMOVED***
					log.Fatalf("Could not close stream: %v", err)
				***REMOVED***
				return
			***REMOVED***
			if err != nil ***REMOVED***
				log.Printf("Could not read from stdin: %v", err)
				continue
			***REMOVED***
		***REMOVED***
	***REMOVED***()

	for ***REMOVED***
		resp, err := stream.Recv()
		if err == io.EOF ***REMOVED***
			break
		***REMOVED***
		if err != nil ***REMOVED***
			log.Fatalf("Cannot stream results: %v", err)
		***REMOVED***
		if err := resp.Error; err != nil ***REMOVED***
			// Workaround while the API doesn't give a more informative error.
			if err.Code == 3 || err.Code == 11 ***REMOVED***
				log.Print("WARNING: Speech recognition request exceeded limit of 60 seconds.")
			***REMOVED***
			log.Fatalf("Could not recognize: %v", err)
		***REMOVED***
		for _, result := range resp.Results ***REMOVED***
			log.Printf("%+v\n", result)

		***REMOVED***

	***REMOVED***
***REMOVED***
