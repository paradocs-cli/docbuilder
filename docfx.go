package docbuilder
import (
	"fmt"
	"github.com/go-git/go-git/v5"
	gendocs "github.com/johhess40/generatedocs"
	"github.com/paradocs-cli/gengit"
	"log"
	"os"
	"path/filepath"
)

type Repositories []git.Repository

func CloneGitlab(g GitlabData, urls []string)Repositories{
	var r Repositories
	for _, v := range urls {
		cln, err := gengit.CloneRepo(gengit.GitOptions{
			RemoteOptions: struct {
				UserName string
				Pat      string
				Provider string
				RepoUrl  string
			}{UserName: g.UserName, Pat: g.Token, RepoUrl: v },
			LocalOptions: struct{ ClonePath string }{ClonePath: "../DocFxData"},
		})
		if err != nil {
			return r
		}
		r = append(r, *cln)
	}
	return r
}
func CreateDocFxDir() error{
	err := os.Mkdir("../DocFxData", 0644)
	if err != nil {
		return fmt.Errorf("%v", err.Error())
	}
	return nil
}


func GenerateDocFxStructure(r Repositories) error{
	dirs, err := gendocs.GetDirs("../DocFxData")
	if err != nil {
		return fmt.Errorf("error reading directories for gendocs.GetDirs")
	}
	for _,v := range dirs {
		path := filepath.Join("../DocFxData", fmt.Sprintf("%v", v))
		data, errs := gendocs.GetData(v)
		if errs != nil {
			log.Fatalf(errs.Error())
		}
		errss := os.Chdir(fmt.Sprintf("%v", path))
		if errss != nil {
			return fmt.Errorf(errss.Error())

		}
		gendocs.WriteMarkdownTerra(data)
	}
	return nil
}

func BuildGitLabDocs(g GitlabData){
	err := CreateDocFxDir()
	if err != nil {
		return
	}
	data, err := GetGitLabProjectData(g)
	if err != nil {
		return
	}

	repos, err := GetGitlabRepos(data)
	if err != nil {
		return
	}
	reps := CloneGitlab(g, repos)

	err = GenerateDocFxStructure(reps)
	if err != nil {
		return
	}
}