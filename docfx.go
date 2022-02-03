package docbuilder

import (
	"fmt"
	"github.com/go-git/go-git/v5"
	gendocs "github.com/johhess40/generatedocs"
	"github.com/paradocs-cli/gengit"
	"log"
	"os"
	"strings"
)

type Repositories []git.Repository

func CloneGitlab(g GitlabData, urls []string) Repositories {
	var r Repositories
	for _, v := range urls {
		spl := strings.Split(v, "/")
		cln, err := gengit.CloneRepo(gengit.GitOptions{
			RemoteOptions: struct {
				UserName string
				Pat      string
				Provider string
				RepoUrl  string
			}{UserName: g.UserName, Pat: g.Token, RepoUrl: v},
			LocalOptions: struct{ ClonePath string }{ClonePath: strings.ReplaceAll(spl[len(spl)-1], ".git", "")},
		})
		if err != nil {
			return r
		}
		r = append(r, *cln)
	}
	return r
}
func CreateDocFxDir() error {
	err := os.Mkdir("../DocFxData", 0644)
	if err != nil {
		return fmt.Errorf("%v", err.Error())
	}
	return nil
}

func GenerateDocFxStructure(g GitlabData, p GitlabProjectData) error {
		err := os.Chdir(fmt.Sprintf("%v", "../DocFxData"))
		if err != nil {
			return fmt.Errorf(err.Error())

		}

		mak := os.MkdirAll(fmt.Sprintf("%s", p.Name), 0644)
		if mak != nil {
			return fmt.Errorf("%v", mak.Error())
		}

		err = os.Chdir(fmt.Sprintf("%s", p.Name))
		if err != nil {
			return fmt.Errorf(err.Error())

		}

		CloneGitlab(g, []string{p.HttpUrlToRepo})

		dirs, err := gendocs.GetDirs(fmt.Sprintf("%s", p.Name))
		if err != nil {
			return fmt.Errorf("error reading directories for gendocs.GetDirs")
		}
		for _, v := range dirs {
			data, errs := gendocs.GetData(v)
			if errs != nil {
				log.Fatalf(errs.Error())
			}
			gendocs.WriteMarkdownTerra(data)
		}
	return nil
}

func BuildGitLabDocs(g GitlabData) {

	err := CreateDocFxDir()
	if err != nil {
		return
	}

	data, err := GetGitLabProjectData(g)
	if err != nil {
		return
	}

	//repos, err := GetGitlabRepos(data)
	//if err != nil {
	//	return
	//}

	for _, v := range data {
		err = GenerateDocFxStructure(g, v)
		if err != nil {
			return
		}
	}


}
