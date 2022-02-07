package main

import "os/exec"
import "bufio"
import "bytes"
import "fmt"
import "strings"
import "encoding/json"

type Data struct {
	Id     string `json:"id"`
	Source string `json:"source"`
	Target string `json:"target"`

	Name     string `json:"name"`
	Nodetype string `json:"type"`
}

type Entry struct {
	Group string `json:"group"`
	Data  Data   `json:"data"`
}

type Inspect struct {
	Image string
}

func getImages() (map[string]string, error) {
	out, err := exec.Command("docker", "images", "--no-trunc").Output()
	if err != nil {
		return nil, err
	}
	s := bufio.NewScanner(bytes.NewReader(out))

	imageMap := map[string]string{}

	for s.Scan() {
		col := strings.Fields(s.Text())
		if col[2] == "IMAGE" {
			continue
		}
		imageMap[col[2]] = col[0] + ":" + col[1]
	}
	if s.Err() != nil {
		return nil, s.Err()
	}

	return imageMap, nil
}

func getContainers() (map[string]string, error) {
	out, err := exec.Command("docker", "ps", "-a", "--no-trunc", "--format", "{{.Names}} {{.ID}}").Output()
	if err != nil {
		return nil, err
	}

	containers := map[string]string{}

	s := bufio.NewScanner(bytes.NewReader(out))
	for s.Scan() {
		col := strings.Fields(s.Text())
		out2, err := exec.Command("docker", "inspect", col[1]).Output()
		if err != nil {
			return nil, err
		}
		var detail []Inspect
		err = json.Unmarshal(out2, &detail)
		if err != nil {
			return nil, err
		}
		containers[col[0]] = detail[0].Image
	}
	if s.Err() != nil {
		return nil, s.Err()
	}

	return containers, err
}

func findParentName(images map[string]string, imageID string) (string, error) {
	out, err := exec.Command("docker", "history", "-q", "--no-trunc", imageID).Output()
	if err != nil {
		return "", err
	}
	s := bufio.NewScanner(bytes.NewReader(out))
	isFirst := true
	for s.Scan() {
		if isFirst {
			isFirst = false
			continue
		}
		if name, ok := images[s.Text()]; ok {
			return name, nil
		}
	}
	if s.Err() != nil {
		return "", s.Err()
	}

	return "", nil
}

func MakeDep2() ([]Entry, error) {

	// (image ID) -> (image Name)
	images, err := getImages()
	if err != nil {
		return nil, err
	}

	// (container name -> image ID)
	containers, err := getContainers()
	if err != nil {
		return nil, err
	}

	en := 1
	var ents []Entry

	for k, v := range images {
		ents = append(ents, Entry{
			"nodes",
			Data{
				Id:       v,
				Name:     v,
				Nodetype: "image",
			},
		})
		parentName, err := findParentName(images, k)
		if err != nil {
			return nil, err
		}
		if parentName != "" {
			ents = append(ents, Entry{
				"edges",
				Data{
					Id:     fmt.Sprintf("e%v", en),
					Source: v,
					Target: parentName,
				},
			})
			en = en + 1
		}
	}

	for k, v := range containers {
		ents = append(ents, Entry{
			"nodes",
			Data{
				Id:       k,
				Name:     k,
				Nodetype: "container",
			},
		})
		ents = append(ents, Entry{
			"edges",
			Data{
				Id:     fmt.Sprintf("e%v", en),
				Source: k,
				Target: images[v],
			},
		})
		en = en + 1
	}
	return ents, nil
}
