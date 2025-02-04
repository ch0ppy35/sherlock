package ec2

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
)

type DescribeRegionsInput = ec2.DescribeRegionsInput
type DescribeRegionsOutput = ec2.DescribeRegionsOutput
type DescribeVpcsInput = ec2.DescribeVpcsInput
type DescribeVpcsOutput = ec2.DescribeVpcsOutput
type DeleteVpcInput = ec2.DeleteVpcInput
type DeleteVpcOutput = ec2.DeleteVpcOutput
type DescribeSecurityGroupsInput = ec2.DescribeSecurityGroupsInput
type DescribeSecurityGroupsOutput = ec2.DescribeSecurityGroupsOutput
type DeleteSecurityGroupInput = ec2.DeleteSecurityGroupInput
type DeleteSecurityGroupOutput = ec2.DeleteSecurityGroupOutput
type DescribeSubnetsInput = ec2.DescribeSubnetsInput
type DescribeSubnetsOutput = ec2.DescribeSubnetsOutput
type DeleteSubnetInput = ec2.DeleteSubnetInput
type DeleteSubnetOutput = ec2.DeleteSubnetOutput
type DescribeRouteTablesInput = ec2.DescribeRouteTablesInput
type DescribeRouteTablesOutput = ec2.DescribeRouteTablesOutput
type DeleteRouteTableInput = ec2.DeleteRouteTableInput
type DeleteRouteTableOutput = ec2.DeleteRouteTableOutput
type DescribeInternetGatewaysInput = ec2.DescribeInternetGatewaysInput
type DescribeInternetGatewaysOutput = ec2.DescribeInternetGatewaysOutput
type DetachInternetGatewayInput = ec2.DetachInternetGatewayInput
type DetachInternetGatewayOutput = ec2.DetachInternetGatewayOutput
type DeleteInternetGatewayInput = ec2.DeleteInternetGatewayInput
type DeleteInternetGatewayOutput = ec2.DeleteInternetGatewayOutput
type DescribeNetworkAclsInput = ec2.DescribeNetworkAclsInput
type DescribeNetworkAclsOutput = ec2.DescribeNetworkAclsOutput
type DeleteNetworkAclInput = ec2.DeleteNetworkAclInput
type DeleteNetworkAclOutput = ec2.DeleteNetworkAclOutput

// EC2 Client
type EC2API interface {
	DescribeRegions(ctx context.Context, input *ec2.DescribeRegionsInput, optFns ...func(*ec2.Options)) (*ec2.DescribeRegionsOutput, error)
	DescribeVpcs(ctx context.Context, input *ec2.DescribeVpcsInput, optFns ...func(*ec2.Options)) (*ec2.DescribeVpcsOutput, error)
	DeleteVpc(ctx context.Context, input *ec2.DeleteVpcInput, optFns ...func(*ec2.Options)) (*ec2.DeleteVpcOutput, error)
	DescribeSecurityGroups(ctx context.Context, input *ec2.DescribeSecurityGroupsInput, optFns ...func(*ec2.Options)) (*ec2.DescribeSecurityGroupsOutput, error)
	DeleteSecurityGroup(ctx context.Context, input *ec2.DeleteSecurityGroupInput, optFns ...func(*ec2.Options)) (*ec2.DeleteSecurityGroupOutput, error)
	DescribeSubnets(ctx context.Context, input *ec2.DescribeSubnetsInput, optFns ...func(*ec2.Options)) (*ec2.DescribeSubnetsOutput, error)
	DeleteSubnet(ctx context.Context, input *ec2.DeleteSubnetInput, optFns ...func(*ec2.Options)) (*ec2.DeleteSubnetOutput, error)
	DescribeRouteTables(ctx context.Context, input *ec2.DescribeRouteTablesInput, optFns ...func(*ec2.Options)) (*ec2.DescribeRouteTablesOutput, error)
	DeleteRouteTable(ctx context.Context, input *ec2.DeleteRouteTableInput, optFns ...func(*ec2.Options)) (*ec2.DeleteRouteTableOutput, error)
	DescribeInternetGateways(ctx context.Context, input *ec2.DescribeInternetGatewaysInput, optFns ...func(*ec2.Options)) (*ec2.DescribeInternetGatewaysOutput, error)
	DetachInternetGateway(ctx context.Context, input *ec2.DetachInternetGatewayInput, optFns ...func(*ec2.Options)) (*ec2.DetachInternetGatewayOutput, error)
	DeleteInternetGateway(ctx context.Context, input *ec2.DeleteInternetGatewayInput, optFns ...func(*ec2.Options)) (*ec2.DeleteInternetGatewayOutput, error)
	DescribeNetworkAcls(ctx context.Context, input *ec2.DescribeNetworkAclsInput, optFns ...func(*ec2.Options)) (*ec2.DescribeNetworkAclsOutput, error)
	DeleteNetworkAcl(ctx context.Context, input *ec2.DeleteNetworkAclInput, optFns ...func(*ec2.Options)) (*ec2.DeleteNetworkAclOutput, error)
}

