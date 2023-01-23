package version

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	//"log"
	"os/exec"
)

type requestBody struct {
	Url string `json:"url"`
}

type Reference struct {
	CommitId  string `json:"commitId"`
	Reference string `json:"reference"`
}

type Response struct {
	Tags  []Reference `json:"tags"`
	Heads []Reference `json:"heads"`
}

func parseReferences(refs []byte) []Reference {
	lines := bytes.Split(refs, []byte("\n"))

	var parsedRefs []Reference
	for j := range lines {
		rows := bytes.Split(lines[j], []byte("\t"))
		if len(rows) == 2 {
			commitId := string(rows[0])

			reference := string(rows[1])
			reference = strings.TrimPrefix(reference, "refs/tags/")
			reference = strings.TrimPrefix(reference, "refs/heads/")

			parsedRefs = append(parsedRefs, Reference{commitId, reference})
		}
	}

	return parsedRefs
}

func buildJSONResponse(tags []byte, heads []byte) []byte {
	response := Response{
		Tags:  parseReferences(tags),
		Heads: parseReferences(heads),
	}

	responseJSON, err := json.Marshal(response)

	if err != nil {
		fmt.Println("Error building the response")
	}

	return responseJSON
}

func GetVersions(w http.ResponseWriter, r *http.Request) {
	var req requestBody

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&req)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	u, err := GetAuthenticatedUser(r)
	if err != nil || u == nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	cmdTags := exec.Command("git", "ls-remote", "--tags", req.Url)
	cmdHeads := exec.Command("git", "ls-remote", "--heads", req.Url)

	tags, err1 := cmdTags.Output()
	heads, err1 := cmdHeads.Output()

	if err1 != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	jsonResponse := buildJSONResponse(tags, heads)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err2 := w.Write(jsonResponse)

	if err2 != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
