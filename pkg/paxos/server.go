package paxos

import (
	"project/pkg/base"
)

const (
	Propose = "propose"
	Accept  = "accept"
	Decide  = "decide"
)

type Proposer struct {
	N             int
	Phase         string
	N_a_max       int
	V             interface{}
	SuccessCount  int
	ResponseCount int
	// To indicate if response from peer is received, should be initialized as []bool of len(server.peers)
	Responses []bool
	// Use this field to check if a message is latest.
	SessionId int

	// in case node will propose again - restore initial value
	InitialValue interface{}
}

type ServerAttribute struct {
	peers []base.Address
	me    int

	// Paxos parameter
	n_p int
	n_a int
	v_a interface{}

	// final result
	agreedValue interface{}

	// Propose parameter
	proposer Proposer

	// retry
	timeout *TimeoutTimer
}

type Server struct {
	base.CoreNode
	ServerAttribute
}

func NewServer(peers []base.Address, me int, proposedValue interface{}) *Server {
	response := make([]bool, len(peers))
	return &Server{
		CoreNode: base.CoreNode{},
		ServerAttribute: ServerAttribute{
			peers: peers,
			me:    me,
			proposer: Proposer{
				InitialValue: proposedValue,
				Responses:    response,
			},
			timeout: &TimeoutTimer{},
		},
	}
}

func (server *Server) MessageHandler(message base.Message) []base.Node {
	//TODO: implement it
	nodes := make([]base.Node, 0)
	switch message.(type) {
	case *ProposeRequest:
		propose := message.(*ProposeRequest)
		newNode := server.copy()
		if propose.N > newNode.n_p {
			newNode.n_p = propose.N
			newNode.Response = []base.Message{
				&ProposeResponse{
					CoreMessage: base.MakeCoreMessage(propose.To(), propose.From()),
					Ok:          true,
					N_p:         propose.N,
					N_a:		 newNode.n_a,
					V_a:		 newNode.v_a,
					SessionId:   propose.SessionId,
				},
			}
		} else {
			newNode.Response = []base.Message{
				&ProposeResponse{
					CoreMessage: base.MakeCoreMessage(propose.To(), propose.From()),
					Ok:          false,
					N_p:         newNode.n_p,
					SessionId:   propose.SessionId,
				},
			}
		}
		nodes = append(nodes, newNode)

	case *ProposeResponse:
		response := message.(*ProposeResponse)
		index := server.Index(response.From())
		if !server.proposer.Responses[index] {
			if response.Ok && server.proposer.SuccessCount + 1 >= len(server.peers) / 2 + 1 {
				// node1: accept
				newNode := server.copy()
				newNode.proposer.Phase = Accept
				newNode.proposer.ResponseCount = 0
				newNode.proposer.SuccessCount = 0
				newNode.proposer.Responses = make([]bool, len(newNode.peers))
				if response.N_a > newNode.proposer.N_a_max {
					newNode.proposer.N_a_max = response.N_a
					newNode.proposer.V = response.V_a
				}
				responses := make([]base.Message, 0)
				for _, peer := range server.peers {
					propose := &AcceptRequest{
						CoreMessage: base.MakeCoreMessage(newNode.Address(), peer),
						N:           response.N_p,
						V:           newNode.proposer.V,
						SessionId:   response.SessionId,
					}
					responses = append(responses, propose)
				}
				newNode.Response = responses
				nodes = append(nodes, newNode)
			}
			if true || server.proposer.ResponseCount + 1 < len(server.peers) {
				// node2: continue
				newNode := server.copy()
				newNode.proposer.ResponseCount++
				newNode.proposer.Responses[index] = true
				if response.Ok {
					newNode.proposer.SuccessCount++
					if response.N_a > newNode.proposer.N_a_max {
						newNode.proposer.N_a_max = response.N_a
						newNode.proposer.V = response.V_a
					}
				}
				nodes = append(nodes, newNode)
			}
		} else {
			newNode := server.copy()
			nodes = append(nodes, newNode)
		}

	case *AcceptRequest:
		accept := message.(*AcceptRequest)
		newNode := server.copy()
		if accept.N >= newNode.n_p {
			newNode.n_p = accept.N
			newNode.n_a	= accept.N
			newNode.v_a = accept.V
			newNode.Response = []base.Message{
				&AcceptResponse{
					CoreMessage: base.MakeCoreMessage(accept.To(), accept.From()),
					Ok:          true,
					N_p:         accept.N,
					SessionId:   accept.SessionId,
				},
			}
		} else {
			newNode.Response = []base.Message{
				&AcceptResponse{
					CoreMessage: base.MakeCoreMessage(accept.To(), accept.From()),
					Ok:          false,
					N_p:         newNode.n_p,
					SessionId:   accept.SessionId,
				},
			}
		}
		nodes = append(nodes, newNode)

	case *AcceptResponse:
		response := message.(*AcceptResponse)
		index := server.Index(response.From())
		if !server.proposer.Responses[index] {
			if response.Ok && server.proposer.SuccessCount + 1 >= len(server.peers) / 2 + 1 {
				// node1: decide
				newNode := server.copy()
				newNode.proposer.Phase = Decide
				//newNode.timeout = nil
				newNode.agreedValue = newNode.proposer.V
				newNode.proposer.ResponseCount = 0
				newNode.proposer.SuccessCount = 0
				newNode.proposer.Responses = make([]bool, len(newNode.peers))
				responses := make([]base.Message, 0)
				for _, peer := range server.peers {
					propose := &DecideRequest{
						CoreMessage: base.MakeCoreMessage(newNode.Address(), peer),
						V:           newNode.proposer.V,
						SessionId:   response.SessionId,
					}
					responses = append(responses, propose)
				}
				newNode.Response = responses
				nodes = append(nodes, newNode)
			}
			if true || server.proposer.ResponseCount + 1 < len(server.peers) {
				// node2: continue
				newNode := server.copy()
				newNode.proposer.ResponseCount++
				newNode.proposer.Responses[index] = true
				if response.Ok {
					newNode.proposer.SuccessCount++
				}
				nodes = append(nodes, newNode)
			}
		} else {
			newNode := server.copy()
			nodes = append(nodes, newNode)
		}

	case *DecideRequest:
		decide := message.(*DecideRequest)
		newNode := server.copy()
		newNode.proposer.Phase = Decide
		//newNode.timeout = nil
		newNode.agreedValue = decide.V
		nodes = append(nodes, newNode)
	}
	return nodes
}

