package main_test

import (
	"context"
	"math/rand"
	"testing"
	"time"

	"github.com/pion/webrtc/v2"
	"github.com/pion/webrtc/v2/pkg/media"
	"github.com/stretchr/testify/assert"
)

func sendVideoUntilDone(done <-chan struct{}, t *testing.T, tracks []*webrtc.Track) {
	for {
		select {
		case <-time.After(20 * time.Millisecond):
			for _, track := range tracks {
				assert.NoError(t, track.WriteSample(media.Sample{Data: []byte{0x00}, Samples: 1}))
			}
		case <-done:
			return
		}
	}
}

func run(t *testing.T, trickle bool) {
	webrtcConfig := webrtc.Configuration{}
	settingEngine := webrtc.SettingEngine{}
	settingEngine.SetTrickle(trickle)
	mediaEngine := webrtc.MediaEngine{}
	mediaEngine.RegisterDefaultCodecs()
	api := webrtc.NewAPI(
		webrtc.WithMediaEngine(mediaEngine),
		webrtc.WithSettingEngine(settingEngine),
	)
	pcOffer, err := api.NewPeerConnection(webrtcConfig)
	assert.NoError(t, err)
	pcAnswer, err := api.NewPeerConnection(webrtcConfig)
	assert.NoError(t, err)

	c := make(chan struct{})
	offerCandidates := make(chan webrtc.ICECandidateInit, 10)
	answerCandidates := make(chan webrtc.ICECandidateInit, 10)

	pcOffer.OnICECandidate(func(c *webrtc.ICECandidate) {
		t.Log("pcOffer.OnICECandidate", c)
		if c != nil {
			offerCandidates <- c.ToJSON()
		}
	})

	pcAnswer.OnICECandidate(func(c *webrtc.ICECandidate) {
		t.Log("pcAnswer.OnICECandidate", c)
		if c != nil {
			answerCandidates <- c.ToJSON()
		}
	})

	pcOffer.OnICEGatheringStateChange(func(state webrtc.ICEGathererState) {
		t.Log("pcOffer.OnICEGatheringStateChange", state)
	})

	pcAnswer.OnICEGatheringStateChange(func(state webrtc.ICEGathererState) {
		t.Log("pcAnswer.OnICEGatheringStateChange", state)
	})

	_, err = pcOffer.CreateDataChannel("data", nil)
	assert.NoError(t, err)

	pcOffer.OnICEConnectionStateChange(func(state webrtc.ICEConnectionState) {
		t.Log("pcOffer.OnICEConnectionStateChange", state)
		if state == webrtc.ICEConnectionStateConnected {
			c <- struct{}{}
		}
	})

	pcAnswer.OnICEConnectionStateChange(func(state webrtc.ICEConnectionState) {
		t.Log("pcAnswer.OnICEConnectionStateChange", state)
		if state == webrtc.ICEConnectionStateConnected {
			c <- struct{}{}
		}
	})

	negotiate := func() {
		t.Log("pcAnswer.AddTransceiverFromKind()")
		_, err = pcAnswer.AddTransceiverFromKind(
			webrtc.RTPCodecTypeVideo,
			webrtc.RtpTransceiverInit{Direction: webrtc.RTPTransceiverDirectionRecvonly},
		)
		assert.NoError(t, err)

		tracksCh := make(chan *webrtc.Track)
		onTrackFired, onTrackFiredFunc := context.WithCancel(context.Background())
		pcAnswer.OnTrack(func(track *webrtc.Track, r *webrtc.RTPReceiver) {
			onTrackFiredFunc()
			tracksCh <- track
		})

		t.Log("pcOffer.NewTrack()")
		vp8Track, err := pcOffer.NewTrack(webrtc.DefaultPayloadTypeVP8, rand.Uint32(), "foo", "bar")
		assert.NoError(t, err)

		t.Log("pcOffer.AddTrack()")
		_, err = pcOffer.AddTrack(vp8Track)
		assert.NoError(t, err)

		go sendVideoUntilDone(onTrackFired.Done(), t, []*webrtc.Track{vp8Track})

		t.Log("pcOffer.CreateOffer()")
		offer, err := pcOffer.CreateOffer(nil)
		assert.NoError(t, err)

		t.Log("pcOffer.SetLocalDescription()")
		assert.NoError(t, pcOffer.SetLocalDescription(offer))
		t.Log("pcAnswer.SetRemoteDescription()")
		assert.NoError(t, pcAnswer.SetRemoteDescription(offer))

		t.Log("pcAnswer.CreateAnswer()")
		answer, err := pcAnswer.CreateAnswer(nil)
		assert.NoError(t, err)

		t.Log("pcAnswer.SetLocalDescription()")
		pcAnswer.SetLocalDescription(answer)

		t.Log("pcOffer.SetRemoteDescription()")
		assert.NoError(t, pcOffer.SetRemoteDescription(answer))

	loop:
		for {
			select {
			case <-time.After(3 * time.Second):
				t.Fatal("Timed out")
			case offerCandidate := <-offerCandidates:
				t.Log("add pcOffer's ice candidate to pcAnswer:", offerCandidate)
				pcAnswer.AddICECandidate(offerCandidate)
			case answerCandidate := <-answerCandidates:
				t.Log("add pcAnswer's ice candidate to pcOffer:", answerCandidate)
				pcOffer.AddICECandidate(answerCandidate)
			case <-tracksCh:
				t.Log("pcAnswer received track")
				break loop
			}
		}
	}

	t.Log("-- negotiate --")
	negotiate()

	t.Log("-- renegotiate --")
	negotiate()
}

func TestTrickle_enabled(t *testing.T) {
	run(t, true)
}

func TestTrickle_disabled(t *testing.T) {
	run(t, false)
}
