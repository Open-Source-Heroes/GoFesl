package theater

import (
	"github.com/ReviveNetwork/GoFesl/log"
	"github.com/SpencerSharkey/GoFesl/GameSpy"
)

// PENT - SERVER sent up when a player joins (entitle player?)
func (tM *TheaterManager) PENT(event GameSpy.EventClientFESLCommand) {
	if !event.Client.IsActive {
		return
	}

	pid := event.Command.Message["PID"]

	// Get 4 stats for PID
	rows, err := tM.getStatsStatement(4).Query(pid, "c_kit", "c_team", "elo", "level")
	if err != nil {
		log.Errorln("Failed gettings stats for hero "+pid, err.Error())
	}

	stats := make(map[string]string)

	for rows.Next() {
		var userID, heroID, heroName, statsKey, statsValue string
		err := rows.Scan(&userID, &heroID, &heroName, &statsKey, &statsValue)
		if err != nil {
			log.Errorln("Issue with database:", err.Error())
		}
		stats[statsKey] = statsValue
	}

	teamKey := "team_" + stats["c_team"]

	stmt, err := tM.db.Prepare("UPDATE games SET players_connected = players_connected + 1, players_joining = players_joining - 1,  " + teamKey + " = " + teamKey + " + 1, updated_at = NOW() WHERE gid = ? AND shard = ?")
	if err != nil {
		log.Panicln(err)
	}
	_, err = stmt.Exec(event.Command.Message["GID"], Shard)
	if err != nil {
		log.Panicln(err)
	}

	// This allows all right now, I think.
	answer := make(map[string]string)
	answer["TID"] = event.Command.Message["TID"]
	answer["PID"] = event.Command.Message["PID"]
	event.Client.WriteFESL("PENT", answer, 0x0)
}
