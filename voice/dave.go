package voice

import (
	"log/slog"

	"github.com/disgoorg/snowflake/v2"
)

type (
	Dave interface {
		CreateSession(userID snowflake.ID, callbacks Callbacks) DaveSession
	}

	Callbacks struct {
		SendMLSKeyPackage        func(mlsKeyPackage []byte) error
		SendMLSCommitWelcome     func(mlsCommitWelcome []byte) error
		SendReadyForTransition   func(transitionID uint16) error
		SendInvalidCommitWelcome func(transitionID uint16) error
	}

	DaveSession interface {
		// MaxSupportedProtocolVersion returns the maximum supported DAVE version for this session
		MaxSupportedProtocolVersion() int

		SetChannelID(channelID snowflake.ID)

		// AssignSsrcToCodec maps a given SSRC to a specific Codec
		// FIXME: Consider having a global Codec
		AssignSsrcToCodec(codec int, ssrc uint32)

		// EncryptOpus encrypts an OPUS frame
		EncryptOpus(ssrc uint32, frame []byte) ([]byte, error)

		// DecryptOpus decrypts an OPUS frame
		DecryptOpus(userID snowflake.ID, frame []byte) ([]byte, error)

		// AddUser adds a user to the MLS group
		AddUser(userID snowflake.ID)

		// RemoveUser removes a use from the MLS group
		RemoveUser(userID snowflake.ID)

		// OnSelectProtocolAck is to be called when SELECT_PROTOCOL_ACK (4) is received
		OnSelectProtocolAck(protocolVersion uint16)

		// OnDavePrepareTransition is to be called when DAVE_PROTOCOL_PREPARE_TRANSITION (21) is received
		OnDavePrepareTransition(transitionID uint16, protocolVersion uint16)

		// OnDaveExecuteTransition is to be called when DAVE_PROTOCOL_EXECUTE_TRANSITION (22) is received
		OnDaveExecuteTransition(protocolVersion uint16)

		// OnDavePrepareEpoch is to be called when DAVE_PROTOCOL_PREPARE_EPOCH (24) is received
		OnDavePrepareEpoch(epoch int, protocolVersion uint16)

		// OnDaveMLSExternalSenderPackage is to be called when DAVE_MLS_EXTERNAL_SENDER_PACKAGE (25) is received
		OnDaveMLSExternalSenderPackage(externalSenderPackage []byte)

		// OnDaveMLSProposals is to be called when DAVE_MLS_PROPOSALS (27) is received
		OnDaveMLSProposals(proposals []byte)

		// OnDaveMLSPrepareCommitTransition is to be called when DAVE_MLS_ANNOUNCE_COMMIT_TRANSITION (29) is received
		OnDaveMLSPrepareCommitTransition(transitionID uint16, commitMessage []byte)

		// OnDaveMLSWelcome is to be called when DAVE_MLS_WELCOME (30) is received
		OnDaveMLSWelcome(transitionID uint16, welcomeMessage []byte)
	}
)

func NewNoopDave() Dave {
	return &noOpDave{}
}

type noOpDave struct{}

func (n *noOpDave) CreateSession(userID snowflake.ID, callbacks Callbacks) DaveSession {
	slog.Warn("Using passthrough dave session. Please migrate to an implementation of libdave or your audio connections will stop working on 01.03.2026")

	return &noOpDaveSession{}
}

type noOpDaveSession struct{}

func (n *noOpDaveSession) MaxSupportedProtocolVersion() int {
	return 0
}

func (n *noOpDaveSession) EncryptOpus(ssrc uint32, frame []byte) ([]byte, error) {
	return frame, nil
}

func (n *noOpDaveSession) DecryptOpus(userID snowflake.ID, frame []byte) ([]byte, error) {
	return frame, nil
}

func (n *noOpDaveSession) SetChannelID(channelID snowflake.ID)                                 {}
func (n *noOpDaveSession) AssignSsrcToCodec(codec int, ssrc uint32)                            {}
func (n *noOpDaveSession) AddUser(userID snowflake.ID)                                         {}
func (n *noOpDaveSession) RemoveUser(userID snowflake.ID)                                      {}
func (n *noOpDaveSession) OnSelectProtocolAck(protocolVersion uint16)                          {}
func (n *noOpDaveSession) OnDavePrepareTransition(transitionID uint16, protocolVersion uint16) {}
func (n *noOpDaveSession) OnDaveExecuteTransition(protocolVersion uint16)                      {}
func (n *noOpDaveSession) OnDavePrepareEpoch(epoch int, protocolVersion uint16)                {}
func (n *noOpDaveSession) OnDaveMLSExternalSenderPackage(externalSenderPackage []byte)         {}
func (n *noOpDaveSession) OnDaveMLSProposals(proposals []byte)                                 {}
func (n *noOpDaveSession) OnDaveMLSPrepareCommitTransition(transitionID uint16, commit []byte) {}
func (n *noOpDaveSession) OnDaveMLSWelcome(transitionID uint16, welcomeMessage []byte)         {}