// To start a new round of Paxos.
func (server *Server) StartPropose() {
	//TODO: implement it
	n := server.proposer.N + 1
	session := server.proposer.SessionId + 1
	responses := make([]base.Message, 0)
	for _, peer := range server.peers {
		propose := &ProposeRequest{
			CoreMessage: base.MakeCoreMessage(server.Address(), peer),
			N:           n,
			SessionId:   session,
		}
		responses = append(responses, propose)
	}
	server.Response = responses
	server.proposer.Phase = Propose
	server.proposer.N = n
	server.proposer.V = server.proposer.InitialValue
	server.proposer.SessionId = session
	server.proposer.ResponseCount = 0
	server.proposer.SuccessCount = 0
	server.proposer.Responses = make([]bool, len(server.peers))
}

// Returns a deep copy of server node
func (server *Server) copy() *Server {
	response := make([]bool, len(server.peers))
	for i, flag := range server.proposer.Responses {
		response[i] = flag
	}

	var copyServer Server
	copyServer.me = server.me
	// shallow copy is enough, assuming it won't change
	copyServer.peers = server.peers
	copyServer.n_a = server.n_a
	copyServer.n_p = server.n_p
	copyServer.v_a = server.v_a
	copyServer.agreedValue = server.agreedValue
	copyServer.proposer = Proposer{
		N:             server.proposer.N,
		Phase:         server.proposer.Phase,
		N_a_max:       server.proposer.N_a_max,
		V:             server.proposer.V,
		SuccessCount:  server.proposer.SuccessCount,
		ResponseCount: server.proposer.ResponseCount,
		Responses:     response,
		InitialValue:  server.proposer.InitialValue,
		SessionId:     server.proposer.SessionId,
	}

	// doesn't matter, timeout timer is state-less
	copyServer.timeout = server.timeout

	return &copyServer
}

func (server *Server) NextTimer() base.Timer {
	return server.timeout
}

// A TimeoutTimer tick simulates the situation where a proposal procedure times out.
// It will close the current Paxos round and start a new one if no consensus reached so far,
// i.e. the server after timer tick will reset and restart from the first phase if Paxos not decided.
// The timer will not be activated if an agreed value is set.
func (server *Server) TriggerTimer() []base.Node {
	if server.timeout == nil {
		return nil
	}

	subNode := server.copy()
	subNode.StartPropose()

	return []base.Node{subNode}
}

func (server *Server) Attribute() interface{} {
	return server.ServerAttribute
}

func (server *Server) Copy() base.Node {
	return server.copy()
}

func (server *Server) Hash() uint64 {
	return base.Hash("paxos", server.ServerAttribute)
}

func (server *Server) Equals(other base.Node) bool {
	otherServer, ok := other.(*Server)

	if !ok || server.me != otherServer.me ||
		server.n_p != otherServer.n_p || server.n_a != otherServer.n_a || server.v_a != otherServer.v_a ||
		(server.timeout == nil) != (otherServer.timeout == nil) {
		return false
	}

	if server.proposer.N != otherServer.proposer.N || server.proposer.V != otherServer.proposer.V ||
		server.proposer.N_a_max != otherServer.proposer.N_a_max || server.proposer.Phase != otherServer.proposer.Phase ||
		server.proposer.InitialValue != otherServer.proposer.InitialValue ||
		server.proposer.SuccessCount != otherServer.proposer.SuccessCount ||
		server.proposer.ResponseCount != otherServer.proposer.ResponseCount {
		return false
	}

	for i, response := range server.proposer.Responses {
		if response != otherServer.proposer.Responses[i] {
			return false
		}
	}

	return true
}

func (server *Server) Address() base.Address {
	return server.peers[server.me]
}

func (server *Server) Index(address base.Address) int {
	for i, addr := range server.peers {
		if addr == address {
			return i
		}
	}
	return -1
}