# pion-trickle-error


```
go test ./
```


```
$ go test ./...
--- FAIL: TestTrickle_enabled (0.05s)
    main_test.go:147: -- negotiate --
    main_test.go:86: pcAnswer.AddTransceiverFromKind()
    main_test.go:100: pcOffer.NewTrack()
    main_test.go:104: pcOffer.AddTrack()
    main_test.go:110: pcOffer.CreateOffer()
    main_test.go:114: pcOffer.SetLocalDescription()
    main_test.go:61: pcOffer.OnICEGatheringStateChange gathering
    main_test.go:116: pcAnswer.SetRemoteDescription()
    main_test.go:119: pcAnswer.CreateAnswer()
    main_test.go:47: pcOffer.OnICECandidate host 192.168.1.5:60658
    main_test.go:61: pcOffer.OnICEGatheringStateChange complete
    main_test.go:47: pcOffer.OnICECandidate <nil>
    main_test.go:79: pcAnswer.OnICEConnectionStateChange checking
    main_test.go:123: pcAnswer.SetLocalDescription()
    main_test.go:65: pcAnswer.OnICEGatheringStateChange gathering
    main_test.go:126: pcOffer.SetRemoteDescription()
    main_test.go:135: add pcOffer's ice candidate to pcAnswer: {candidate:foundation 1 udp 2130706431 192.168.1.5 60658 typ host <nil> 0xc00009e0b8 }
    main_test.go:72: pcOffer.OnICEConnectionStateChange checking
    main_test.go:54: pcAnswer.OnICECandidate host 192.168.1.5:58699
    main_test.go:65: pcAnswer.OnICEGatheringStateChange complete
    main_test.go:54: pcAnswer.OnICECandidate <nil>
    main_test.go:138: add pcAnswer's ice candidate to pcOffer: {candidate:foundation 1 udp 2130706431 192.168.1.5 58699 typ host <nil> 0xc00028c0d8 }
    main_test.go:79: pcAnswer.OnICEConnectionStateChange connected
    main_test.go:72: pcOffer.OnICEConnectionStateChange connected
    main_test.go:141: pcAnswer received track
    main_test.go:150: -- renegotiate --
    main_test.go:86: pcAnswer.AddTransceiverFromKind()
    main_test.go:100: pcOffer.NewTrack()
    main_test.go:104: pcOffer.AddTrack()
    main_test.go:110: pcOffer.CreateOffer()
    main_test.go:114: pcOffer.SetLocalDescription()
    main_test.go:61: pcOffer.OnICEGatheringStateChange gathering
    main_test.go:115:
                Error Trace:    main_test.go:115
                                                        main_test.go:151
                                                        main_test.go:155
                Error:          Received unexpected error:
                                attempting to gather candidates during gathering state
                Test:           TestTrickle_enabled
    main_test.go:116: pcAnswer.SetRemoteDescription()
    main_test.go:119: pcAnswer.CreateAnswer()
    main_test.go:123: pcAnswer.SetLocalDescription()
    main_test.go:65: pcAnswer.OnICEGatheringStateChange gathering
    main_test.go:126: pcOffer.SetRemoteDescription()
    main_test.go:141: pcAnswer received track
FAIL
FAIL    github.com/jeremija/pion-trickle-error  0.424s
FAIL
```
