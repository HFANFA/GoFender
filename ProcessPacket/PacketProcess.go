package ProcessPacket

import (
	"GoFender/Database"
	"GoFender/MachineLearning"
	"GoFender/SuricataMatch"
	"GoFender/Utils"
)

func PacketProcess(packet Utils.CommonPacket) {
	// malicious traffic detection and storage
	isMalicious(packet)
}

func isMalicious(packet Utils.CommonPacket) {
	matchscore, in, matchtype := SuricataMatch.PacketMatch(&packet, SuricataMatch.RuleSet)
	mlscore, mltype := MachineLearning.Attacktype(packet)
	if matchscore+mlscore >= 79.99 {
		if in {
			Packetinfo := Utils.EvilPacket{
				CommInfo:   packet,
				Type:       "Evil",
				AttackType: matchtype[0],
			}
			Database.Insert(Packetinfo)
		} else {
			Packetinfo := Utils.EvilPacket{
				CommInfo:   packet,
				Type:       "Evil",
				AttackType: mltype,
			}
			Database.Insert(Packetinfo)
		}
	} else {
		Packetinfo := Utils.NomPacket{
			CommInfo: packet,
			Type:     "Normal",
		}
		Database.Insert(Packetinfo)
	}
}
