package backlog4fzf

import (
	"fmt"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
	"os"
)

var (
	app              = kingpin.New("backlog4fzf", "top command")
	appProfiles      = app.Flag("profiles", "profiles").Default("").String()
	appProfileAll    = app.Flag("profile-all", "profile-all").Bool()
	appProfileConfig = app.Flag("profile-config", "profile-config").Default(os.Getenv("HOME") + "/.backlog/profiles").String()
	appCacheDir      = app.Flag("cache-dir", "cache-dir").Default(os.Getenv("HOME") + "/.backlog/cache").String()
)

func Run() int {
	cmdIssueList := app.Command("issue-list", "issue-list")
	cmdIssueListWithDescription := cmdIssueList.Flag("desc", "description").Bool()

	cmdIssueUrls := app.Command("issue-urls", "issue-urls")
	cmdIssueUrlArgs := cmdIssueUrls.Arg("profile-issues", "profile-issues").Strings()

	cmdDeleteIssueCache := app.Command("delete-issue-cache", "delete-issue-cache")
	cmdDeleteIssueCacheProfiles := cmdDeleteIssueCache.Arg("profile-issues", "profile-issues").Strings()

	cmdUpdateIssueStatus := app.Command("update-issue-status", "update-issue-status")
	cmdUpdateIssueStatusTarget := cmdUpdateIssueStatus.Arg("target-status", "target-status").String()
	cmdUpdateIssueStatusIssues := cmdUpdateIssueStatus.Arg("profile-issues", "profile-issues").Strings()

	cmdIssueDescription := app.Command("issue-description", "issue-description")
	cmdIssueDescriptionIssue := cmdIssueDescription.Arg("profile-issue", "profile-issue").String()

	cmdIssueComments := app.Command("issue-comments", "issue-comments")
	cmdIssueCommentsIssue := cmdIssueComments.Arg("profile-issue", "profile-issue").String()

	cmdIssueDescriptionAndComments := app.Command("issue-description-and-comments", "issue-description-and-comments")
	cmdIssueDescriptionAndCommentsIssue := cmdIssueDescriptionAndComments.Arg("profile-issue", "profile-issue").String()

	cmdPullrequestList := app.Command("pullrequest-list", "pullrequest-list")
	cmdPullrequestListWithDescription := cmdPullrequestList.Flag("desc", "description").Bool()

	cmdPullrequestUrls := app.Command("pullrequest-urls", "pullrequest-urls")
	cmdPullrequestUrlArgs := cmdPullrequestUrls.Arg("profile-repository-pullrequests", "profile-repository-pullrequests").Strings()

	cmdDeletePullrequestCache := app.Command("delete-pullrequest-cache", "delete-pullrequest-cache")
	cmdDeletePullrequestCacheProfiles := cmdDeletePullrequestCache.Arg("profile-pullrequests", "profile-pullrequests").Strings()

	cmdPullrequestDescription := app.Command("pullrequest-description", "pullrequest-description")
	cmdPullrequestDescriptionPullrequest := cmdPullrequestDescription.Arg("profile-repository-pullrequest", "profile-repository-pullrequest").String()

	switch kingpin.MustParse(app.Parse(os.Args[1:])) {
	case cmdIssueList.FullCommand():
		exit, err := printIssueList(*cmdIssueListWithDescription)
		if err != nil {
			fmt.Println(err)
		}
		return exit
	case cmdIssueUrls.FullCommand():
		exit, err := printIssueUrls(*cmdIssueUrlArgs)
		if err != nil {
			fmt.Println(err)
		}
		return exit
	case cmdDeleteIssueCache.FullCommand():
		exit, err := deleteIssueCache(*cmdDeleteIssueCacheProfiles)
		if err != nil {
			fmt.Println(err)
		}
		return exit
	case cmdUpdateIssueStatus.FullCommand():
		exit, err := updateIssueStatus(*cmdUpdateIssueStatusTarget, *cmdUpdateIssueStatusIssues)
		if err != nil {
			fmt.Println(err)
		}
		return exit
	case cmdIssueDescription.FullCommand():
		exit, err := printIssueDescription(*cmdIssueDescriptionIssue)
		if err != nil {
			fmt.Println(err)
		}
		return exit
	case cmdIssueComments.FullCommand():
		exit, err := printIssueComments(*cmdIssueCommentsIssue)
		if err != nil {
			fmt.Println(err)
		}
		return exit
	case cmdIssueDescriptionAndComments.FullCommand():
		exit, err := printIssueDescriptionAndComments(*cmdIssueDescriptionAndCommentsIssue)
		if err != nil {
			fmt.Println(err)
		}
		return exit
	case cmdPullrequestList.FullCommand():
		exit, err := printPullrequestList(*cmdPullrequestListWithDescription)
		if err != nil {
			fmt.Println(err)
		}
		return exit
	case cmdPullrequestUrls.FullCommand():
		exit, err := printPullrequestUrls(*cmdPullrequestUrlArgs)
		if err != nil {
			fmt.Println(err)
		}
		return exit
	case cmdDeletePullrequestCache.FullCommand():
		exit, err := deletePullrequestCache(*cmdDeletePullrequestCacheProfiles)
		if err != nil {
			fmt.Println(err)
		}
		return exit
	case cmdPullrequestDescription.FullCommand():
		exit, err := printPullrequestDescription(*cmdPullrequestDescriptionPullrequest)
		if err != nil {
			fmt.Println(err)
		}
		return exit
	}

	/*
		profilesPtr := flag.String("profiles", "", "comma separated profiles")
		refreshAllPtrPtr := flag.Bool("refresh-all", false, "reload forcely")
		flag.Parse()

		if len(flag.Args()) == 0 {
			log.Fatal("引数がありません")
		}

		cacheDir := ""
		if val, ok := os.LookupEnv("FZF_BACKLOG_CACHE_DIR"); ok {
			cacheDir = val
		} else {
			cacheDir = os.Getenv("HOME") + "/.backlog/cache"
		}

		if flag.Args()[0] == "issues" {
			profiles := *profilesPtr
			if val, ok := os.LookupEnv("FZF_BACKLOG_PROFILES"); profiles == "" && ok {
				profiles = val
			}

			if len(profiles) == 0 {
				return 1
			}
			for _, profile := range strings.Split(profiles, ",") {
				backlogBaseUrl, projectId, apiKey, err := getBacklogProfile(profile, os.Getenv("HOME")+"/.backlog/profiles")
				if err != nil {
					log.Fatal(err)
				}
				printIssues(profile, backlogBaseUrl, projectId, apiKey, cacheDir, *refreshAllPtrPtr)
			}
		} else if flag.Args()[0] == "issues-across-projects" {
			profileList, err := getProfileList(os.Getenv("HOME") + "/.backlog/profiles")
			if err != nil {
				log.Fatal(err)
			}
			for _, profile := range profileList {
				backlogBaseUrl, projectId, apiKey, err := getBacklogProfile(profile, os.Getenv("HOME")+"/.backlog/profiles")
				printIssues(profile, backlogBaseUrl, projectId, apiKey, cacheDir, *refreshAllPtrPtr)
				if err != nil {
					log.Fatal(err)
				}
			}
		} else if flag.Args()[0] == "issue-urls" {
			urls := []string{}
			for _, issue := range flag.Args()[1:] {
				sp := strings.Split(issue, ":")
				profile := sp[0]
				issueId := sp[1]
				backlogBaseUrl, _, _, err := getBacklogProfile(profile, os.Getenv("HOME")+"/.backlog/profiles")
				if err != nil {
					log.Fatal(err)
				}
				urls = append(urls, "https://"+backlogBaseUrl+"/view/"+issueId)
			}
			fmt.Println(strings.Join(urls, " "))
		} else if flag.Args()[0] == "pullrequests" {
			profiles := *profilesPtr
			if val, ok := os.LookupEnv("FZF_BACKLOG_PROFILES"); profiles == "" && ok {
				profiles = val
			}

			if len(profiles) == 0 {
				return 1
			}
			for _, profile := range strings.Split(profiles, ",") {
				backlogBaseUrl, projectId, apiKey, err := getBacklogProfile(profile, os.Getenv("HOME")+"/.backlog/profiles")
				if err != nil {
					log.Fatal(err)
				}
				printPullrequests(profile, backlogBaseUrl, projectId, apiKey, cacheDir, *refreshAllPtrPtr)
			}
		} else if flag.Args()[0] == "pullrequests-across-projects" {
			profileList, err := getProfileList(os.Getenv("HOME") + "/.backlog/profiles")
			if err != nil {
				log.Fatal(err)
			}
			for _, profile := range profileList {
				backlogBaseUrl, projectId, apiKey, err := getBacklogProfile(profile, os.Getenv("HOME")+"/.backlog/profiles")
				printPullrequests(profile, backlogBaseUrl, projectId, apiKey, cacheDir, *refreshAllPtrPtr)
				if err != nil {
					log.Fatal(err)
				}
			}
		} else if flag.Args()[0] == "pullrequest-urls" {
			urls := []string{}
			for _, issue := range flag.Args()[1:] {
				sp := strings.Split(issue, ":")
				profile := sp[0]
				repository := sp[1]
				number := sp[2][1:]
				backlogBaseUrl, projectId, _, err := getBacklogProfile(profile, os.Getenv("HOME")+"/.backlog/profiles")
				if err != nil {
					log.Fatal(err)
				}
				urls = append(urls, "https://"+backlogBaseUrl+"/git/"+projectId+"/"+repository+"/pullRequests/"+number)
			}
			fmt.Println(strings.Join(urls, " "))
		} else if flag.Args()[0] == "wikis" {
			profiles := *profilesPtr
			if val, ok := os.LookupEnv("FZF_BACKLOG_PROFILES"); profiles == "" && ok {
				profiles = val
			}

			if len(profiles) == 0 {
				return 1
			}
			for _, profile := range strings.Split(profiles, ",") {
				backlogBaseUrl, projectId, apiKey, err := getBacklogProfile(profile, os.Getenv("HOME")+"/.backlog/profiles")
				if err != nil {
					log.Fatal(err)
				}
				printWikis(profile, backlogBaseUrl, projectId, apiKey, cacheDir, *refreshAllPtrPtr)
			}
		} else if flag.Args()[0] == "wikis-across-projects" {
			profileList, err := getProfileList(os.Getenv("HOME") + "/.backlog/profiles")
			if err != nil {
				log.Fatal(err)
			}
			for _, profile := range profileList {
				backlogBaseUrl, projectId, apiKey, err := getBacklogProfile(profile, os.Getenv("HOME")+"/.backlog/profiles")
				printWikis(profile, backlogBaseUrl, projectId, apiKey, cacheDir, *refreshAllPtrPtr)
				if err != nil {
					log.Fatal(err)
				}
			}
		} else if flag.Args()[0] == "wiki-content" {
			sp := strings.Split(flag.Args()[1], ":")
			profile := sp[0]
			wikiId := sp[1]
			backlogBaseUrl, projectId, apiKey, err := getBacklogProfile(profile, os.Getenv("HOME")+"/.backlog/profiles")
			if err != nil {
				log.Fatal(err)
			}
			printWiki(profile, wikiId, backlogBaseUrl, projectId, apiKey, cacheDir, *refreshAllPtrPtr)
		} else if flag.Args()[0] == "wiki-urls" {
			urls := []string{}
			for _, wiki := range flag.Args()[1:] {
				sp := strings.Split(wiki, ":")
				profile := sp[0]
				id := sp[1]
				backlogBaseUrl, _, _, err := getBacklogProfile(profile, os.Getenv("HOME")+"/.backlog/profiles")
				if err != nil {
					log.Fatal(err)
				}
				urls = append(urls, "https://"+backlogBaseUrl+"/alias/wiki/"+id)
			}
			fmt.Println(strings.Join(urls, " "))
		} else if flag.Args()[0] == "delete-issues-cache" {
			profileSet := map[string]bool{}
			for _, profile := range strings.Split(flag.Args()[1], " ") {
				profileSet[profile] = true
			}
			for profile, _ := range profileSet {
				backlogBaseUrl, projectId, apiKey, err := getBacklogProfile(profile, os.Getenv("HOME")+"/.backlog/profiles")
				if err != nil {
					log.Fatal(err)
				}
				deleteIssuesCache(profile, backlogBaseUrl, projectId, apiKey, cacheDir, *refreshAllPtrPtr)
			}
		} else if flag.Args()[0] == "delete-wiki-content-cache" {
			sp := strings.Split(flag.Args()[1], ":")
			profile := sp[0]
			wikiId := sp[1]
			backlogBaseUrl, projectId, apiKey, err := getBacklogProfile(profile, os.Getenv("HOME")+"/.backlog/profiles")
			if err != nil {
				log.Fatal(err)
			}
			deleteWikiContentCache(profile, wikiId, backlogBaseUrl, projectId, apiKey, cacheDir, *refreshAllPtrPtr)
		} else if flag.Args()[0] == "update-issue-status" {
			targetStatus := flag.Args()[1]
			for _, issue := range flag.Args()[2:] {
				sp := strings.Split(issue, ":")
				profile := sp[0]
				issueId := sp[1]
				backlogBaseUrl, projectId, apiKey, err := getBacklogProfile(profile, os.Getenv("HOME")+"/.backlog/profiles")
				if err != nil {
					log.Fatal(err)
				}
				updateIssueStatus(profile, issueId, backlogBaseUrl, projectId, apiKey, cacheDir, *refreshAllPtrPtr, targetStatus)
			}
		}
	*/
	return 0
}

