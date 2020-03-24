package service

import (
	"bufio"
	"bytes"
	texttospeech "cloud.google.com/go/texttospeech/apiv1"
	"context"
	"fmt"
	"github.com/faiface/beep"
	"github.com/faiface/beep/speaker"
	"github.com/faiface/beep/wav"
	texttospeechpb "google.golang.org/genproto/googleapis/cloud/texttospeech/v1"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"
)

///	TEXT TO SPEECH ///

/// GO TO LINE #164 TO FIND SPEECH TO TEXT ///

/// Structs ///

type SpeechRequest struct {
	Text         string
	LanguageCode string
	SsmlGender   string
	VoiceName    string
}

type SpeechExampleError struct {
	Message string
}

/// Core func's ///

//SpeakToFile uses the values within the receiver to open a connection to GCP, create a request, and then take the response and put it into a file.
func (st *SpeechRequest) SpeakToFile(outputFile string) {

	//Create a go context, a key component of nearly all golang web requests
	ctx := context.Background()

	//create a client connection to the GCP TTS backend
	client, err := texttospeech.NewClient(ctx)
	checkErr(err)

	// we 'defer' the closing of the client connection until the function has exited
	defer client.Close()

	//Craft a request using the parameters specified
	req, serr := st.CraftTextSpeechRequest()
	checkSpeechErr(serr)

	//Receive a response from GCP TTS
	resp, err := client.SynthesizeSpeech(ctx, &req)
	checkErr(err)

	//Write the contents of the response body to a file

	err = ioutil.WriteFile(outputFile, resp.AudioContent, 0644)
	log.Println(err)
	checkErr(err)

	fmt.Printf("TTS Successfully written to %s\n\n", outputFile)
}

func (pt *SpeechRequest) SpeakFromFileToFile(inputFile string, outputFile string) {

	//Read the contents of the input file and pass them to the struct
	content, err := ioutil.ReadFile(inputFile)
	checkErr(err)
	pt.Text = string(content)

	//Proceed to translate the contents
	pt.SpeakToFile(outputFile)
}

func (st *SpeechRequest) SpeakToStream() []byte {

	ctx := context.Background()
	client, err := texttospeech.NewClient(ctx)
	checkErr(err)
	defer client.Close()

	req, er := st.CraftTextSpeechRequest()
	checkSpeechErr(er)

	resp, err := client.SynthesizeSpeech(ctx, &req)
	checkErr(err)

	return resp.AudioContent
}

/// Supporting funcs ///

func (st *SpeechRequest) CraftTextSpeechRequest() (texttospeechpb.SynthesizeSpeechRequest, SpeechExampleError) {

	//make sure st has the required values

	if st.Text == "" {
		return texttospeechpb.SynthesizeSpeechRequest{}, SpeechExampleError{Message: "TTS Request Has Empty Text"}
	}

	if st.LanguageCode == "" {
		return texttospeechpb.SynthesizeSpeechRequest{}, SpeechExampleError{Message: "TTS Request Has Empty Language Code"}
	}

	if st.SsmlGender == "" {
		return texttospeechpb.SynthesizeSpeechRequest{}, SpeechExampleError{Message: "TTS Request Has Empty Ssml Gender"}
	}

	if st.VoiceName == "" {
		return texttospeechpb.SynthesizeSpeechRequest{}, SpeechExampleError{Message: "TTS Request Has Empty Voice Name"}
	}

	// convert input strings to the proper types
	gender := texttospeechpb.SsmlVoiceGender_FEMALE
	if strings.Contains(st.SsmlGender, "MALE") {
		gender = texttospeechpb.SsmlVoiceGender_MALE
	}
	if st.SsmlGender == "" {
		gender = texttospeechpb.SsmlVoiceGender_NEUTRAL
	}

	input := &texttospeechpb.SynthesisInput{InputSource: &texttospeechpb.SynthesisInput_Text{Text: st.Text}}
	if strings.Contains(st.Text, "<speak>") {
		input = &texttospeechpb.SynthesisInput{InputSource: &texttospeechpb.SynthesisInput_Ssml{Ssml: st.Text}}
	}

	//create and return the request
	return texttospeechpb.SynthesizeSpeechRequest{
		AudioConfig: &texttospeechpb.AudioConfig{
			AudioEncoding:        texttospeechpb.AudioEncoding_LINEAR16,
			SpeakingRate:         0,
			Pitch:                0,
			VolumeGainDb:         0,
			SampleRateHertz:      16000,
			EffectsProfileId:     nil,
			XXX_NoUnkeyedLiteral: struct{}{},
			XXX_unrecognized:     nil,
			XXX_sizecache:        0,
		},

		Voice: &texttospeechpb.VoiceSelectionParams{
			LanguageCode: st.LanguageCode,
			Name:         st.VoiceName,
			SsmlGender:   gender,
		},

		Input: input,
	}, SpeechExampleError{}

}

//SpeakAloud is a function which accepts a string of text as its only parameter, it then matches that text
//to a file name within our audio store. If i cannot find the audio file it generates it and stores it.
func SpeakAloud(text string) bool {

	if !CheckForFile("audio/" + text) {
		req := SpeechRequest{
			Text:         text,
			LanguageCode: "en-US",
			SsmlGender:   "FEMALE",
			VoiceName:    "en-us-Wavenet-C",
		}

		log.Println("Attempting to write to file, dir = audio/\"" + text + "\".wav")
		req.SpeakToFile("/home/pi/Hardware/audio/\"" + text + "\".wav")

	}

	f, err := os.Open("/home/pi/Hardware/audio/\"" + text + "\".wav")
	if err != nil {
		log.Println(err)
		log.Fatal("Could not complete required audio I/O.")
	}

	streamer, format, err := wav.Decode(f)
	if err != nil {
		log.Println(err)
		return false
	}

	defer streamer.Close()

	//start the speaker, specify the files sample rate and the buffer size. larger the buffer the better the stability, the smaller the lower the latency.
	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second*2/3))
	speaker.Play(streamer)
	log.Println("Beginning To speak. \n\n")

	done := make(chan bool)
	speaker.Play(beep.Seq(streamer, beep.Callback(func() {
		done <- true
	})))

	<-done
	log.Println("Finished Speaking.")
	return true
}

