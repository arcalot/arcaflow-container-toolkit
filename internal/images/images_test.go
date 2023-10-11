package images_test

import (
	"fmt"
	"testing"

	"go.arcalot.io/arcaflow-container-toolkit/internal/docker"
	"go.arcalot.io/arcaflow-container-toolkit/internal/dto"
	"go.arcalot.io/arcaflow-container-toolkit/internal/images"
	mocks "go.arcalot.io/arcaflow-container-toolkit/mocks/ce_service"
	"go.arcalot.io/assert"
	"go.arcalot.io/log/v2"
	"go.uber.org/mock/gomock"
)

func TestBuildImage(t *testing.T) {
	logger := log.NewLogger(log.LevelInfo, log.NewNOOPLogger())
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	cec := mocks.NewMockContainerEngineService(ctrl)
	cec.EXPECT().
		Build("use", "the", []string{"forks"}, "amd64", docker.DefaultBuildOptions()).
		Return(nil).
		Times(1)
	assert.Nil(t, images.BuildImage(true, true, cec, "use", "the", "forks", "amd64", docker.DefaultBuildOptions(), logger))
}

func TestPushImage(t *testing.T) {
	logger := log.NewLogger(log.LevelInfo, log.NewNOOPLogger())
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	cec := mocks.NewMockContainerEngineService(ctrl)
	rg1 := dto.Registry{
		Url:       "reg1.io",
		Username:  "user1",
		Password:  "secret1",
		Namespace: "allyourbases",
	}
	image_name := "usethe"
	image_tag := "forks"

	destination := fmt.Sprintf("%s/%s/%s:%s", rg1.Url, rg1.Namespace, image_name, image_tag)
	cec.EXPECT().
		Tag(fmt.Sprintf("%s:%s", image_name, image_tag), destination).
		Return(nil).
		Times(1)
	cec.EXPECT().
		Push(destination, rg1.Username, rg1.Password, rg1.Url).
		Return(nil).
		Times(1)
	assert.Nil(t, images.PushImage(true, true, true, cec, image_name, image_tag,
		rg1.Username, rg1.Password, rg1.Url, rg1.Namespace, logger))
}