type EC2Client struct {
	Client *ec2.Client
}

func ClientConnect() (context.Context, *EC2Client) {
	ctx := context.Background()
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		fmt.Printf("Unable to load AWS SDK config: %v", err)
		os.Exit(1)
	}

	return ctx, &EC2Client{Client: ec2.NewFromConfig(cfg)}
}

func (c *EC2Client) DescribeRegions(ctx context.Context, input *ec2.DescribeRegionsInput, optFns ...func(*ec2.Options)) (*ec2.DescribeRegionsOutput, error) {
	return c.Client.DescribeRegions(ctx, input, optFns...)
}

func (c *EC2Client) DescribeVpcs(ctx context.Context, input *ec2.DescribeVpcsInput, optFns ...func(*ec2.Options)) (*ec2.DescribeVpcsOutput, error) {
	return c.Client.DescribeVpcs(ctx, input, optFns...)
}

func (c *EC2Client) DeleteVpc(ctx context.Context, input *ec2.DeleteVpcInput, optFns ...func(*ec2.Options)) (*ec2.DeleteVpcOutput, error) {
	return c.Client.DeleteVpc(ctx, input, optFns...)
}

func (c *EC2Client) DescribeSecurityGroups(ctx context.Context, input *ec2.DescribeSecurityGroupsInput, optFns ...func(*ec2.Options)) (*ec2.DescribeSecurityGroupsOutput, error) {
	return c.Client.DescribeSecurityGroups(ctx, input, optFns...)
}

func (c *EC2Client) DeleteSecurityGroup(ctx context.Context, input *ec2.DeleteSecurityGroupInput, optFns ...func(*ec2.Options)) (*ec2.DeleteSecurityGroupOutput, error) {
	return c.Client.DeleteSecurityGroup(ctx, input, optFns...)
}

func (c *EC2Client) DescribeSubnets(ctx context.Context, input *ec2.DescribeSubnetsInput, optFns ...func(*ec2.Options)) (*ec2.DescribeSubnetsOutput, error) {
	return c.Client.DescribeSubnets(ctx, input, optFns...)
}

func (c *EC2Client) DeleteSubnet(ctx context.Context, input *ec2.DeleteSubnetInput, optFns ...func(*ec2.Options)) (*ec2.DeleteSubnetOutput, error) {
	return c.Client.DeleteSubnet(ctx, input, optFns...)
}

func (c *EC2Client) DescribeRouteTables(ctx context.Context, input *ec2.DescribeRouteTablesInput, optFns ...func(*ec2.Options)) (*ec2.DescribeRouteTablesOutput, error) {
	return c.Client.DescribeRouteTables(ctx, input, optFns...)
}

func (c *EC2Client) DeleteRouteTable(ctx context.Context, input *ec2.DeleteRouteTableInput, optFns ...func(*ec2.Options)) (*ec2.DeleteRouteTableOutput, error) {
	return c.Client.DeleteRouteTable(ctx, input, optFns...)
}

func (c *EC2Client) DescribeInternetGateways(ctx context.Context, input *ec2.DescribeInternetGatewaysInput, optFns ...func(*ec2.Options)) (*ec2.DescribeInternetGatewaysOutput, error) {
	return c.Client.DescribeInternetGateways(ctx, input, optFns...)
}