func SpeakAloudWithChan(text string, done chan bool) {

	if !CheckForFile("audio/" + text) {
		req := SpeechRequest{
			Text:         text,
			LanguageCode: "en-US",
			SsmlGender:   "FEMALE",
			VoiceName:    "en-us-Wavenet-C",
		}

		log.Println("Attempting to write to file, dir = audio/\"" + text + "\".wav")
		req.SpeakToFile("/home/pi/Hardware/audio/\"" + text + "\".wav")

	}

	f, err := os.Open("/home/pi/Hardware/audio/\"" + text + "\".wav")
	if err != nil {
		log.Println(err)
		log.Fatal("Could not complete required audio I/O.")
	}

	streamer, format, err := wav.Decode(f)
	if err != nil {
		log.Println(err)
	}

	defer streamer.Close()

	//start the speaker, specify the files sample rate and the buffer size. larger the buffer the better the stability, the smaller the lower the latency.
	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second*2/3))
	speaker.Play(streamer)
	log.Println("Beginning To speak. \n\n")

	speaker.Play(beep.Seq(streamer, beep.Callback(func() {
		done <- true
	})))

}

//short hand for quick error checks
func checkErr(e error) {
	if e != nil {
		fmt.Print(e)
	}
}

//short hand for speechExampleError checks
func checkSpeechErr(exampleError SpeechExampleError) {
	if exampleError.Message != "" {
		fmt.Println(exampleError.Message)
	}
}

func TranscriptionConfidence(transcription string, exact string) float64 {

	wordsTranscribed := strings.Split(transcription, " ")
	wordsExpected := strings.Split(exact, " ")
	totalWords := float64(len(wordsTranscribed))
	wordsFound := 0.0
	//a fancy golang foreach
	for wordTranscribed := range wordsTranscribed {
		for wordExpected := range wordsExpected {
			if wordTranscribed == wordExpected {
				wordsFound++
			}
		}
	}

	return wordsFound / totalWords
}

/// SPEECH TO TEXT ///

//This function calls the google speech to text api
//We use gstreamer to get audio input from the rpi
//we then pass that audio stream to the Stdin of another script
//which then formats the audio stream and passes it to the google cloud platform
//We continue to stream to the GCP until we receive a non-empty response body
//at which point we return the contents and kill the streaming process.
func Recognize() (string, float64, error) {

	//First we need to craft the command we want to execute
	cmdName := "/home/pi/Hardware/service/recognize"
	cmdArgs := []string{}
	cmd := exec.Command(cmdName, cmdArgs...)

	//If we want to return the values that are returned from the above script
	//we need to declare the return values at a higher scope
	transcription := ""
	confidence := 0.0

	stderr, err := cmd.StderrPipe() //Stderr for whatever reason.

	//We need to start our goroutine from the main thread
	fmt.Println("STARTING.")
	err = cmd.Start()
	if err != nil {
		fmt.Println("Error starting Cmd", err)
		os.Exit(1)
	}

	go func() {
		//We need to create a reader for the stdout of this script
		//A scanner is created to read the stdout of the above command
		scanner := bufio.NewScanner(stderr)

		for scanner.Scan() {
			fmt.Println("Response Recognized ON STD OUT...")
			fmt.Println(scanner.Text())
			//A Third party can interrupt this streaming process by simply saying "stop"
			//useful when you want to stop the test, but don't want orphan processes
			if strings.Contains(scanner.Text(), "stop") {
				if err := cmd.Process.Kill(); err != nil {
					log.Fatal("failed to kill process: ", err)
				}
			}

			//Regular expressions are used to parse out the transcription and confidence score
			//from the return body of our API request
			transcriptRegx := regexp.MustCompile("(\"([^\"]|\"\")*\")")
			match := transcriptRegx.FindStringSubmatch(scanner.Text())

			confRegx := regexp.MustCompile("([+-]?[0-9]*\\.[0-9]*)")
			conf := confRegx.FindStringSubmatch(scanner.Text())

			//We need to ensure that an empty response body doesn't stop our recognition
			if len(match) != 0 {
				transcription = match[1]
				confidence, err = strconv.ParseFloat(conf[1], 64)
				if err != nil {
					fmt.Println("Could not convert confidence from string to float64")
					fmt.Println(conf[1])
					fmt.Println(match[1])
				}

				//Now that we have our transcription we can stop the recognition process
				if err := cmd.Process.Signal(os.Kill); err != nil {
					log.Println(cmd.ProcessState.String())
				}
			}
		}
	}()

	//We need to wait for a transcription before we can return said transcription
	err = cmd.Wait()
	var out bytes.Buffer
	if err != nil && transcription == "" {
		fmt.Println("Recognition crash")
		fmt.Println(fmt.Sprint(err) + ": " + out.String())
		return transcription, confidence, err
	}
	fmt.Println(cmd.ProcessState.String())
	return transcription, confidence, nil
}

/// Supporting Funcs ///

func c(err error) {
	if err != nil {
		panic(err)
	}

}
