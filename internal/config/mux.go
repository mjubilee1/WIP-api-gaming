package config

import (
	"log"
	"api-gaming/internal/util"
	"github.com/muxinc/mux-go"
)

var (
	muxid = util.ViperEnvVariable("MUX_TOKEN_ID")
	muxtoken = util.ViperEnvVariable("MUX_TOKEN_SECRET")
	stream muxgo.LiveStreamResponse
	muxClient *muxgo.APIClient
)

// InitMux - initialize mux video streaming
func InitMux() *muxgo.APIClient {
	muxClient = muxgo.NewAPIClient(
		muxgo.NewConfiguration(
				muxgo.WithBasicAuth(muxid, muxtoken),
		))

	return muxClient
}

// CreateLiveStream - init and create live streaming response
func CreateLiveStream() muxgo.LiveStreamResponse {
	muxClient = muxgo.NewAPIClient(
		muxgo.NewConfiguration(
				muxgo.WithBasicAuth(muxid, muxtoken),
		))

	// Create live stream
	car := muxgo.CreateAssetRequest{PlaybackPolicy: []muxgo.PlaybackPolicy{muxgo.PUBLIC}}
	csr := muxgo.CreateLiveStreamRequest{NewAssetSettings: car, PlaybackPolicy: []muxgo.PlaybackPolicy{muxgo.PUBLIC}}
	stream, err := muxClient.LiveStreamsApi.CreateLiveStream(csr)
	if err != nil {
		log.Println("Error creating live stream", err)
	}
	return stream
}

// GetLiveStream - get live streaming data response
func GetLiveStream(streamID string) muxgo.LiveStreamResponse {
	gs, err := muxClient.LiveStreamsApi.GetLiveStream(streamID)
	if err != nil {
		log.Println("Error getting live stream", err)
	}
	return gs
}