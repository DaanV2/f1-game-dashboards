package game

import (
	"errors"
	"fmt"
	"net"
	"strings"

	"github.com/DaanV2/f1-game-dashboards/server/pkg/hooks"
	"github.com/DaanV2/f1-game-dashboards/server/sessions"
	"github.com/DaanV2/go-f1-library/encoding"
	"github.com/DaanV2/go-f1-library/enums"
	f1_2023 "github.com/DaanV2/go-f1-library/packets/2023"
	"github.com/DaanV2/go-f1-library/packets/general"
	"github.com/charmbracelet/log"
)

const (
	max_packet_size = f1_2023.MAX_PACKET_SIZE
	min_packet_size = f1_2023.PacketHeaderSize
)

type PacketProcessor struct {
	options packetProcessorOptions

	parser   *f1_2023.PacketParser
	pipeline *PacketPipeline

	chairs map[int]*chairSession
}

type chairSession struct {
	chair *sessions.Chair
	conn  *net.UDPConn
}

type chairProcessor struct {
	session   *chairSession
	processor *PacketProcessor
}

func NewPacketProcessor(chairs *sessions.ChairManager, options ...PacketOption) *PacketProcessor {
	opts := packetProcessorOptions{}
	opts._default()
	opts.apply(options...)

	processor := &PacketProcessor{
		chairs:   make(map[int]*chairSession),
		options:  opts,
		parser:   f1_2023.NewPacketParser(),
		pipeline: NewPacketPipeline(),
	}

	chairs.OnChairAdded.Add(processor.HandleChairAdded)
	chairs.OnChairRemoved.Add(processor.HandleChairRemoved)

	for _, chair := range chairs.Chairs() {
		processor.HandleChairAdded(chair)
	}

	return processor
}

func (pp *PacketProcessor) Close() {
	for _, session := range pp.chairs {
		if err := session.Stop(); err != nil {
			log.Error("error stopping session", "error", err, "port", session.chair.Port)
		}
	}
}

func (pp *PacketProcessor) HandleChairAdded(chair *sessions.Chair) {
	logger := log.With("port", chair.Port)

	logger.Info("starting server on chair...")
	session, ok := pp.chairs[chair.Port]
	// If sessions already exists and the connection is still open, return
	if session != nil && ok && session.conn != nil {
		logger.Info("skipping start server on chair, already exists and seems healthy")
		return
	}

	// If session is nil, create a new session
	if session == nil {
		session = &chairSession{
			chair: chair,
		}
		pp.chairs[chair.Port] = session
	}

	// If crashed or closed, ensure the connection is nil
	if session.conn == nil {
		data := &chairProcessor{
			session:   session,
			processor: pp,
		}
		go data.Start()
	}
}

func (pp *PacketProcessor) HandleChairRemoved(chair *sessions.Chair) {
	logger := log.With("port", chair.Port)
	logger.Info("closing server on chair...")

	session, ok := pp.chairs[chair.Port]
	if !ok {
		return
	}
	defer func() {
		delete(pp.chairs, chair.Port)
		logger.Info("closed server on chair")
	}()

	err := session.conn.Close()
	if strings.Contains(err.Error(), "use of closed network connection") {
		return
	}
	if err != nil {
		logger.Error("error closing connection", "error", err)
	}
}

func (cs *chairSession) Stop() error {
	defer func() {
		cs.conn = nil
	}()
	if cs.conn == nil {
		return nil
	}

	err := cs.conn.Close()
	// Conn closed
	if strings.Contains(err.Error(), "use of closed network connection") {
		return nil
	}

	return err
}