func (c *EC2Client) DetachInternetGateway(ctx context.Context, input *ec2.DetachInternetGatewayInput, optFns ...func(*ec2.Options)) (*ec2.DetachInternetGatewayOutput, error) {
	return c.Client.DetachInternetGateway(ctx, input, optFns...)
}

func (c *EC2Client) DeleteInternetGateway(ctx context.Context, input *ec2.DeleteInternetGatewayInput, optFns ...func(*ec2.Options)) (*ec2.DeleteInternetGatewayOutput, error) {
	return c.Client.DeleteInternetGateway(ctx, input, optFns...)
}

func (c *EC2Client) DescribeNetworkAcls(ctx context.Context, input *ec2.DescribeNetworkAclsInput, optFns ...func(*ec2.Options)) (*ec2.DescribeNetworkAclsOutput, error) {
	return c.Client.DescribeNetworkAcls(ctx, input, optFns...)
}

func (c *EC2Client) DeleteNetworkAcl(ctx context.Context, input *ec2.DeleteNetworkAclInput, optFns ...func(*ec2.Options)) (*ec2.DeleteNetworkAclOutput, error) {
	return c.Client.DeleteNetworkAcl(ctx, input, optFns...)
}

// EC2 Client Mocks
type MockEC2Client struct {
	describeRegionsFunc          func(ctx context.Context, input *ec2.DescribeRegionsInput, optFns ...func(*ec2.Options)) (*ec2.DescribeRegionsOutput, error)
	describeVpcsFunc             func(ctx context.Context, input *ec2.DescribeVpcsInput, optFns ...func(*ec2.Options)) (*ec2.DescribeVpcsOutput, error)
	deleteVpcFunc                func(ctx context.Context, input *ec2.DeleteVpcInput, optFns ...func(*ec2.Options)) (*ec2.DeleteVpcOutput, error)
	describeSubnetsFunc          func(ctx context.Context, input *ec2.DescribeSubnetsInput, optFns ...func(*ec2.Options)) (*ec2.DescribeSubnetsOutput, error)
	deleteSubnetFunc             func(ctx context.Context, input *ec2.DeleteSubnetInput, optFns ...func(*ec2.Options)) (*ec2.DeleteSubnetOutput, error)
	describeRouteTablesFunc      func(ctx context.Context, input *ec2.DescribeRouteTablesInput, optFns ...func(*ec2.Options)) (*ec2.DescribeRouteTablesOutput, error)
	deleteRouteTableFunc         func(ctx context.Context, input *ec2.DeleteRouteTableInput, optFns ...func(*ec2.Options)) (*ec2.DeleteRouteTableOutput, error)
	describeInternetGatewaysFunc func(ctx context.Context, input *ec2.DescribeInternetGatewaysInput, optFns ...func(*ec2.Options)) (*ec2.DescribeInternetGatewaysOutput, error)
	detachInternetGatewayFunc    func(ctx context.Context, input *ec2.DetachInternetGatewayInput, optFns ...func(*ec2.Options)) (*ec2.DetachInternetGatewayOutput, error)
	deleteInternetGatewayFunc    func(ctx context.Context, input *ec2.DeleteInternetGatewayInput, optFns ...func(*ec2.Options)) (*ec2.DeleteInternetGatewayOutput, error)
	describeSecurityGroupsFunc   func(ctx context.Context, input *ec2.DescribeSecurityGroupsInput, optFns ...func(*ec2.Options)) (*ec2.DescribeSecurityGroupsOutput, error)
	deleteSecurityGroupFunc      func(ctx context.Context, input *ec2.DeleteSecurityGroupInput, optFns ...func(*ec2.Options)) (*ec2.DeleteSecurityGroupOutput, error)
	describeNetworkAclsFunc      func(ctx context.Context, input *ec2.DescribeNetworkAclsInput, optFns ...func(*ec2.Options)) (*ec2.DescribeNetworkAclsOutput, error)
	deleteNetworkAclFunc         func(ctx context.Context, input *ec2.DeleteNetworkAclInput, optFns ...func(*ec2.Options)) (*ec2.DeleteNetworkAclOutput, error)
}

