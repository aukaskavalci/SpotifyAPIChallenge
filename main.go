package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/zmb3/spotify"
)

type trackRecord struct {
	frequency int
	name      string
}

var (
	autherizatetionHeader = "Get token from https://developer.spotify.com/console/get-playlists/?user_id=gulsahguray&limit=&offset= "
	desiredPlaylistsURL   = "https://api.spotify.com/v1/users/gulsahguray/playlists"
)

func addHeaders(request *http.Request) {
	request.Header.Add("Accept", "application/json")
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("Authorization", autherizatetionHeader)
}

func getRequestResponse(url string) ([]byte, error) {
	request, err := http.NewRequest("GET", url, nil)
	addHeaders(request)
	client := http.DefaultClient
	response, err := client.Do(request)
	if err != nil {
		fmt.Println("Error happened while getting response from first playlists request", err.Error())
		return nil, err
	}
	responseStatus := response.StatusCode
	for responseStatus != http.StatusOK {
		fmt.Printf("Received status code %d, will try in a sec\n", response.StatusCode)
		time.Sleep(1 * time.Second)
		response, err := client.Do(request)
		if err != nil {
			fmt.Println("Error happened while getting response from request", err.Error())
			return nil, err
		}
		responseStatus = response.StatusCode
	}

	prettyResponse, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Error happened while reading the response body", err.Error())
		return nil, err
	}
	return prettyResponse, nil
}

func main() {
	var playlistList spotify.SimplePlaylistPage
	var playlistTracks spotify.PlaylistTrackPage
	trackFrequency := make(map[string]trackRecord)

	playlistsResponse, err := getRequestResponse(desiredPlaylistsURL)
	err = json.Unmarshal(playlistsResponse, &playlistList)
	if err != nil {
		fmt.Printf("Error happened while unmarshalling playlist response: %s\n", err.Error())
		return
	}

	nextpage := "start"
	for nextpage != playlistList.Next || nextpage == "start" {

		for _, plist := range playlistList.Playlists {
			tracksResponse, err := getRequestResponse(plist.Tracks.Endpoint)
			err = json.Unmarshal(tracksResponse, &playlistTracks)
			if err != nil {
				fmt.Printf("Error happened while unmarshalling track response: %s\n", err.Error())
				return
			}

			nextpageTracks := "startTrack"
			for nextpageTracks != playlistTracks.Next || nextpageTracks == "startTrack" {
				for _, track := range playlistTracks.Tracks {

					if counter, ok := trackFrequency[string(track.Track.ID)]; ok {
						trackFrequency[string(track.Track.ID)] = trackRecord{frequency: counter.frequency + 1, name: track.Track.Name}
					} else {
						trackFrequency[string(track.Track.ID)] = trackRecord{frequency: 1, name: track.Track.Name}
					}
				}
				if playlistTracks.Next == "" {
					break
				}
				nextPageTracksResponse, err := getRequestResponse(playlistTracks.Next)
				nextpageTracks = playlistTracks.Next
				err = json.Unmarshal(nextPageTracksResponse, &playlistTracks)
				if err != nil {
					fmt.Printf("Error happened while unmarshalling playlist tracks response: %s\n", err.Error())
					return
				}
			}
		}
		nextPageREsponse, err := getRequestResponse(playlistList.Next)
		nextpage = playlistList.Next
		err = json.Unmarshal(nextPageREsponse, &playlistList)
		if err != nil {
			fmt.Printf("Error happened while unmarshalling next page playlist response: %s\n", err.Error())
			return
		}

	}
	fmt.Println(len(trackFrequency))
	fmt.Println("Winners!!!!!!!")
	for _, v := range trackFrequency {
		if v.frequency > 20 {
			fmt.Println("counter:", v.frequency, "track name:", v.name)
		}

	}
}
