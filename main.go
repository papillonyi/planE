package main

import (
	"flag"
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/spf13/viper"
	"github.com/xjh22222228/ip"
	"log"
	"time"
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
	SleepTime       int
}

func main() {
	configFile := flag.String("config", "config.yml", "config file")
	flag.Parse()

	viper.SetConfigFile(*configFile)
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

	for {
		run(configuration)
		log.Printf("done, will be called  %dmin late", configuration.SleepTime)
		time.Sleep(time.Duration(configuration.SleepTime) * time.Minute)
	}

}

func run(config Config) {
	outboundIP, err := getOutboundIP()
	if err != nil {
		fmt.Printf("failed to get outbound ip %s \n", err)
		return
	}
	log.Printf("local ip is  %s", outboundIP)

	permissions, err := getPermission(config)
	if err != nil {
		fmt.Printf("failed to get permissions %s \n", err)
		return
	}
	filteredPermissions, err := filterPermissionsByDescription(permissions, config)
	if err != nil {
		fmt.Printf("failed to get filtered permissions %s \n", err)
		return
	}

	fmt.Printf("%d permissions may need to upate \n", len(filteredPermissions))
	for _, permission := range filteredPermissions {
		err := handlePermission(permission, config, outboundIP)
		if err != nil {
			fmt.Printf("failed to handle permission permissions %s \n", err)

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
			fmt.Printf("failed to rm security group permission %s \n", err)
			return
		}
		err = addSecurityGroupPermission(permission, config, localIp)
		if err != nil {
			fmt.Printf("failed to add security group permission %s \n", err)
			return
		}
		fmt.Printf("permission %s has been changed \n", permission.Description)
	} else {
		fmt.Printf("permission %s don't need to change \n", permission.Description)
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
	fmt.Printf("success(%d)! \n", response.GetHttpStatus())
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
