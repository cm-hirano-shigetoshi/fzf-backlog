package backlog4fzf

import (
	"fmt"
	toml "github.com/pelletier/go-toml"
	"os"
	"strings"
)

var BACKLOG_BASE_URL string
var PROJECT_ID string
var API_KEY string
var CACHE_DIR string

func setEnv(backlogBaseUrl string, projectId string, apiKey string, cacheDir string) {
	BACKLOG_BASE_URL = backlogBaseUrl
	PROJECT_ID = projectId
	API_KEY = apiKey
	CACHE_DIR = cacheDir
}

type BacklogProfile struct {
	baseUrl   string
	projectId string
	apiKey    string
}

func getProfiles() string {
	if *appProfileAll {
		profiles, _ := getAllProfiles()
		return strings.Join(profiles, ",")
	}
	if len(*appProfiles) > 0 {
		return *appProfiles
	} else if val, ok := os.LookupEnv("FZF_BACKLOG_PROFILES"); ok {
		return val
	} else {
		return ""
	}
}

func getAllProfiles() ([]string, error) {
	tree, err := toml.LoadFile(*appProfileConfig)
	if err != nil {
		return []string{}, fmt.Errorf("profileファイルの読み込みに失敗")
	}
	keys := tree.Keys()
	return keys, nil
}

func getBacklogProfile(profile string, profileFile string) (BacklogProfile, error) {
	tree, err := toml.LoadFile(profileFile)
	if err != nil {
		return BacklogProfile{}, fmt.Errorf("profileファイルの読み込みに失敗")
	}
	if !tree.Has(profile) {
		return BacklogProfile{}, fmt.Errorf("profile: " + profile + " がありません")
	}
	backlogBaseUrl := tree.GetDefault(profile+".base_url", "").(string)
	projectId := tree.GetDefault(profile+".project_id", "").(string)
	apiKey := tree.GetDefault(profile+".api_key", "").(string)
	if len(backlogBaseUrl)*len(projectId)*len(apiKey) == 0 {
		return BacklogProfile{}, fmt.Errorf("profileに必要な設定がありません")
	}
	return BacklogProfile{backlogBaseUrl, projectId, apiKey}, nil
}