func (cp *chairProcessor) Start() (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic: %v", r)
		}
	}()

	var (
		buf     [max_packet_size]byte
		n       int
		address *net.UDPAddr
	)
	logger := log.With("port", cp.session.chair.Port, "name", cp.session.chair.Name)
	udpAddr, err := net.ResolveUDPAddr("udp", cp.session.chair.Name)
	if err != nil {
		return err
	}
	conn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		panic(err)
	}
	cp.session.conn = conn

	// TODO add signal handling for when the server is shutting down

	for {
		n, address, err = cp.session.conn.ReadFromUDP(buf[0:])
		if err != nil {
			// Conn closed
			if strings.Contains(err.Error(), "use of closed network connection") {
				return
			}

			logger.Error("error reading from udp", "error", err)

		// If the chair is not active, skip the packet
		} else if cp.session.chair.Active {
			err := cp.handlePacket(buf[:n])
			if err != nil {
				logger.Error("error handling packet", "error", err, "ip", address.IP, "port", address.Port)
			}
		}
	}
}

func (cp *chairProcessor) handlePacket(packet []byte) error {
	//NOTE: packet is owned by the caller, so we need to copy it or process it immediately
	if len(packet) < min_packet_size {
		return errors.New("packet too small")
	}

	header := general.ParsePacketHeader(packet)
	switch header.PacketFormat {
	case enums.PF_F1_2023:
		return cp.handle2023Packet(packet)
	}

	return fmt.Errorf("unknown packet format: %d", header.PacketFormat)
}

func (cp *chairProcessor) handle2023Packet(packet []byte) error {
	pipeline := cp.processor.pipeline
	parser := cp.processor.parser
	decoder := encoding.NewDecoder(packet)
	header, err := parser.PacketHeader(decoder)
	if err != nil {
		return err
	}

	switch header.PacketId {
	case enums.PID_Motion:
		return process(cp, decoder, header, pipeline.Motion, parser.PacketMotionData)
	case enums.PID_Session:
		return process(cp, decoder, header, pipeline.Session, parser.PacketSessionData)
	case enums.PID_LapData:
		return process(cp, decoder, header, pipeline.LapData, parser.PacketLapData)
	case enums.PID_Event:
		return process(cp, decoder, header, pipeline.Event, parser.PacketEventData)
	case enums.PID_Participants:
		return process(cp, decoder, header, pipeline.Participants, parser.PacketParticipantsData)
	case enums.PID_CarSetups:
		return process(cp, decoder, header, pipeline.CarSetups, parser.PacketCarSetupData)
	case enums.PID_CarTelemetry:
		return process(cp, decoder, header, pipeline.CarTelemetry, parser.PacketCarTelemetryData)
	case enums.PID_CarStatus:
		return process(cp, decoder, header, pipeline.CarStatus, parser.PacketCarStatusData)
	case enums.PID_FinalClassification:
		return process(cp, decoder, header, pipeline.FinalClassification, parser.PacketFinalClassificationData)
	case enums.PID_LobbyInfo:
		return process(cp, decoder, header, pipeline.LobbyInfo, parser.PacketLobbyInfoData)
	case enums.PID_CarDamage:
		return process(cp, decoder, header, pipeline.CarDamage, parser.PacketCarDamageData)
	case enums.PID_SessionHistory:
		return process(cp, decoder, header, pipeline.SessionHistory, parser.PacketSessionHistoryData)
	case enums.PID_TyreSets:
		return process(cp, decoder, header, pipeline.TyreSets, parser.PacketTyreSetsData)
	case enums.PID_MotionEx:
		return process(cp, decoder, header, pipeline.MotionEx, parser.PacketMotionExData)
	}

	return fmt.Errorf("unknown packet id: %d", header.PacketId)
}

func process[T any](cp *chairProcessor, decoder *encoding.Decoder, header f1_2023.PacketHeader, hook hooks.Hook[PacketWithChair[T]], get func(decoder *encoding.Decoder, header f1_2023.PacketHeader) (T, error)) error {
	if !hook.Active() {
		return nil
	}
	packet, err := get(decoder, header)
	if err != nil {
		return err
	}

	data := PacketWithChair[T]{
		Chair:  cp.session.chair,
		Packet: packet,
	}

	hook.Call(data)

	return nil
}