/*
func getProfileList(profilePath string) ([]string, error) {
	tree, err := toml.LoadFile(profilePath)
	if err != nil {
		return []string{}, fmt.Errorf("profileファイルの読み込みに失敗")
	}
	profileList := []string{}
	for _, a := range tree.Keys() {
		profileList = append(profileList, a)
	}
	return profileList, nil
}

func printPullrequests(profile string, backlogBaseUrl string, projectId string, apiKey string, cacheDir string, refreshAll bool) {
	setEnv(backlogBaseUrl, projectId, apiKey, cacheDir)
	allPullrequests, err := getAllPullrequests(refreshAll)
	if err != nil {
		log.Fatal(err)
	}
	jsonObj := map[string]interface{}{}
	jsonObj["profile"] = profile
	jsonObj["repositories"] = allPullrequests
	_ = json.NewEncoder(os.Stdout).Encode(jsonObj)
}

func printWikis(profile string, backlogBaseUrl string, projectId string, apiKey string, cacheDir string, refreshAll bool) {
	setEnv(backlogBaseUrl, projectId, apiKey, cacheDir)
	allWikis, err := getAllWikis(refreshAll)
	if err != nil {
		log.Fatal(err)
	}
	jsonObj := map[string]interface{}{}
	jsonObj["profile"] = profile
	jsonObj["wikis"] = allWikis
	_ = json.NewEncoder(os.Stdout).Encode(jsonObj)
}

func printWiki(profile string, wikiId string, backlogBaseUrl string, projectId string, apiKey string, cacheDir string, refreshAll bool) {
	setEnv(backlogBaseUrl, projectId, apiKey, cacheDir)
	wikiContent, err := getWikiContent(wikiId, refreshAll)
	if err != nil {
		log.Fatal(err)
	}
	jsonObj := map[string]interface{}{}
	jsonObj["profile"] = profile
	jsonObj["content"] = wikiContent
	_ = json.NewEncoder(os.Stdout).Encode(jsonObj)
}

func deleteIssuesCache(profile string, backlogBaseUrl string, projectId string, apiKey string, cacheDir string, refreshAll bool) {
	setEnv(backlogBaseUrl, projectId, apiKey, cacheDir)
	cachePath := getIssuesCache()
	fmt.Println("Remove " + cachePath)
	if err := os.Remove(cachePath); err != nil {
		fmt.Println(err)
	}
}

func deleteWikiContentCache(profile string, wikiId string, backlogBaseUrl string, projectId string, apiKey string, cacheDir string, refreshAll bool) {
	setEnv(backlogBaseUrl, projectId, apiKey, cacheDir)
	cachePath := getWikisCache()
	fmt.Println("Remove " + cachePath)
	if err := os.Remove(cachePath); err != nil {
		fmt.Println(err)
	}
}
*/
