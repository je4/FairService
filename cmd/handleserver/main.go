package main

import (
	"fmt"
	"github.com/je4/FairService/v2/pkg/handle_net"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	srv, err := handle_net.NewHandleServer([]string{"handle"})
	if err != nil {
		panic(err)
	}
	srv.RegisterFunction(handle_net.OpCodeOCGetSiteInfo, func(msg *handle_net.Message) (*handle_net.Message, error) {
		var result = &handle_net.Message{
			Envelope: handle_net.MessageEnvelope{
				MajorVersion: 0,
				MinorVersion: 0,
				MessageFlag: handle_net.MessageFlag{
					CP: false,
					EC: false,
					TC: false,
				},
				SessionId:      msg.Envelope.SessionId,
				RequestId:      msg.Envelope.RequestId,
				SequenceNumber: 0,
				MessageLength:  0,
			},
			Header:           handle_net.MessageHeader{},
			Body:             handle_net.MessageBody{},
			CredentialLength: 0,
			Credential:       handle_net.Credential{},
		}
		fmt.Printf("received message: %v", msg)
		return result, nil
	})
	if err := srv.Start("tcp4", ":2641"); err != nil {
		panic(err)
	}
	end := make(chan bool, 1)

	// process waiting for interrupt signal (TERM or KILL)
	go func() {
		sigint := make(chan os.Signal, 1)

		// interrupt signal sent from terminal
		signal.Notify(sigint, os.Interrupt)

		signal.Notify(sigint, syscall.SIGTERM)
		signal.Notify(sigint, syscall.SIGKILL)

		<-sigint

		// We received an interrupt signal, shut down.
		fmt.Sprintf("shutdown requested")
		srv.Stop()

		end <- true
	}()

	<-end
	fmt.Sprintf("server stopped")

}