func (m *MockEC2Client) DescribeRegions(ctx context.Context, input *ec2.DescribeRegionsInput, optFns ...func(*ec2.Options)) (*ec2.DescribeRegionsOutput, error) {
	return m.describeRegionsFunc(ctx, input, optFns...)
}

func (m *MockEC2Client) DescribeVpcs(ctx context.Context, input *ec2.DescribeVpcsInput, optFns ...func(*ec2.Options)) (*ec2.DescribeVpcsOutput, error) {
	return m.describeVpcsFunc(ctx, input, optFns...)
}

func (m *MockEC2Client) DeleteVpc(ctx context.Context, input *ec2.DeleteVpcInput, optFns ...func(*ec2.Options)) (*ec2.DeleteVpcOutput, error) {
	return m.deleteVpcFunc(ctx, input, optFns...)
}

func (m *MockEC2Client) DescribeSubnets(ctx context.Context, input *ec2.DescribeSubnetsInput, optFns ...func(*ec2.Options)) (*ec2.DescribeSubnetsOutput, error) {
	return m.describeSubnetsFunc(ctx, input, optFns...)
}

func (m *MockEC2Client) DeleteSubnet(ctx context.Context, input *ec2.DeleteSubnetInput, optFns ...func(*ec2.Options)) (*ec2.DeleteSubnetOutput, error) {
	return m.deleteSubnetFunc(ctx, input, optFns...)
}

func (m *MockEC2Client) DescribeRouteTables(ctx context.Context, input *ec2.DescribeRouteTablesInput, optFns ...func(*ec2.Options)) (*ec2.DescribeRouteTablesOutput, error) {
	return m.describeRouteTablesFunc(ctx, input, optFns...)
}

func (m *MockEC2Client) DeleteRouteTable(ctx context.Context, input *ec2.DeleteRouteTableInput, optFns ...func(*ec2.Options)) (*ec2.DeleteRouteTableOutput, error) {
	return m.deleteRouteTableFunc(ctx, input, optFns...)
}

func (m *MockEC2Client) DescribeInternetGateways(ctx context.Context, input *ec2.DescribeInternetGatewaysInput, optFns ...func(*ec2.Options)) (*ec2.DescribeInternetGatewaysOutput, error) {
	return m.describeInternetGatewaysFunc(ctx, input, optFns...)
}

func (m *MockEC2Client) DetachInternetGateway(ctx context.Context, input *ec2.DetachInternetGatewayInput, optFns ...func(*ec2.Options)) (*ec2.DetachInternetGatewayOutput, error) {
	return m.detachInternetGatewayFunc(ctx, input, optFns...)
}

func (m *MockEC2Client) DeleteInternetGateway(ctx context.Context, input *ec2.DeleteInternetGatewayInput, optFns ...func(*ec2.Options)) (*ec2.DeleteInternetGatewayOutput, error) {
	return m.deleteInternetGatewayFunc(ctx, input, optFns...)
}

func (m *MockEC2Client) DescribeSecurityGroups(ctx context.Context, input *ec2.DescribeSecurityGroupsInput, optFns ...func(*ec2.Options)) (*ec2.DescribeSecurityGroupsOutput, error) {
	return m.describeSecurityGroupsFunc(ctx, input, optFns...)
}

func (m *MockEC2Client) DeleteSecurityGroup(ctx context.Context, input *ec2.DeleteSecurityGroupInput, optFns ...func(*ec2.Options)) (*ec2.DeleteSecurityGroupOutput, error) {
	return m.deleteSecurityGroupFunc(ctx, input, optFns...)
}

func (m *MockEC2Client) DescribeNetworkAcls(ctx context.Context, input *ec2.DescribeNetworkAclsInput, optFns ...func(*ec2.Options)) (*ec2.DescribeNetworkAclsOutput, error) {
	return m.describeNetworkAclsFunc(ctx, input, optFns...)
}

func (m *MockEC2Client) DeleteNetworkAcl(ctx context.Context, input *ec2.DeleteNetworkAclInput, optFns ...func(*ec2.Options)) (*ec2.DeleteNetworkAclOutput, error) {
	return m.deleteNetworkAclFunc(ctx, input, optFns...)
}
