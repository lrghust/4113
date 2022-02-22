package paxos

import (
	"project/pkg/base"
)

// Fill in the function to lead the program to a state where A2 rejects the Accept Request of P1
func ToA2RejectP1() []func(s *base.State) bool {
	p1PreparePhase := func(s *base.State) bool {
		s1 := s.Nodes()["s1"].(*Server)
		return s1.proposer.Phase == Propose
	}

	p1AcceptPhase := func(s *base.State) bool {
		s1 := s.Nodes()["s1"].(*Server)
		return s1.proposer.Phase == Accept
	}

	p3PreparePhase := func(s *base.State) bool {
		s3 := s.Nodes()["s3"].(*Server)
		return s3.proposer.Phase == Propose
	}

	p3AcceptPhase := func(s *base.State) bool {
		s3 := s.Nodes()["s3"].(*Server)
		return s3.proposer.Phase == Accept
	}

	p1GetA2Reject := func(s *base.State) bool {
		s1 := s.Nodes()["s1"].(*Server)
		valid := false
		if s1.proposer.Phase == Accept && s1.proposer.ResponseCount == 2 && s1.proposer.SuccessCount == 1 {
			valid = true
		}
		return valid
	}

	return []func(s *base.State) bool{p1PreparePhase, p1AcceptPhase, p3PreparePhase, p3AcceptPhase, p1GetA2Reject}
}

// Fill in the function to lead the program to a state where a consensus is reached in Server 3.
func ToConsensusCase5() []func(s *base.State) bool {
	p3KnowConsensus := func(s *base.State) bool {
		s3 := s.Nodes()["s3"].(*Server)
		return s3.agreedValue == "v3"
	}
	return []func(s *base.State) bool{p3KnowConsensus}
}

// Fill in the function to lead the program to a state where all the Accept Requests of P1 are rejected
func NotTerminate1() []func(s *base.State) bool {
	p1PreparePhase := func(s *base.State) bool {
		s1 := s.Nodes()["s1"].(*Server)
		return s1.proposer.Phase == Propose
	}

	p1PreparePhase2 := func(s *base.State) bool {
		s1 := s.Nodes()["s1"].(*Server)
		return s1.proposer.Phase == Propose && s1.proposer.ResponseCount == 2 && s1.proposer.SuccessCount == 2
	}

	p1AcceptPhase := func(s *base.State) bool {
		s1 := s.Nodes()["s1"].(*Server)
		return s1.proposer.Phase == Accept && s1.proposer.ResponseCount == 0
	}

	p3PreparePhase := func(s *base.State) bool {
		s3 := s.Nodes()["s3"].(*Server)
		return s3.proposer.Phase == Propose
	}

	p3PreparePhase2 := func(s *base.State) bool {
		s3 := s.Nodes()["s3"].(*Server)
		return s3.proposer.Phase == Propose && s3.proposer.ResponseCount == 2 && s3.proposer.SuccessCount == 2
	}

	p3AcceptPhase := func(s *base.State) bool {
		s3 := s.Nodes()["s3"].(*Server)
		return s3.proposer.Phase == Accept && s3.proposer.ResponseCount == 0
	}

	p1GetReject := func(s *base.State) bool {
		s1 := s.Nodes()["s1"].(*Server)
		return s1.proposer.Phase == Accept && s1.proposer.ResponseCount == 1 && s1.proposer.SuccessCount == 0
	}

	return []func(s *base.State) bool{p1PreparePhase, p1PreparePhase2, p1AcceptPhase,
		p3PreparePhase, p3PreparePhase2, p3AcceptPhase, p1GetReject}
}

// Fill in the function to lead the program to a state where all the Accept Requests of P3 are rejected
func NotTerminate2() []func(s *base.State) bool {
	p1PreparePhase := func(s *base.State) bool {
		s1 := s.Nodes()["s1"].(*Server)
		return s1.proposer.Phase == Propose
	}

	p1PreparePhase2 := func(s *base.State) bool {
		s1 := s.Nodes()["s1"].(*Server)
		return s1.proposer.Phase == Propose && s1.proposer.ResponseCount == 2 && s1.proposer.SuccessCount == 2
	}

	p1AcceptPhase := func(s *base.State) bool {
		s1 := s.Nodes()["s1"].(*Server)
		return s1.proposer.Phase == Accept && s1.proposer.ResponseCount == 0
	}

	p3GetReject := func(s *base.State) bool {
		s3 := s.Nodes()["s3"].(*Server)
		return s3.proposer.Phase == Accept && s3.proposer.ResponseCount == 1 && s3.proposer.SuccessCount == 0
	}

	return []func(s *base.State) bool{p1PreparePhase, p1PreparePhase2, p1AcceptPhase, p3GetReject}
}

