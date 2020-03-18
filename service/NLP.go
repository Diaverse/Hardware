package service

import (
	"bufio"
	"bytes"
	texttospeech "cloud.google.com/go/texttospeech/apiv1"
	"context"
	"fmt"
	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
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

type SpeechRequest struct ***REMOVED***
	Text         string
	LanguageCode string
	SsmlGender   string
	VoiceName    string
***REMOVED***

type SpeechExampleError struct ***REMOVED***
	Message string
***REMOVED***

/// Core func's ///

//SpeakToFile uses the values within the receiver to open a connection to GCP, create a request, and then take the response and put it into a file.
func (st *SpeechRequest) SpeakToFile(outputFile string) ***REMOVED***

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
	checkErr(err)

	fmt.Printf("TTS Successfully written to %s", outputFile)
***REMOVED***

func (pt *SpeechRequest) SpeakFromFileToFile(inputFile string, outputFile string) ***REMOVED***

	//Read the contents of the input file and pass them to the struct
	content, err := ioutil.ReadFile(inputFile)
	checkErr(err)
	pt.Text = string(content)

	//Proceed to translate the contents
	pt.SpeakToFile(outputFile)
***REMOVED***

func (st *SpeechRequest) SpeakToStream() []byte ***REMOVED***

	ctx := context.Background()
	client, err := texttospeech.NewClient(ctx)
	checkErr(err)
	defer client.Close()

	req, er := st.CraftTextSpeechRequest()
	checkSpeechErr(er)

	resp, err := client.SynthesizeSpeech(ctx, &req)
	checkErr(err)

	return resp.AudioContent
***REMOVED***

/// Supporting funcs ///

func (st *SpeechRequest) CraftTextSpeechRequest() (texttospeechpb.SynthesizeSpeechRequest, SpeechExampleError) ***REMOVED***

	//make sure st has the required values

	if st.Text == "" ***REMOVED***
		return texttospeechpb.SynthesizeSpeechRequest***REMOVED******REMOVED***, SpeechExampleError***REMOVED***Message: "TTS Request Has Empty Text"***REMOVED***
	***REMOVED***

	if st.LanguageCode == "" ***REMOVED***
		return texttospeechpb.SynthesizeSpeechRequest***REMOVED******REMOVED***, SpeechExampleError***REMOVED***Message: "TTS Request Has Empty Language Code"***REMOVED***
	***REMOVED***

	if st.SsmlGender == "" ***REMOVED***
		return texttospeechpb.SynthesizeSpeechRequest***REMOVED******REMOVED***, SpeechExampleError***REMOVED***Message: "TTS Request Has Empty Ssml Gender"***REMOVED***
	***REMOVED***

	if st.VoiceName == "" ***REMOVED***
		return texttospeechpb.SynthesizeSpeechRequest***REMOVED******REMOVED***, SpeechExampleError***REMOVED***Message: "TTS Request Has Empty Voice Name"***REMOVED***
	***REMOVED***

	// convert input strings to the proper types
	gender := texttospeechpb.SsmlVoiceGender_FEMALE
	if strings.Contains(st.SsmlGender, "MALE") ***REMOVED***
		gender = texttospeechpb.SsmlVoiceGender_MALE
	***REMOVED***
	if st.SsmlGender == "" ***REMOVED***
		gender = texttospeechpb.SsmlVoiceGender_NEUTRAL
	***REMOVED***

	input := &texttospeechpb.SynthesisInput***REMOVED***InputSource: &texttospeechpb.SynthesisInput_Text***REMOVED***Text: st.Text***REMOVED******REMOVED***
	if strings.Contains(st.Text, "<speak>") ***REMOVED***
		input = &texttospeechpb.SynthesisInput***REMOVED***InputSource: &texttospeechpb.SynthesisInput_Ssml***REMOVED***Ssml: st.Text***REMOVED******REMOVED***
	***REMOVED***

	//create and return the request
	return texttospeechpb.SynthesizeSpeechRequest***REMOVED***
		AudioConfig: &texttospeechpb.AudioConfig***REMOVED***
			AudioEncoding: texttospeechpb.AudioEncoding_MP3,
		***REMOVED***,

		Voice: &texttospeechpb.VoiceSelectionParams***REMOVED***
			LanguageCode: st.LanguageCode,
			Name:         st.VoiceName,
			SsmlGender:   gender,
		***REMOVED***,

		Input: input,
	***REMOVED***, SpeechExampleError***REMOVED******REMOVED***

***REMOVED***

func SpeakAloud(text string) ***REMOVED***

	if !CheckForFile("/audio/" + text) ***REMOVED***
		SpeechRequest***REMOVED***
			Text:         text,
			LanguageCode: "en-US",
			SsmlGender:   "FEMALE",
			VoiceName:    "en-us-Wavenet-C",
		***REMOVED***.SpeakToFile("audio/" + text + ".mp3")
	***REMOVED***

	f, err := os.Open("audio/" + text)
	if err != nil ***REMOVED***
		log.Fatal("Could not complete required audio I/O.")
	***REMOVED***
	streamer, format, err := mp3.Decode(f)
	if err != nil ***REMOVED***
		log.Fatal("Could not construct streamer from decoded input file")
	***REMOVED***

	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
	done := make(chan bool)
	speaker.Play(streamer, beep.Callback(func() ***REMOVED***
		done <- true
	***REMOVED***))

	//wait for the file to stop playing
	<-done

	//must be called.
	streamer.Close()
***REMOVED***

//short hand for quick error checks
func checkErr(e error) ***REMOVED***
	if e != nil ***REMOVED***
		fmt.Print(e)
	***REMOVED***
***REMOVED***

//short hand for speechExampleError checks
func checkSpeechErr(exampleError SpeechExampleError) ***REMOVED***
	if exampleError.Message != "" ***REMOVED***
		fmt.Println(exampleError.Message)
	***REMOVED***
***REMOVED***

func TranscriptionConfidence(transcription string, exact string) float64 ***REMOVED***

	wordsTranscribed := strings.Split(transcription, " ")
	wordsExpected := strings.Split(exact, " ")
	totalWords := float64(len(wordsTranscribed))
	wordsFound := 0.0
	//a fancy golang foreach
	for wordTranscribed := range wordsTranscribed ***REMOVED***
		for wordExpected := range wordsExpected ***REMOVED***
			if wordTranscribed == wordExpected ***REMOVED***
				wordsFound++
			***REMOVED***
		***REMOVED***
	***REMOVED***

	return wordsFound / totalWords
***REMOVED***

/// SPEECH TO TEXT ///

//This function calls the google speech to text api
//We use gstreamer to get audio input from the rpi
//we then pass that audio stream to the Stdin of another script
//which then formats the audio stream and passes it to the google cloud platform
//We continue to stream to the GCP until we receive a non-empty response body
//at which point we return the contents and kill the streaming process.
func Recognize() (string, float64, error) ***REMOVED***

	//First we need to craft the command we want to execute
	cmdName := "SpeechToTextExamples/scripts/recognize"
	cmdArgs := []string***REMOVED***""***REMOVED***
	cmd := exec.Command(cmdName, cmdArgs...)
	//We need to create a reader for the stdout of this script
	cmdReader, err := cmd.StdoutPipe()
	if err != nil ***REMOVED***
		fmt.Println("Error creating StdoutPipe for Cmd", err)
		os.Exit(1)
	***REMOVED***

	//If we want to return the values that are returned from the above script
	//we need to declare the return values at a higher scope
	transcription := ""
	confidence := 0.0

	//A scanner is created to read the stdout of the above command
	scanner := bufio.NewScanner(cmdReader)

	//A new go thread is created to handle the audio streaming and subsequent response bodies
	go func() ***REMOVED***
		for scanner.Scan() ***REMOVED***

			fmt.Println("Response Recognized...")
			//A Third party can interrupt this streaming process by simply saying "stop"
			//useful when you want to stop the test, but don't want orphan processes
			if strings.Contains(scanner.Text(), "stop") ***REMOVED***
				if err := cmd.Process.Kill(); err != nil ***REMOVED***
					log.Fatal("failed to kill process: ", err)
				***REMOVED***
			***REMOVED***

			//Regular expressions are used to parse out the transcription and confidence score
			//from the return body of our API request
			transcriptRegx := regexp.MustCompile("(\"([^\"]|\"\")*\")")
			match := transcriptRegx.FindStringSubmatch(scanner.Text())

			confRegx := regexp.MustCompile("([+-]?[0-9]*\\.[0-9]*)")
			conf := confRegx.FindStringSubmatch(scanner.Text())

			//We need to ensure that an empty response body doesn't stop our recognition
			if len(match) != 0 ***REMOVED***
				transcription = match[1]
				confidence, err = strconv.ParseFloat(conf[1], 64)
				if err != nil ***REMOVED***
					fmt.Println("Could not convert confidence from string to float64")
					fmt.Println(conf[1])
					fmt.Println(match[1])
				***REMOVED***

				//Now that we have our transcription we can stop the recognition process
				if err := cmd.Process.Kill(); err != nil ***REMOVED***
					log.Fatal("failed to kill process: ", err)
				***REMOVED***
			***REMOVED***
		***REMOVED***
	***REMOVED***()

	//We need to start our goroutine from the main thread
	err = cmd.Start()
	if err != nil ***REMOVED***
		fmt.Println("Error starting Cmd", err)
		os.Exit(1)
	***REMOVED***

	//We need to wait for a transcription before we can return said transcription
	err = cmd.Wait()
	var out bytes.Buffer
	if err != nil && transcription == "" ***REMOVED***
		fmt.Println("Recognition crash")
		fmt.Println("tried to run the command : ./scripts/recognize")
		fmt.Println(fmt.Sprint(err) + ": " + out.String())
		return transcription, confidence, err
	***REMOVED***
	return transcription, confidence, nil
***REMOVED***

/// Supporting Funcs ///

func c(err error) ***REMOVED***
	if err != nil ***REMOVED***
		panic(err)
	***REMOVED***

***REMOVED***