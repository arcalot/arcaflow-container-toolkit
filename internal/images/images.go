package images

import (
	"strings"

	"go.arcalot.io/arcaflow-container-toolkit/internal/ce_service"
	"go.arcalot.io/arcaflow-container-toolkit/internal/docker"
	"go.arcalot.io/log"
)

func BuildImage(build_img bool, all_checks bool, cec ce_service.ContainerEngineService, abspath string, image_name string,
	image_tag string, architype string, options *docker.BuildOptions, logger log.Logger) error {

	if all_checks && build_img {
		logger.Infof("Passed all requirements: %s %s\n", image_name, image_tag)
		logger.Infof("Building %s %s from %v\n", image_name, image_tag, abspath)
		if err := cec.Build(abspath, image_name, []string{image_tag}, architype, options); err != nil {
			return err
		}
	}
	return nil
}

func PushImage(all_checks, build_image, push_image bool, cec ce_service.ContainerEngineService, name, version,
	username, password, registry_address, registry_namespace string, logger log.Logger) error {

	if all_checks && build_image && push_image {
		image_name_tag := name + ":" + version
		destination := strings.Join(
			[]string{registry_address, registry_namespace, image_name_tag},
			"/")
		logger.Infof("Pushing %s to %s", name, destination)
		err := cec.Tag(image_name_tag, destination)
		if err != nil {
			return err
		}
		err = cec.Push(destination, username, password, registry_address)
		if err != nil {
			return err
		}
	}
	return nil
}
