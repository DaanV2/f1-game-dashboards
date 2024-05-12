package game

import (
	"github.com/DaanV2/f1-game-dashboards/server/pkg/hooks"
	"github.com/DaanV2/f1-game-dashboards/server/sessions"
	f1_2023 "github.com/DaanV2/go-f1-library/packets/2023"
)

type PacketWithChair[T any] struct {
	Chair  *sessions.Chair // The chair the packet came from
	Packet T               // The packet
}

type PacketPipeline struct {
	Motion              hooks.Hook[PacketWithChair[f1_2023.PacketMotionData]]
	Session             hooks.Hook[PacketWithChair[f1_2023.PacketSessionData]]
	LapData             hooks.Hook[PacketWithChair[f1_2023.PacketLapData]]
	Event               hooks.Hook[PacketWithChair[f1_2023.PacketEventData]]
	Participants        hooks.Hook[PacketWithChair[f1_2023.PacketParticipantsData]]
	CarSetups           hooks.Hook[PacketWithChair[f1_2023.PacketCarSetupsData]]
	CarTelemetry        hooks.Hook[PacketWithChair[f1_2023.PacketCarTelemetryData]]
	CarStatus           hooks.Hook[PacketWithChair[f1_2023.PacketCarStatusData]]
	FinalClassification hooks.Hook[PacketWithChair[f1_2023.PacketFinalClassificationData]]
	LobbyInfo           hooks.Hook[PacketWithChair[f1_2023.PacketLobbyInfoData]]
	CarDamage           hooks.Hook[PacketWithChair[f1_2023.PacketCarDamageData]]
	SessionHistory      hooks.Hook[PacketWithChair[f1_2023.PacketSessionHistoryData]]
	TyreSets            hooks.Hook[PacketWithChair[f1_2023.PacketTyreSetsData]]
	MotionEx            hooks.Hook[PacketWithChair[f1_2023.PacketMotionExData]]
}

func NewPacketPipeline() *PacketPipeline {
	return &PacketPipeline{
		Motion:              hooks.NewHook[PacketWithChair[f1_2023.PacketMotionData]](),
		Session:             hooks.NewHook[PacketWithChair[f1_2023.PacketSessionData]](),
		LapData:             hooks.NewHook[PacketWithChair[f1_2023.PacketLapData]](),
		Event:               hooks.NewHook[PacketWithChair[f1_2023.PacketEventData]](),
		Participants:        hooks.NewHook[PacketWithChair[f1_2023.PacketParticipantsData]](),
		CarSetups:           hooks.NewHook[PacketWithChair[f1_2023.PacketCarSetupsData]](),
		CarTelemetry:        hooks.NewHook[PacketWithChair[f1_2023.PacketCarTelemetryData]](),
		CarStatus:           hooks.NewHook[PacketWithChair[f1_2023.PacketCarStatusData]](),
		FinalClassification: hooks.NewHook[PacketWithChair[f1_2023.PacketFinalClassificationData]](),
		LobbyInfo:           hooks.NewHook[PacketWithChair[f1_2023.PacketLobbyInfoData]](),
		CarDamage:           hooks.NewHook[PacketWithChair[f1_2023.PacketCarDamageData]](),
		SessionHistory:      hooks.NewHook[PacketWithChair[f1_2023.PacketSessionHistoryData]](),
		TyreSets:            hooks.NewHook[PacketWithChair[f1_2023.PacketTyreSetsData]](),
		MotionEx:            hooks.NewHook[PacketWithChair[f1_2023.PacketMotionExData]](),
	}
}
