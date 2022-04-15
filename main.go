package main

import (
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/spf13/viper"
	"github.com/xjh22222228/ip"
	"log"
)

type Access struct {
	Key    string
	Secret string
}

type Config struct {
	Access          Access
	RegionId        string
	SecurityGroupId string
	Description     string
}

func main() {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	var configuration Config

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}
	err := viper.Unmarshal(&configuration)
	if err != nil {
		log.Fatalf("unable to decode into struct, %v", err)
	}
	log.Printf("key is %s", configuration.Access.Key)
	log.Printf("secret is %s", configuration.Access.Secret)

	outboundIP, err := getOutboundIP()
	if err != nil {
		return
	}
	log.Printf("local ip is  %s", outboundIP)

	permissions, err := getPermission(configuration)
	if err != nil {
		return
	}
	filteredPermissions, err := filterPermissionsByDescription(permissions, configuration)
	if err != nil {
		return
	}
	for _, permission := range filteredPermissions {
		err := handlePermission(permission, configuration, outboundIP)
		if err != nil {
			return
		}
	}

}

func getOutboundIP() (string, error) {
	return ip.V4()
}

func handlePermission(permission ecs.Permission, config Config, localIp string) (err error) {
	if permission.SourceCidrIp != localIp {
		err = rmSecurityGroupPermission(permission, config)
		if err != nil {
			return
		}
		err = addSecurityGroupPermission(permission, config, localIp)
		if err != nil {
			return
		}
	}
	return
}

func filterPermissionsByDescription(permissions ecs.Permissions, config Config) (output []ecs.Permission, err error) {
	for _, s := range permissions.Permission {
		if s.Description == config.Description {
			output = append(output, s)
		}
	}
	return
}

//
func getPermission(config Config) (ecs.Permissions, error) {
	client, err := ecs.NewClientWithAccessKey(config.RegionId, config.Access.Key, config.Access.Secret)
	if err != nil {
		return ecs.Permissions{}, err
	}
	request := ecs.CreateDescribeSecurityGroupAttributeRequest()
	request.SecurityGroupId = config.SecurityGroupId

	response, err := client.DescribeSecurityGroupAttribute(request)
	if err != nil {
		// 异常处理
		return ecs.Permissions{}, err
	}
	fmt.Printf("success(%d)! instanceId = %s\n", response.GetHttpStatus(), response.Permissions)
	return response.Permissions, nil
}

func rmSecurityGroupPermission(permission ecs.Permission, config Config) (err error) {
	client, err := ecs.NewClientWithAccessKey(config.RegionId, config.Access.Key, config.Access.Secret)
	if err != nil {
		return
	}
	request := ecs.CreateRevokeSecurityGroupRequest()
	request.Policy = permission.Policy
	request.Description = permission.Description
	request.Priority = permission.Priority
	request.NicType = permission.NicType
	request.PortRange = permission.PortRange
	request.SourceCidrIp = permission.SourceCidrIp
	request.IpProtocol = permission.IpProtocol
	request.SecurityGroupId = config.SecurityGroupId
	response, err := client.RevokeSecurityGroup(request)
	if 200 == response.GetHttpStatus() {
		fmt.Printf("success(%d)!  revoke ip %s\n", response.GetHttpStatus(), permission.SourceCidrIp)
	}
	return
}

func addSecurityGroupPermission(permission ecs.Permission, config Config, localIp string) (err error) {
	client, err := ecs.NewClientWithAccessKey(config.RegionId, config.Access.Key, config.Access.Secret)
	if err != nil {
		return
	}
	request := ecs.CreateAuthorizeSecurityGroupRequest()
	request.Policy = permission.Policy
	request.Description = permission.Description
	request.Priority = permission.Priority
	request.NicType = permission.NicType
	request.PortRange = permission.PortRange
	request.SourceCidrIp = localIp
	request.IpProtocol = permission.IpProtocol
	request.SecurityGroupId = config.SecurityGroupId
	response, err := client.AuthorizeSecurityGroup(request)
	if 200 == response.GetHttpStatus() {
		fmt.Printf("success(%d)!  add ip %s\n", response.GetHttpStatus(), localIp)
	}
	return
}