// Fill in the function to lead the program to a state where all the Accept Requests of P1 are rejected again.
func NotTerminate3() []func(s *base.State) bool {
	p3PreparePhase := func(s *base.State) bool {
		s3 := s.Nodes()["s3"].(*Server)
		return s3.proposer.Phase == Propose
	}

	p3PreparePhase2 := func(s *base.State) bool {
		s3 := s.Nodes()["s3"].(*Server)
		return s3.proposer.Phase == Propose && s3.proposer.ResponseCount == 2 && s3.proposer.SuccessCount == 2
	}

	p3AcceptPhase := func(s *base.State) bool {
		s3 := s.Nodes()["s3"].(*Server)
		return s3.proposer.Phase == Accept && s3.proposer.ResponseCount == 0
	}

	p1GetReject := func(s *base.State) bool {
		s1 := s.Nodes()["s1"].(*Server)
		return s1.proposer.Phase == Accept && s1.proposer.ResponseCount == 1 && s1.proposer.SuccessCount == 0
	}

	return []func(s *base.State) bool{p3PreparePhase, p3PreparePhase2, p3AcceptPhase, p1GetReject}
}

// Fill in the function to lead the program to make P1 propose first, then P3 proposes, but P1 get rejects in
// Accept phase
func concurrentProposer1() []func(s *base.State) bool {
	p1PreparePhase := func(s *base.State) bool {
		s1 := s.Nodes()["s1"].(*Server)
		return s1.proposer.Phase == Propose
	}

	p1PreparePhase2 := func(s *base.State) bool {
		s1 := s.Nodes()["s1"].(*Server)
		return s1.proposer.Phase == Propose && s1.proposer.ResponseCount == 2 && s1.proposer.SuccessCount == 2
	}

	p1AcceptPhase := func(s *base.State) bool {
		s1 := s.Nodes()["s1"].(*Server)
		return s1.proposer.Phase == Accept && s1.proposer.ResponseCount == 0
	}

	p3PreparePhase := func(s *base.State) bool {
		s3 := s.Nodes()["s3"].(*Server)
		return s3.proposer.Phase == Propose
	}

	p3PreparePhase2 := func(s *base.State) bool {
		s3 := s.Nodes()["s3"].(*Server)
		return s3.proposer.Phase == Propose && s3.proposer.ResponseCount == 2 && s3.proposer.SuccessCount == 2
	}

	p3AcceptPhase := func(s *base.State) bool {
		s3 := s.Nodes()["s3"].(*Server)
		return s3.proposer.Phase == Accept && s3.proposer.ResponseCount == 0
	}

	p1GetReject := func(s *base.State) bool {
		s1 := s.Nodes()["s1"].(*Server)
		return s1.proposer.Phase == Accept && s1.proposer.ResponseCount == 1 && s1.proposer.SuccessCount == 0
	}

	return []func(s *base.State) bool{p1PreparePhase, p1PreparePhase2, p1AcceptPhase,
		p3PreparePhase, p3PreparePhase2, p3AcceptPhase, p1GetReject}
}

// Fill in the function to lead the program continue  P3's proposal  and reaches consensus at the value of "v3".
func concurrentProposer2() []func(s *base.State) bool {
	p3GetOK := func(s *base.State) bool {
		s3 := s.Nodes()["s3"].(*Server)
		return s3.proposer.Phase == Accept && s3.proposer.ResponseCount == 1 && s3.proposer.SuccessCount == 1
	}

	p3KnowConsensus := func(s *base.State) bool {
		s3 := s.Nodes()["s3"].(*Server)
		return s3.agreedValue == "v3"
	}

	return []func(s *base.State) bool{p3GetOK, p3KnowConsensus}
}
