package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// Define Go struct to match the JSON structure
type MatchInfo struct {
	Meta struct {
		DataVersion string `json:"data_version"`
		Created     string `json:"created"`
		Revision    int    `json:"revision"`
	} `json:"meta"`
	Info struct {
		BallsPerOver int    `json:"balls_per_over"`
		City         string `json:"city"`
		Dates        []string
		Event        struct {
			MatchNumber int    `json:"match_number"`
			Name        string `json:"name"`
		} `json:"event"`
		Gender       string `json:"gender"`
		MatchType    string `json:"match_type"`
		MatchTypeNum int    `json:"match_type_number"`
		Officials    struct {
			MatchReferees  []string `json:"match_referees"`
			ReserveUmpires []string `json:"reserve_umpires"`
			TVUmpires      []string `json:"tv_umpires"`
			Umpires        []string `json:"umpires"`
		} `json:"officials"`
		Outcome struct {
			By struct {
				Wickets int `json:"wickets"`
			} `json:"by"`
			Winner string `json:"winner"`
		} `json:"outcome"`
		Overs         int `json:"overs"`
		PlayerOfMatch []string
		Players       struct {
			Australia []string
			Pakistan  []string
		} `json:"players"`
		Registry struct {
			People map[string]string
		} `json:"registry"`
		Info struct {
			// ... (other fields as before)
			Season interface{} `json:"season"`
		} `json:"info"`
		TeamType string `json:"team_type"`
		Teams    []string
		Toss     struct {
			Decision string `json:"decision"`
			Winner   string `json:"winner"`
		} `json:"toss"`
		Venue string `json:"venue"`
	} `json:"info"`
	Innings []struct {
		Team  string `json:"team"`
		Overs []struct {
			Over       int `json:"over"`
			Deliveries []struct {
				Batter string
				Bowler string
				Extras struct {
					Wides int
				} `json:"extras"`
				NonStriker string `json:"non_striker"`
				Runs       struct {
					Batter int
					Extras int
					Total  int
				} `json:"runs"`
			} `json:"deliveries"`
		} `json:"overs"`
		Powerplays []struct {
			From float64
			To   float64
			Type string
		} `json:"powerplays"`
	} `json:"innings"`
}

func main() {
	// Specify the directory where JSON files are located
	jsonDir := "/Users/ranga/Documents/fun_scorecard/odis_json"

	// Define filter conditions
	teamFilter := []string{"India", "Pakistan"}
	eventNameFilter := ""

	// Create and open the "stats.html" file
	htmlFile, err := os.Create("stats.html")
	if err != nil {
		fmt.Println("Error creating the HTML file:", err)
		return
	}
	defer htmlFile.Close()

	// Write the HTML structure to the file
	htmlFile.WriteString("<html><body>")
	htmlFile.WriteString("<table>")
	htmlFile.WriteString("<tr><th>Match Winner</th><th>City</th><th>Venue</th><th>Event Name</th></tr>") // Add more headers as needed

	// List JSON files in the directory
	files, err := filepath.Glob(filepath.Join(jsonDir, "*.json"))
	if err != nil {
		fmt.Println("Error listing JSON files:", err)
		return
	}

	// Iterate through the JSON files
	for _, file := range files {
		// Open the JSON file
		jsonFile, err := os.Open(file)
		if err != nil {
			fmt.Println("Error opening the file:", err)
			continue
		}
		defer jsonFile.Close()

		// Create a reader for the JSON file
		jsonReader := io.Reader(jsonFile)

		// Unmarshal the JSON data into a Go struct
		var matchInfo MatchInfo
		decoder := json.NewDecoder(jsonReader)
		err = decoder.Decode(&matchInfo)
		if err != nil {
			fmt.Println("Error unmarshaling JSON:", err)
			continue
		}

		// Apply filters
		teamFoundCount := 0
		for _, team := range teamFilter {
			if contains(matchInfo.Info.Teams, team) {
				teamFoundCount++
			}
		}

		if teamFoundCount == len(teamFilter) && (eventNameFilter == "" || matchInfo.Info.Event.Name == eventNameFilter) {
			// Print information in tabular format
			htmlFile.WriteString(fmt.Sprintf("<tr><td>%s</td><td>%s</td><td>%s</td><td>%s</td></tr>",
				matchInfo.Info.Outcome.Winner, matchInfo.Info.City, matchInfo.Info.Venue, matchInfo.Info.Event.Name))
		}

		// group by over per match
		groupedData := groupByOver(matchInfo)

		fmt.Println("Over\tTotal Runs")
		for over, totalRuns := range groupedData {
			fmt.Printf("%d\t%d\n", over, totalRuns)
		}
	}
	// Close the HTML table and body
	htmlFile.WriteString("</table>")
	htmlFile.WriteString("</body></html>")

}

// Function to group by overs and calculate the total runs per over
func groupByOver(match MatchInfo) map[int]int {
	groupedData := make(map[int]int)

	for _, inning := range match.Innings {
		for _, over := range inning.Overs {
			overNumber := over.Over
			totalRunsInOver := 0

			for _, delivery := range over.Deliveries {
				totalRunsInOver += delivery.Runs.Total
			}

			groupedData[overNumber] += totalRunsInOver
		}
	}

	return groupedData
}

// Helper function to check if a slice contains a string
func contains(slice []string, target string) bool {
	for _, value := range slice {
		if value == target {
			return true
		}
	}
	return false
}
