package backlog4fzf

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
