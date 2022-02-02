package docbuilder
import (
	"fmt"
	"github.com/go-git/go-git/v5"
	gendocs "github.com/johhess40/generatedocs"
	"github.com/paradocs-cli/gengit"
	"log"
	"os"
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
			LocalOptions: struct{ ClonePath string }{ClonePath: ".."},
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
		err := os.Mkdir(fmt.Sprintf("../DocFxData/%s", v), 0644)
		if err != nil {
			return fmt.Errorf("error making directories for os.Mkdir")

		}
		data, err := gendocs.GetData(v)
		if err != nil {
			log.Fatalf(err.Error())
		}
		err = os.Chdir(fmt.Sprintf("../DocFxData/%s", v))
		if err != nil {
			return fmt.Errorf("error switching directories for markdown creation")

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