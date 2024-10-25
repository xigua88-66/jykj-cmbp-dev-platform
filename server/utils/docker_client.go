// -*- coding: utf-8 -*-
// @Time    : 2024/10/21 16:43
// @Author  : chuzhongtian
// @File    : docker_client.go
// @Software: GoLand
// @Project : jykj-cmbp-dev-platform
// @Contact : chuzhongtian1688@gmail.com

package utils

import (
	"errors"
	"fmt"
	"github.com/docker/docker/api/types"
	contai "github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"
	"golang.org/x/net/context"
	"jykj-cmbp-dev-platform/server/global"
)

func GetDockerClient() (cli *client.Client, err error) {
	// 创建一个新的 Docker 客户端
	cli, err = client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		global.CMBP_LOG.Error(fmt.Sprintf("无法创建 Docker 客户端: %v", err.Error()))
		return nil, errors.New(fmt.Sprintf("无法创建 Docker 客户端: %v", err))
	}
	return cli, nil
}

// ImageExists 检查指定的镜像是否存在
func ImageExists(imageName string) (bool, error) {
	cli, err := GetDockerClient()
	defer cli.Close()
	if err != nil {
		return false, err
	}
	// 获取所有镜像
	images, err := cli.ImageList(context.Background(), image.ListOptions{})
	if err != nil {
		return false, err
	}

	// 遍历所有镜像，查找指定的镜像
	for _, img := range images {
		for _, tag := range img.RepoTags {
			if tag == imageName {
				return true, nil
			}
		}
	}

	return false, nil
}

func DeleteImages(imageName string) (bool, error) {
	cli, err := GetDockerClient()
	if err != nil {
		return false, err
	}
	defer cli.Close()

	// 获取所有基于该镜像的容器
	containerList, err := getContainersByImage(cli, imageName)
	if err != nil {
		global.CMBP_LOG.Error(fmt.Sprintf("获取容器列表失败: %v", err.Error()))
		return false, err
	}

	// 停止并删除所有基于该镜像的容器
	for _, container := range containerList {
		if err = stopAndRemoveContainer(cli, container.ID); err != nil {
			global.CMBP_LOG.Error(fmt.Sprintf("停止和删除容器 %s 失败: %v", container.ID, err.Error()))
			return false, err
		} else {
			global.CMBP_LOG.Error(fmt.Sprintf("容器 %s 已成功停止并删除", container.ID))
		}
	}

	// 删除镜像
	if err = removeImage(cli, imageName); err != nil {
		global.CMBP_LOG.Error(fmt.Sprintf("删除镜像 %s 失败: %v", imageName, err.Error()))
		return false, err
	} else {
		global.CMBP_LOG.Info(fmt.Sprintf("镜像 %s 已成功删除", imageName))
	}
	return true, nil
}

// getContainersByImage 获取基于指定镜像的所有容器
func getContainersByImage(cli *client.Client, imageName string) ([]types.Container, error) {
	filters := filters.NewArgs()
	filters.Add("ancestor", imageName)

	options := contai.ListOptions{
		All:     true,
		Filters: filters,
	}
	containers, err := cli.ContainerList(context.Background(), options)
	if err != nil {
		return nil, err
	}
	var containerList []types.Container

	for _, container := range containers {
		inspect, err := cli.ContainerInspect(context.Background(), container.ID)
		if err != nil {
			return nil, err
		}
		if inspect.Config.Image == imageName {
			containerList = append(containerList, container)
		}
	}
	return containerList, nil

}

// stopAndRemoveContainer 停止并删除指定的容器
func stopAndRemoveContainer(cli *client.Client, containerID string) error {
	// 停止容器
	if err := cli.ContainerStop(context.Background(), containerID, contai.StopOptions{}); err != nil && !client.IsErrNotFound(err) {
		return fmt.Errorf("停止容器 %s 失败: %w", containerID, err)
	}

	// 删除容器
	if err := cli.ContainerRemove(context.Background(), containerID, contai.RemoveOptions{}); err != nil && !client.IsErrNotFound(err) {
		return fmt.Errorf("删除容器 %s 失败: %w", containerID, err)
	}

	return nil
}

// removeImage 删除指定的镜像
func removeImage(cli *client.Client, imageName string) error {
	_, err := cli.ImageRemove(context.Background(), imageName, image.RemoveOptions{})
	if err != nil && !client.IsErrNotFound(err) {
		return fmt.Errorf("删除镜像 %s 失败: %w", imageName, err)
	}

	return nil
}
